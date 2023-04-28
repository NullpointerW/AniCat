package detection

import (
	qbt "github.com/NullpointerW/go-qbittorrent-apiv2"
	"github.com/NullpointerW/mikanani/subject"
	"time"
)

func init() {
	go detect()
}

func detect() {
	subject.Wg.Wait()
	for {
		//TODO
		time.Sleep(5 * time.Minute)
	}

}
