package conf

import (
	// "io"
	"runtime"

	"log"
	"os"

	"github.com/NullpointerW/anicat/errs"
	util "github.com/NullpointerW/anicat/utils"
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
		Timeout      int    `yaml:"timeout"`
	} `yaml:"qbittorrent"`
	Proxies          []string `yaml:"proxies"`
	SubjPath         string   `yaml:"path"`
	DropOnDumplicate bool     `yaml:"dropDumplicate"`
	Pusher           struct {
		Email struct {
			Host         string `yaml:"host"`
			Port         int    `yaml:"port"`
			Username     string `yaml:"username"`
			Password     string `yaml:"password"`
			TemplatePath string `yaml:"template"`
			SkipSSL      bool   `yaml:"skipssl"`
		} `yaml:"email"`
	} `yaml:"push"`
}

func init() {
	flaginit()
	b, err := os.ReadFile(EnvPath)
	errs.PanicErr(err)
	errs.PanicErr(yaml.Unmarshal(b, &Env))
	if Env.Qbt.Timeout <= 0 {
		Env.Qbt.Timeout = 3000
	}
	log.Printf("env:\n%#+v\n", Env)
}

func loginit(debug bool) {
	flag := log.Ldate | log.Lmicroseconds
	if debug {
		flag |= log.Lshortfile
	}
	log.SetFlags(flag)

	if runtime.GOOS == "windows"&&!IDEdebugging {
		f, err := os.OpenFile("./output.log", os.O_TRUNC|os.O_CREATE, 0777)
		if err != nil {
			log.Println(err)
			return
		}
		log.SetOutput(f)
	}
	// else {
	// 	mio := io.MultiWriter(os.Stderr, f)
	// 	log.SetOutput(mio)
	// }
	util.Debugln("os:", runtime.GOOS, "debug mode")
}
