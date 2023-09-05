package main

import (
	"github.com/NullpointerW/anicat/download/detection"
	netsrv "github.com/NullpointerW/anicat/net/server"
	"github.com/NullpointerW/anicat/subject"
)

func main() {
	subject.Scan()
	go subject.StartManagement()
	go detection.Detect()
	go netsrv.Listen()
	select {}
}
