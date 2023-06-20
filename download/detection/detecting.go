package detection

import (
	// "fmt"
	"log"
	"strconv"
	"strings"
	"time"

	qbt "github.com/NullpointerW/go-qbittorrent-apiv2"

	DL "github.com/NullpointerW/anicat/download"
	"github.com/NullpointerW/anicat/download/torrent"
	"github.com/NullpointerW/anicat/errs"
	"github.com/NullpointerW/anicat/subject"
	"github.com/NullpointerW/anicat/util"
)

func Detect() {
	for {
		sync, err := DL.Qbt.GetMainData()
		if err == nil {
			for h, t := range sync.Torrents {
				if t.Progress == 1 {
					var (
						sid int
						err error
					)
					torr, err := torrent.Get(h)
					if err != nil {
						log.Println(err)
						continue
					}
					util.Debugf("detcting---->torrfn:%s,savepath:%s,tag:%s \n", torr.Name, torr.SavePath, torr.Tags)
					if strings.Contains(torr.Tags, subject.QbtTag_prefix) {
						s := strings.ReplaceAll(torr.Tags, subject.QbtTag_prefix, "")
						sid, err = strconv.Atoi(s)
						if err != nil {
							log.Println(err)
							goto viaSP
						}
						err = send(sid, torr)
						if err != nil {
							log.Println(err)
						}
						continue
					} else {
						goto viaSP
					}
				viaSP:
					sid = subject.Manager.GetSidViaSp(util.FileSeparatorConv(torr.SavePath))
					if err != nil {
						log.Println(err)
						continue
					}
					err = send(sid, torr)
					if err != nil {
						log.Println(err)
					}
					continue
				}
			}
		} else {
			log.Println(err)
		}
		time.Sleep(20 * time.Second)
	}

}

func send(sid int, torr qbt.Torrent) error {
	s := subject.Manager.Get(sid)

	if s == nil {
		return errs.Custom("%w:sid:%d", errs.ErrSubjectNotFound, sid)
	}
	if s.Terminate {
		return nil
	}

	log.Printf("pushing--->torrfn:%s,savepath:%s,tag:%s \n", torr.Name, torr.SavePath, torr.Tags)
	
	select {
	case <-s.Exited:
	default:
		s.PushChan <- torr
	}
	return nil
}
