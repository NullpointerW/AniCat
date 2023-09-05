package main

import (
	"github.com/NullpointerW/anicat/download/detection"
	netsrv "github.com/NullpointerW/anicat/net/server"
	"github.com/NullpointerW/anicat/subject"
)

func main() {
	Run := func() {
		subject.Scan()
		go subject.StartManagement()
		go detection.Detect()
		go netsrv.Listen()
		//select {}
	}
	p := program{Run: Run}
	p.service()
}
