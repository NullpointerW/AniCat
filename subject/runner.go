package subject

import (
	"context"
	"log"
	"time"

	CFG "github.com/NullpointerW/anicat/conf"
	qbt "github.com/NullpointerW/go-qbittorrent-apiv2"

	DL "github.com/NullpointerW/anicat/download"
	"github.com/NullpointerW/anicat/download/rss"
	TORR "github.com/NullpointerW/anicat/download/torrent"
	"github.com/NullpointerW/anicat/errs"
	P "github.com/NullpointerW/anicat/pusher"
	"github.com/NullpointerW/anicat/pusher/email"
	"github.com/NullpointerW/anicat/util"
)

// before gorountie handle it init inner channels and ctxfunc
func (s *Subject) runtimeInit(reload bool) {
	if s.Terminate {
		Manager.Add(s)
		return
	}
	c := context.Background()
	ctx, exit := context.WithCancel(c)
	s.Exit = exit
	s.PushChan = make(chan qbt.Torrent, 1024)
	s.Exited = make(chan struct{})
	if s.Pushed == nil {
		s.Pushed = make(map[string]string)
	}
	Manager.Add(s)
	go s.run(ctx, reload)
}

func (s *Subject) run(ctx context.Context, reload bool) {
	if reload {
		util.Debugf("subj reload sid:%d", s.SubjId)
		s.checkDL()
	}
	t := time.NewTicker(30 * time.Minute)
	for {
		select {
		case torr := <-s.PushChan:
			err := s.push(torr, email.Sender{})
			if err != nil {
				log.Println(err)
			}
		case <-ctx.Done():
			util.Debugf("subj exited sid:%d", s.SubjId)
			exit(s)
			return
		case <-t.C:
			util.Debugf("subj update mission started sid:%d", s.SubjId)
			err := s.update()
			if err != nil {
				log.Println(err)
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
		log.Println(err)
	}
	close(s.Exited)
	close(s.PushChan)
}

func (s *Subject) checkDL() (err error) {
	if s.ResourceTyp == Torrent {
		util.Debugln("subj:", s.SubjId, "is torr typ,start check DL")
		compl, err := TORR.DLcompl(s.TorrentHash)
		if err != nil {
			return err
		} else if compl {
			util.Debugln("subj:", s.SubjId, "DL fin terminate now")
			s.terminate()
			return err
		}
	} else if s.ResourceTyp == RSS && s.Finished {
		util.Debugln("subj:", s.SubjId, "is rss typ,start check DL")
		if s.Typ == TV && s.EndTime != "" {
			util.Debugln("subj:", s.SubjId, "is TV and epi fin ")
			e, err := util.ParseTime(s.EndTime)
			util.Debugln("subj:", s.SubjId, "epi endtime is ", util.ParseTimeStr(e))
			if err != nil {
				return err
			}
			if time.Since(e) >= util.Day {
				util.Debugln("subj:", s.SubjId, "The time elapsed since the end of the anime is more than 1 day. ")
				goto checkSync
			}
			util.Debugln("subj:", s.SubjId, "The time elapsed DAY between the end of the anime and now is", time.Since(e).Hours()/24)
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

// call only when the subject epis is fin
func (s *Subject) RssDLSynced() (bool, error) {
	arts, err := rss.GetMatchedArts(s.RssPath())
	if err != nil {
		return false, nil
	}
	tlen := len(arts)
	if tlen == 0 {
		log.Println("there is no arts matched , check the rss match rule!", "sid:", s.SubjId)
		return true, nil
	}
	util.Debugln("rss total len is", tlen, "sid is", s.SubjId)
	// hs, err := TORR.GetViaPath(s.Path)
	// if err != nil {
	// 	return false, err
	// }
	// c := 0
	// for _, h := range hs {
	// 	if h.Progress == 1 {
	// 		c++
	// 	}
	// }
	c := len(s.RssTorrents)
	util.Debugf("subj sid:%d total series:%d local series:%d,local cmpl series:%d ", s.SubjId, tlen, c, c)
	return c >= tlen, nil
}

func (s *Subject) push(torr qbt.Torrent, pusher P.Pusher) error {
	if s.ResourceTyp == Torrent {
		if torr.Hash == s.TorrentHash {
			err := pusher.Push(P.Payload{
				SubjectId:    s.SubjId,
				SubjectName:  s.Name,
				DownLoadName: torr.Name,
				Size:         torr.Size,
			})
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
	s.RssTorrents[torr.Hash] = struct{}{}
	if s.Typ == TV {
		rename, err := Rename(s, torr)
		if err != nil {
			return err
		}
		if th, e := s.Pushed[rename]; e {
			merr := errs.MultiErr{}
			dumpliErr := errs.Custom("%w:origin_name=%s,rename:%s", errs.ErrItemAlreadyPushed, torr.Name, rename)
			merr.Add(dumpliErr)
			if CFG.Env.DropOnDumplicate && th != torr.Hash {
				log.Println("delete ", torr.Name)
				merr.Add(DL.Qbt.DelTorrentsFs(torr.Hash))
			}
			return merr.Err()
		}
		err = DL.Qbt.RenameFile(torr.Hash, torr.Name, rename)
		if err != nil {
			return err
		}
		mErr := errs.MultiErr{}
		s.Pushed[rename] = torr.Hash
		err = pusher.Push(P.Payload{
			SubjectId:    s.SubjId,
			SubjectName:  s.Name,
			DownLoadName: torr.Name,
			Size:         torr.Size,
			Episode:      util.TrimExtensionAndGetEpi(rename),
		})
		mErr.Add(err)
		mErr.Add(s.writeJson())
		return mErr.Err()
	}
	return nil

}

func (s *Subject) terminate() {
	util.Debugf("subj sid:%d terminate ", s.SubjId)
	s.Terminate, s.Finished = true, true
	s.writeJson()
	s.Exit()
}
