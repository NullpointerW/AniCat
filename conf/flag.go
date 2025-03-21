package conf

import (
	"flag"
	"testing"
)

var (
	EnvPath      string
	debug        bool
	IdeDebugging bool
	Testing      bool
)

func flagInit() {
	flag.BoolVar(&IdeDebugging, "idebug", false, "IdeDebugging mode")
	flag.StringVar(&EnvPath, "e", "./env.yaml", "env yaml filepath")
	flag.BoolVar(&debug, "d", false, "debug mode")
	// use Testing flag to skip env.yaml init and qbittorrent cli init
	// to avoid panic on go testing
	// use command:
	// 	go test ... -args -t=true
	flag.BoolVar(&Testing, "t", false, "testing mode")
	testing.Init()
	flag.Parse()

	// Testing=true

	if Testing {
		IdeDebugging = true
	}
	logInit(debug || IdeDebugging)
}
