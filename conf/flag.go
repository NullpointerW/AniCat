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
)

func flaginit() {
	flag.BoolVar(&IDEdebugging, "idebug", false, "IDEdebugging mode")
	flag.StringVar(&EnvPath, "e", "./env.yaml", "env yaml filepath")
	flag.BoolVar(&debug, "d", false, "debug mode")
	testing.Init()
	flag.Parse()
	if debug {
		util.DebugEnv()
	}
	loginit(debug || IDEdebugging)
}
