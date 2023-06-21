package conf

import (
	"flag"
	"testing"

	"github.com/NullpointerW/anicat/util"
)

var (
	EnvPath string
	debug   bool
)

func flaginit() {
	flag.StringVar(&EnvPath, "e", "./env.yaml", "env yaml filepath")
	flag.BoolVar(&debug, "d", false, "debug mode")
	testing.Init()
	flag.Parse()
	if debug {
		util.DebugEnv()
	}
	loginit(debug)
}
