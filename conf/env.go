package conf

import (
	"log"
	"os"

	"github.com/NullpointerW/mikanani/errs"
	"gopkg.in/yaml.v3"
)

var Env Environment

type Environment struct {
	Port int `yaml:"port"`
	Qbt  struct {
		Host         string `yaml:"host"`
		Username     string `yaml:"username"`
		Password     string `yaml:"password"`
		LocalConnect bool   `yaml:"localed"`
	} `yaml:"qbittorrent"`
}

func init() {
	log.SetFlags(log.Lshortfile | log.Ldate | log.Lmicroseconds)
	flaginit()
	b, err := os.ReadFile(EnvPath)
	errs.PanicErr(err)
	errs.PanicErr(yaml.Unmarshal(b, &Env))
}
