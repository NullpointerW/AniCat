package subject

import (
	"context"
	"log"
	"time"

	qbt "github.com/NullpointerW/go-qbittorrent-apiv2"

	"github.com/NullpointerW/mikanani/download/rss"
	TORR "github.com/NullpointerW/mikanani/download/torrent"
	"github.com/NullpointerW/mikanani/errs"
	"github.com/NullpointerW/mikanani/pusher"
	"github.com/NullpointerW/mikanani/util"
)

// before gorountie handle it init inner channels and ctxfunc
func (s *Subject) runtimeInit(reload bool) {
	if s.Terminate {
		Manager.Add(s)
		return
	}
	c := context.Background()
	ctx, exit := context.WithCancel(c)
	s.exit = exit
	s.PushChan = make(chan qbt.Torrent, 1024)
	s.Exited = make(chan struct{})
	Manager.Add(s)
	go s.run(ctx, reload)
}

func (s *Subject) run(ctx context.Context, reload bool) {
	if reload {
		util.Debugf("subj reload sid:%d", s.SubjId)
		s.checkDL()
	}
	t := time.NewTicker(util.Day)
	for {
		select {
		case torr := <-s.PushChan:
			err := s.push(torr)
			log.Println(err)
		case <-ctx.Done():
			util.Debugf("subj exited sid:%d", s.SubjId)
			exit(s)
			return
		case <-t.C:
			util.Debugf("subj update mission started sid:%d", s.SubjId)
			s.update()
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
	close(s.Exited)
	close(s.PushChan)
}

func (s *Subject) checkDL() (err error) {
	if s.ResourceTyp == Torrent {
		compl, err := TORR.DLcompl(s.TorrentHash)
		if err != nil {
			return err
		} else if compl {
			s.terminate()
			return err
		}
	} else if s.ResourceTyp == RSS && s.Finished {
		if s.Typ == TV && s.EndTime != "" {
			e, err := util.ParseTime(s.EndTime)
			if err != nil {
				return err
			}
			if time.Since(e) >= util.Day {
				goto checkSync
			}
		} else if s.Typ == MOVIE {
			goto checkSync
		}
	} else {
		return
	}
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

func (s *Subject) RssDLSynced() (bool, error) {
	arts, err := rss.GetMatchedArts(s.RssPath())
	if err != nil {
		return false, nil
	}
	tlen := len(arts)
	hs, err := TORR.GetViaPath(s.Path)
	if err != nil {
		return false, nil
	}
	llen := len(hs)
	util.Debugf("subj sid:%d total series:%d local series:%d ", s.SubjId, tlen, llen)
	return llen >= tlen, nil
}

func (s *Subject) push(torr qbt.Torrent) error {
	if s.ResourceTyp == Torrent {
		pusher.Push()
		s.terminate()
	}
	return nil
}

func (s *Subject) terminate() {
	util.Debugf("subj sid:%d terminate ", s.SubjId)
	s.Terminate, s.Finished = true, true
	s.writeJson()
	s.exit()
}
