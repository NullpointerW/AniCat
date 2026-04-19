package main

import (
	"github.com/NullpointerW/anicat/downloader/detector"
	netsrv "github.com/NullpointerW/anicat/net/server"
	"github.com/NullpointerW/anicat/subject"
	webserver "github.com/NullpointerW/anicat/web/server"
	_ "github.com/NullpointerW/anicat/debug"
)

func main() {
	subject.Scan()
	go subject.StartManagement()
	go detector.Detect()
	go netsrv.Listen()
	go webserver.Listen()
	select {}
}
