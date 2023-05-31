package conf

import (
	"flag"
	"os"

	"testing"

	"github.com/NullpointerW/mikanani/util"
)

var EnvPath string

func flaginit() {
	flag.StringVar(&EnvPath, "e", "./env.yaml", "env yaml file path")
	var debug bool
	if len(os.Args)>1{
		args:=os.Args[1:]
		for _,a :=range args{
			if debug=a=="d";debug{
				util.DebugEnv()
			}
		}
	}
	loginit(debug)
	testing.Init()
	flag.Parse()
}
