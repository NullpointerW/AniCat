package subject

import (
	"context"
	qbt "github.com/NullpointerW/go-qbittorrent-apiv2"
	"github.com/NullpointerW/mikanani/pusher"
	"log"
	"time"
)

func run(subj *Subject, ctx context.Context, reload bool) {
	if reload {
          
	}
	t := time.NewTicker(24 * time.Hour)
	select {
	case torr := <-subj.PushChan:
		err := checkDL(subj, torr)
		log.Println(err)
	case <-ctx.Done():
		close(subj.exited)
		return
	case <-t.C:
		subj.update()
	}

}

func (s *Subject) update() {

}

func exit(s *Subject) {
	s.exit()
	close(s.exited)
	close(s.PushChan)
}

func checkDL(s *Subject, torr qbt.Torrent) error {
	if s.ResourceTyp == Torrent {
		pusher.Push()
		s.Terminal, s.Finished = true, true
		exit(s)
	}
	return nil
}
