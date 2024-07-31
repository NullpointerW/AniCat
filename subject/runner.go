package subject

import (
	"context"
	"fmt"
	CFG "github.com/NullpointerW/anicat/conf"
	DL "github.com/NullpointerW/anicat/downloader"
	"github.com/NullpointerW/anicat/downloader/rss"
	TORR "github.com/NullpointerW/anicat/downloader/torrent"
	"github.com/NullpointerW/anicat/errs"
	"github.com/NullpointerW/anicat/log"
	eslog "github.com/NullpointerW/anicat/pkg/log"
	P "github.com/NullpointerW/anicat/pusher"
	"github.com/NullpointerW/anicat/pusher/email"
	util "github.com/NullpointerW/anicat/utils"
	qbt "github.com/NullpointerW/go-qbittorrent-apiv2"
	"time"
)

// runtimeInit before goroutine handle it init inner channels and ctx func
func (s *Subject) runtimeInit(reload bool) {
	s.Exited = make(chan struct{})
	if s.Terminate {
		close(s.Exited)
		Mgr.Add(s)
		return
	}
	c := context.Background()
	ctx, exit := context.WithCancel(c)
	s.Exit = exit
	s.PushChan = make(chan qbt.Torrent, 1024)
	if s.Pushed == nil {
		s.Pushed = make(map[string]string)
	}
	Mgr.Add(s)
	go s.run(ctx, reload)
	go JellyfinMetaDataHelper(s.Path, s.FolderName, s.Exited)
}

func (s *Subject) run(ctx context.Context, reload bool) {
	Mgr.wg.Add(1)
	defer Mgr.wg.Done()
	if reload {
		log.Debug(log.Struct{"sid", s.SubjId}, "subject reload")
		//s.checkDL()
	}
	t := time.NewTicker(30 * time.Minute)
	for {
		select {
		case torr := <-s.PushChan:
			err := s.push(torr, email.Poster)
			if err != nil {
				log.Error(log.Struct{"sid", s.SubjId, "err", err}, "push process failed")
			}
		case <-ctx.Done():
			log.Debug(log.Struct{"sid", s.SubjId}, "runner exited")
			exit(s)
			return
		case <-t.C:
			log.Debug(log.Struct{"sid", s.SubjId}, "subject update mission started")
			err := s.update()
			if err != nil {
				log.Error(log.Struct{"sid", s.SubjId, "err", err}, "update mission failed")
			}
		}
	}

}

func (s *Subject) update() error {
	wrap := errs.ErrWrapper{}
	wrap.Handle(func() error {
		return s.FetchInfo()
	})
	wrap.Handle(func() error {
		return s.writeJson()
	})
	wrap.Handle(func() error {
		return s.checkDL()
	})
	return wrap.Error()
}

func exit(s *Subject) {
	err := s.writeJson()
	if err != nil {
		log.Error(log.Struct{"sid", s.SubjId, "err", err}, "write json failed while exited")
	}
	close(s.Exited)
	close(s.PushChan)
	Mgr.Sync()
}

func (s *Subject) checkDL() (err error) {
	if s.ResourceTyp == Torrent {
		log.Debug(log.Struct{"sid", s.SubjId, "type", "torrent"}, "start check DL")
		compl, err := TORR.DLcompl(s.TorrentHash)
		if err != nil {
			return err
		} else if compl {
			log.Debug(log.Struct{"sid", s.SubjId, "type", "torrent"}, "DL fin terminate now")
			s.terminate()
			return err
		}
	} else if s.ResourceTyp == RSS && s.Finished {
		log.Debug(log.Struct{"sid", s.SubjId, "type", "rss"}, "start check DL")
		if s.Typ == TV && s.EndTime != "" {
			log.Debug(log.Struct{"sid", s.SubjId, "resType", "TV"}, "epi fin")
			e, err := util.ParseTime(s.EndTime, util.YMDParseLayout)
			log.Debug(log.Struct{"sid", s.SubjId, "resType", "TV"}, "epi endtime is ", util.ParseTimeStr(e))
			if err != nil {
				return err
			}
			if time.Since(e) >= util.Day {
				log.Debug(log.Struct{"sid", s.SubjId, "resType", "TV"}, "The time elapsed since the end of the anime is more than 1 day. ")
				goto checkSync
			}
			log.Debug(log.Struct{"sid", s.SubjId, "resType", "TV"}, "The time elapsed DAY between the end of the anime and nowtime is ",
				time.Since(e).Hours()/24)
		} else if s.Typ == MOVIE {
			goto checkSync
		}
	}
	return
checkSync:
	sync, err := s.RssDLSynced()
	if err != nil {
		return err
	}
	if sync {
		s.terminate()
	}
	return
}

// RssDLSynced called only when the subject epis is fin
func (s *Subject) RssDLSynced() (bool, error) {
	arts, err := rss.GetMatchedArts(s.RssPath())
	if err != nil {
		return false, nil
	}
	tlen := len(arts)
	if tlen == 0 {
		log.Warn(log.Struct{"sid", s.SubjId, "resType", "RSS"}, "there is no arts matched,check the rss match rule!")
		return true, nil
	}
	log.Debug(log.Struct{"sid", s.SubjId, "rssTotalLen", tlen})
	c := len(s.RssTorrents)
	log.Debug(log.Struct{"sid", s.SubjId, "series", tlen, "localSeries", c, "cmplSeries", c})
	return c >= tlen, nil
}

func (s *Subject) push(torr qbt.Torrent, pusher P.Pusher) error {
	if s.ResourceTyp == Torrent {
		if torr.Hash == s.TorrentHash {
			var err error
			if s.Typ == TV {
				err = renameTorr(s, torr)
				if err != nil {
					goto term
				}
				epi := "S" + s.Season + "E01-" + fmt.Sprintf("%02d", s.Episode)
				err = pusher.Push(P.Payload{
					SubjectId:    s.SubjId,
					SubjectName:  s.Name,
					DownLoadName: torr.Name,
					Size:         torr.Size,
					Episode:      epi,
				})
			} else {
				err = pusher.Push(P.Payload{
					SubjectId:    s.SubjId,
					SubjectName:  s.Name,
					DownLoadName: torr.Name,
					Size:         torr.Size,
					Episode:      "MOVIE",
				})
			}
		term:
			s.terminate()
			return err
		}
		return nil
	}
	// RSS
	if s.Pushed == nil {
		s.Pushed = make(map[string]string)
	}
	if s.RssTorrents == nil {
		s.RssTorrents = map[string]struct{}{}
	}

	// perf: skip rename process
	if _, e := s.RssTorrents[torr.Hash]; e {
		log.Debug(log.Struct{"sid", s.SubjId, "torrfn", torr.Name, "torrHash", torr.Hash}, "skip rename")
		return nil
	}

	s.RssTorrents[torr.Hash] = struct{}{}
	if s.Typ == TV {
		var se = ""
		if checkSingleVideo(torr) {
			rename, err := renameTV(s, torr.Name)
			if err != nil {
				if CFG.Env.BgmiLog {
					CFG.BgmiLogger.Infof(eslog.Struct{"sid", s.SubjId, "name", s.Name}, "episode update(unnamed): %s", torr.Name)
				}
				return err
			}
			se = util.TrimExtensionAndGetEpi(rename)
			if th, e := s.Pushed[se]; e {
				merr := errs.MultiErr{}
				dumpliErr := fmt.Errorf("%w: origin_name=%s,rename=%s", errs.ErrItemAlreadyPushed, torr.Name, rename)
				merr.Add(dumpliErr)
				if CFG.Env.DropOnDuplicate && th != torr.Hash {
					log.Warn(log.Struct{"sid", s.SubjId, "torrfn", torr.Name, "torrHash", torr.Hash, "size", torr.Size}, "delete dumplicate file")
					merr.Add(DL.Qbt.DelTorrentsFs(torr.Hash))
				}
				return merr.Err()
			}
			err = DL.Qbt.RenameFile(torr.Hash, torr.Name, rename)
			if err != nil {
				return err
			}
			s.Pushed[se] = torr.Hash
			if CFG.Env.BgmiLog {
				CFG.BgmiLogger.Infof(eslog.Struct{"sid", s.SubjId, "name", s.Name}, "episode update: %s", rename)
			}
		} else {
			log.Info(log.Struct{"sid", s.SubjId, "torrfn", torr.Name, "torrHash", torr.Hash}, "not a videoFile,may external subtitles")
			ok, rn, err := renameSubRssTorr(s, torr)
			log.Error(log.Struct{"err", err}, "rename RssTorr with subtitles failed")
			if !ok {
				return nil
			}
			se = util.TrimExtensionAndGetEpi(rn)
		}
		mErr := errs.MultiErr{}
		err := pusher.Push(P.Payload{
			SubjectId:    s.SubjId,
			SubjectName:  s.Name,
			DownLoadName: torr.Name,
			Size:         torr.Size,
			Episode:      se,
		})
		mErr.Add(err)
		mErr.Add(s.writeJson())
		// FIXME if rss contains sp ,may be exit early
		episNum := s.Episode
		if episNum != 0 && len(s.Pushed) >= episNum {
			log.Info(log.Struct{"sid", s.SubjId, "resType", "RSS"}, "compiled,exited now")
			s.terminate()
		}
		return mErr.Err()
	} else { //Movie
		if _, e := s.Pushed[torr.Hash]; e {
			return fmt.Errorf("%w: name=%s", errs.ErrItemAlreadyPushed, torr.Name)
		}
		mErr := errs.MultiErr{}
		s.Pushed[torr.Hash] = ""
		err := pusher.Push(P.Payload{
			SubjectId:    s.SubjId,
			SubjectName:  s.Name,
			DownLoadName: torr.Name,
			Size:         torr.Size,
			Episode:      "Movie",
		})
		mErr.Add(err)
		mErr.Add(s.writeJson())
		return mErr.Err()
	}
}

func (s *Subject) terminate() {
	log.Debug(log.Struct{"sid", s.SubjId, "resType", s.ResourceTyp.String()}, "exited")
	s.Terminate, s.Finished = true, true
	err := s.writeJson()
	if err != nil {
		log.Info(log.Struct{"err", err}, "write json failed")
	}
	s.Exit()
}
