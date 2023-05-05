package subject

import (
	"context"
	"log"
	"time"

	qbt "github.com/NullpointerW/go-qbittorrent-apiv2"
	TORR "github.com/NullpointerW/mikanani/download/torrent"
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
	s.PushChan = make(chan qbt.Torrent, 2)
	s.Exited = make(chan struct{})
	Manager.Add(s)
	go s.run(ctx, reload)
}

func (s *Subject) run(ctx context.Context, reload bool) {
	if reload {
		s.checkDL()
	}
	t := time.NewTicker(util.Day)
	for {
		select {
		case torr := <-s.PushChan:
			err := s.push(torr)
			log.Println(err)
		case <-ctx.Done():
			exit(s)
			return
		case <-t.C:
			s.update()
		}
	}

}

func (s *Subject) update() {
         
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
		}
	}
	return err
}

func (s *Subject) push(torr qbt.Torrent) error {
	if s.ResourceTyp == Torrent {
		pusher.Push()
		s.terminate()
	}
	return nil
}

func (s *Subject) terminate() {
	s.Terminate, s.Finished = true, true
	s.writeJson()
	exit(s)
}
