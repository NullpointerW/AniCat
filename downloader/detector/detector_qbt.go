package detector

import (
	"fmt"
	CFG "github.com/NullpointerW/anicat/conf"
	DL "github.com/NullpointerW/anicat/downloader"
	"github.com/NullpointerW/anicat/downloader/torrent"
	"github.com/NullpointerW/anicat/errs"
	"github.com/NullpointerW/anicat/log"
	"github.com/NullpointerW/anicat/subject"
	qbt "github.com/NullpointerW/go-qbittorrent-apiv2"
	"strconv"
	"strings"
	"time"
)

func Detect() {
	if CFG.Env.BuiltinDownloader {
		return
	}
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
						log.Error(log.Struct{"err", err}, "detector: get torrent failed")
						continue
					}
					//log.Debug(log.Struct{"torrfn", torr.Name, "savepath", torr.SavePath, "tag", torr.Tags, "categ", torr.Category},
					//	"detected completed downloader")
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
							log.Error(log.Struct{"err", err}, "detector: can not convert subjectId")
							continue
						}
						err = send(sid, torr)
						if err != nil {
							log.Error(log.Struct{"err", err}, "detector: send downloadEvent failed")
						}
					}
				}
			}
		} else {
			log.Error(log.Struct{"err", err}, "detector: get qbtSyncData failed")
		}
		time.Sleep(20 * time.Second)
		//time.Sleep(10 * time.Second)
	}
}

func send(sid int, torr qbt.Torrent) error {
	s := subject.Mgr.Get(sid)
	if s == nil {
		return fmt.Errorf("%w:sid:%d", errs.ErrSubjectNotFound, sid)
	}
	if s.Terminate {
		return nil
	}
	log.Debug(log.Struct{"torrfn", torr.Name, "savepath", torr.SavePath, "tag", torr.Tags, "categ", torr.Category},
		"pushing completed downloadEvent")
	select {
	case <-s.Exited:
	default:
		s.PushChan <- torr
	}
	return nil
}
