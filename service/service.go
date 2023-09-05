package main

import (
	"fmt"
	"github.com/NullpointerW/anicat/download/detection"
	"github.com/NullpointerW/anicat/log"
	netsrv "github.com/NullpointerW/anicat/net/server"
	"github.com/NullpointerW/anicat/subject"
	util "github.com/NullpointerW/anicat/utils"
	"github.com/kardianos/service"
	"os"
	"path/filepath"
)

func main() {
	Run := func() {
		subject.Scan()
		go subject.StartManagement()
		go detection.Detect()
		go netsrv.Listen()
	}
	p := program{Run: Run}
	p.service()
}

type program struct {
	Run func()
}

func (p *program) Start(s service.Service) error {
	go p.run()
	return nil
}

func (p *program) run() {
	p.Run()
}

func (p *program) Stop(s service.Service) error {
	return nil
}

func (p *program) service() {
	executePath, err := os.Executable()
	if err != nil {
		fmt.Println("get executePath failed:", err)
	}
	executePath, err = filepath.EvalSymlinks(filepath.Dir(executePath))
	if err != nil {
		fmt.Println("get executePath failed:", err)
	}
	executePath += "/env.yaml"
	executePath = util.FileSeparatorConv(executePath)
	svcConfig := &service.Config{
		Name:        "AniCat",
		DisplayName: "AniCat",
		Description: "auto download service",
		Arguments:   []string{"-e", executePath},
	}
	s, err := service.New(p, svcConfig)
	if err != nil {
		log.Error(log.Struct{"err", err}, "create service failed")
		os.Exit(1)
	}
	if len(os.Args) > 1 {
		switch {
		case os.Args[1] == "install":
			err := s.Install()
			if err != nil {
				fmt.Println("error:", err)
				return
			}
			fmt.Println("service install succeeded")
			return
		case os.Args[1] == "uninstall":
			err := s.Uninstall()
			if err != nil {
				fmt.Println("error:", err)
				return
			}
			fmt.Println("service uninstall succeeded")
			return
		case os.Args[1] == "start":
			err := s.Start()
			if err != nil {
				fmt.Println("error:", err)
				return
			}
			fmt.Println("service start succeeded")
			return
		default:
			goto exec
		}
	}
exec:
	err = s.Run()
	log.Error(log.Struct{"err", err}, "service terminal")

}
