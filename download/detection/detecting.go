package detection

import (
	// "fmt"
	"fmt"
	"log"
	"strconv"
	"strings"
	"time"

	qbt "github.com/NullpointerW/go-qbittorrent-apiv2"

	DL "github.com/NullpointerW/anicat/download"
	"github.com/NullpointerW/anicat/download/torrent"
	"github.com/NullpointerW/anicat/errs"
	"github.com/NullpointerW/anicat/subject"
	util "github.com/NullpointerW/anicat/utils"
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
					util.Debugf("detcting---->torrfn:%s,savepath:%s,tag:%s,categ:%s \n", torr.Name, torr.SavePath,
						torr.Tags, torr.Category)
					if istorr, isrss := strings.Contains(torr.Tags, subject.QbtTag_prefix),
						strings.Contains(torr.Category, subject.QbtTag_prefix); istorr || isrss {
						var s string
						if istorr {
							s = strings.ReplaceAll(torr.Tags, subject.QbtTag_prefix, "")
						} else {
							s = strings.ReplaceAll(torr.Category, subject.QbtTag_prefix, "")
						}
						sid, err = strconv.Atoi(s)
						if err != nil {
							log.Println(err)
							continue
						}
						err = send(sid, torr)
						if err != nil {
							log.Println(err)
						}
					}
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
		return fmt.Errorf("%w:sid:%d", errs.ErrSubjectNotFound, sid)
	}
	if s.Terminate {
		return nil
	}

	log.Printf("pushing----> torrfn=%s,\nsavepath=%s,\ntag=%s,\ncateg=%s\n", torr.Name, torr.SavePath,
		torr.Tags, torr.Category)

	select {
	case <-s.Exited:
	default:
		s.PushChan <- torr
	}
	return nil
}
