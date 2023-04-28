package subject

import (
	"context"
	qbt "github.com/NullpointerW/go-qbittorrent-apiv2"
	TORR "github.com/NullpointerW/mikanani/download/torrent"
	"github.com/NullpointerW/mikanani/pusher"
	"log"
	"time"
)

func (s *Subject) run(ctx context.Context, reload bool) {
	if reload {
		s.checkDL()
	}
	t := time.NewTicker(24 * time.Hour)
	select {
	case torr := <-s.PushChan:
		err := s.push(torr)
		log.Println(err)
	case <-ctx.Done():
		close(s.exited)
		return
	case <-t.C:
		s.update()
	}

}

func (s *Subject) update() {

}

func exit(s *Subject) {
	s.exit()
	close(s.exited)
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
