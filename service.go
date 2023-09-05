package main

import (
	"fmt"
	"github.com/NullpointerW/anicat/log"
	"github.com/kardianos/service"
	"os"
)

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
	svcConfig := &service.Config{
		Name:        "AniCat",
		DisplayName: "AniCat",
		Description: "auto download service",
		Arguments:   []string{"-e", "D:\\anicat-srvd\\env.yaml"},
	}
	s, err := service.New(p, svcConfig)
	if err != nil {
		log.Error(log.Struct{"err", err}, "create service failed")
		os.Exit(1)
	}
	if len(os.Args) > 1 {
		if os.Args[1] == "install" {
			err := s.Install()
			if err != nil {
				fmt.Println("error:", err)
				return
			}
			fmt.Println("service install succeeded")
			return
		} else if os.Args[1] == "uninstall" {
			err := s.Uninstall()
			if err != nil {
				fmt.Println("error:", err)
				return
			}
			fmt.Println("service uninstall succeeded")
			return
		}
	}
	err = s.Run()
	log.Error(log.Struct{"err", err}, "service terminal")
}
