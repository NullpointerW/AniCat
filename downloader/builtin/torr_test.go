package main

import (
	"fmt"
	"io"
	"testing"
)

func TestReader(t *testing.T) {
	var r Reader = NewHttpReader()
	_, err := r.Reader("https://mikanani.me/Download/20240629/799328b2d580a66e25640fdea2d17302501eca08.torrent")
	if err != nil {
		fmt.Println(err)
	}
	rr, err := r.Reader("https://mikanani.me/Download/20240629/799328b2d580a66e25640fdea2d17302501eca08.x")
	if err != nil {
		fmt.Println(err)
	}
	all, err := io.ReadAll(rr)
	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Println(string(all))
	}
	_, err = r.Reader("x-man")
	if err != nil {
		fmt.Println(err)
	}
}
