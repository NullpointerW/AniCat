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

var (
	Ver    = "x.x.x"
	projlk = "https://github.com/NullpointerW/AniCat"
)

var Env Environment

type Environment struct {
	Port int `yaml:"port"`
	Qbt  struct {
		Url          string `yaml:"url"`
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

func (env *Environment) Print() {
	log.Println("port:", env.Port)
	log.Println("subject path:", env.SubjPath)

	if len(env.Proxies) != 0 {
		log.Println("scraper proxies:", env.Proxies)
	}
	if env.DropOnDumplicate {
		log.Println("drop dumplicate:", "yes")
	}

	log.Println("qbt weburl:",env.Qbt.Url)
	log.Println("qbt api request timeout(ms):",env.Qbt.Timeout)
}

func (env *Environment) EmailPrint() {
	log.Println("host:", env.Pusher.Email.Host)
	log.Println("port:", env.Pusher.Email.Port)
	log.Println("username:",env.Pusher.Email.Username)
}

func init() {
	flaginit()
	b, err := os.ReadFile(EnvPath)
	errs.PanicErr(err)
	errs.PanicErr(yaml.Unmarshal(b, &Env))
	if Env.Qbt.Timeout <= 0 {
		Env.Qbt.Timeout = 3000
	}
	log.Println("AniCat", "Ver."+Ver, "github:"+projlk)
	// log.Printf("env:\n%#+v\n", Env)
	Env.Print()
}

func loginit(debug bool) {
	flag := log.Ldate | log.Lmicroseconds
	if debug {
		flag |= log.Lshortfile
	}
	log.SetFlags(flag)

	if runtime.GOOS == "windows" && !IDEdebugging {
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
