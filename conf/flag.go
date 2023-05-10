package conf

import (
	"flag"

	"testing"
)

var EnvPath string

func flaginit() {
	flag.StringVar(&EnvPath, "e", "./env.yaml", "env yaml file path")
	testing.Init()
	flag.Parse()
}
