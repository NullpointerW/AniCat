package conf

import (
	"flag"
	"testing"

	util "github.com/NullpointerW/anicat/utils"
)

var (
	EnvPath      string
	debug        bool
	IDEdebugging bool
	Testing      bool
)

func flaginit() {
	flag.BoolVar(&IDEdebugging, "idebug", false, "IDEdebugging mode")
	flag.StringVar(&EnvPath, "e", "./env.yaml", "env yaml filepath")
	flag.BoolVar(&debug, "d", false, "debug mode")
	// use Testing flag to skip env.yaml init and qbittorrent cli init
	// to avoid panic on go testing
	// use command:
	// 	go test ... -args -t=true
	flag.BoolVar(&Testing, "t", false, "testing mode")
	testing.Init()
	flag.Parse()
	if Testing{
		return
	}
	debug = debug || IDEdebugging
	if debug {
		util.DebugEnv()
	}
	loginit(debug)
}
