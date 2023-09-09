package main

import (
	"github.com/NullpointerW/anicat/downloader/detector"
	netsrv "github.com/NullpointerW/anicat/net/server"
	"github.com/NullpointerW/anicat/subject"
)

func main() {
	subject.Scan()
	go subject.StartManagement()
	go detector.Detect()
	go netsrv.Listen()
	select {}
}
