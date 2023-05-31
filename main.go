package main

import (
	"os"
	"os/signal"
	"syscall"

	"github.com/NullpointerW/mikanani/download/detection"
	netsrv "github.com/NullpointerW/mikanani/net/server"
	"github.com/NullpointerW/mikanani/subject"
)

func main() {
	subject.Scan()
	go subject.StartManagement()
	go detection.Detect()
	go netsrv.Listen()
	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)
	<-sigCh
}
