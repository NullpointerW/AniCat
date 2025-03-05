package builtin

import (
	"fmt"
	"sync"
	"testing"
)

func TestReader(t *testing.T) {
	var r TorrentSeeker = NewHttpSeeker()
	_, err := r.Seek("https://mikanani.me/Download/20240629/799328b2d580a66e25640fdea2d17302501eca08.torrent")
	if err != nil {
		fmt.Println(err)
	}
	ts, err := r.Seek("https://mikanani.me/Download/20240629/799328b2d580a66e25640fdea2d17302501eca08.x")
	if err != nil {
		fmt.Println(err)
	}

	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Println(*ts)
	}
	_, err = r.Seek("x-man")
	if err != nil {
		fmt.Println(err)
	}
	s:=sync.Pool{}
	s.Get()
}
