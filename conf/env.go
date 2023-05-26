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
	Proxies  []string `yaml:"proxies"`
	SubjPath string   `yaml:"path"`
	Pusher   struct {
		Email struct{
			Host string `yaml:"host"`
			Port int   `yaml:"port"`
			Username string `yaml:"username"`
			Password string  `yaml:"password"`
		}`yaml:"email"`
	} `yaml:"push"`
}

func init() {
	log.SetFlags(log.Lshortfile | log.Ldate | log.Lmicroseconds)
	flaginit()
	b, err := os.ReadFile(EnvPath)
	errs.PanicErr(err)
	errs.PanicErr(yaml.Unmarshal(b, &Env))
	log.Printf("env:\n%#+v\n", Env)
}
