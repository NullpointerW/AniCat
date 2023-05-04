package detection

import (
	"log"
	"strconv"
	"strings"
	"time"

	qbt "github.com/NullpointerW/go-qbittorrent-apiv2"
	

	DL "github.com/NullpointerW/mikanani/download"
	"github.com/NullpointerW/mikanani/errs"
	"github.com/NullpointerW/mikanani/subject"
	"github.com/NullpointerW/mikanani/util"
)

func init() {
	go detect()
}

func detect() {
	subject.Wg.Wait()
	for {
		sync, err := DL.Qbt.GetMainData()
		if err == nil {
			for _, torr := range sync.Torrents {
				if torr.Progress == 1 {
					var (
						sid int
						err error
					)
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
					}
				viaSP:
					sp := util.FileSeparatorConv(torr.SavePath)
					ss := strings.Split(sp, "/")
					s := ss[len(ss)-1]
					s = strings.ReplaceAll(s, subject.FolderSuffix, "")
					sid, err = strconv.Atoi(s)
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
		time.Sleep(5 * time.Minute)
	}

}

func send(sid int, torr qbt.Torrent) error {
	s := subject.Manager.GetSubject(sid)

	if s == nil {
		return errs.Custom("%w:sid:%d", errs.ErrSubjectNotFound, sid)
	}
	select {
	case <-s.Exited:
	default:
		s.PushChan <- torr
	}

	return nil
}
