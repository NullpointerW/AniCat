package conf

import (
	"github.com/NullpointerW/anicat/errs"
	"github.com/NullpointerW/anicat/log"
	"gopkg.in/yaml.v3"
	"os"
	"runtime"
)

var (
	Ver    = "x.x.x"
	projlk = "https://github.com/NullpointerW/AniCat"
)

var Env Environment

type Environment struct {
	Port             int    `yaml:"port"`
	SubjPath         string `yaml:"path"`
	DropOnDumplicate bool   `yaml:"drop-dumplicate"`
	RssFilter        struct {
		Contain   []string `yaml:"contain"`
		Exclusion []string `yaml:"exclusion"`
	} `yaml:"rss-filter"`
	Crawl struct {
		Proxies []string `yaml:"proxies"`
	} `yaml:"crawl"`
	Qbt struct {
		Url          string `yaml:"url"`
		Username     string `yaml:"username"`
		Password     string `yaml:"password"`
		LocalConnect bool   `yaml:"localed"`
		Timeout      int    `yaml:"timeout"`
		Proxy        struct {
			Addr     string `yaml:"address"`
			Type     string `yaml:"type"`
			Username string `yaml:"username"`
			Password string `yaml:"password"`
			Peer     bool   `yaml:"peer"`
			TorrOnly bool   `yaml:"torrent-only"`
			Hslookup bool   `yaml:"host-lookup"`
		} `yaml:"proxy"`
	} `yaml:"qbittorrent"`
	Pusher struct {
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
	logStruct := log.Struct{"port", env.Port, "subject path", env.SubjPath}
	if env.DropOnDumplicate {
		logStruct.Append("drop dumplicate:", "yes")
	}
	log.Info(logStruct, "basic setting")
	logStruct.Clear()
	if env.EnabledFilter() {
		logStruct.Append("golbal filter", "enable", "contain words", env.RssFilter.Contain, "exclusion words", env.RssFilter.Exclusion)
		log.Info(logStruct, "golbal filter setting")
		logStruct.Clear()
	}

	if len(env.Crawl.Proxies) != 0 {
		logStruct.Append("scraper proxies", env.Crawl.Proxies)
		log.Info(logStruct, "crawling setting")
		logStruct.Clear()
	}
	logStruct.Append("qbt weburl", env.Qbt.Url, "qbt api request timeout(ms)", env.Qbt.Timeout)
	log.Info(logStruct, "qbt setting")
}

func (env *Environment) EmailPrint() {
	log.Info(log.Struct{"host", env.Pusher.Email.Host, "port:", env.Pusher.Email.Port, "username:", env.Pusher.Email.Username}, "SMTP setting")
}

func (env *Environment) EnabledFilter() bool {
	return len(env.RssFilter.Contain) > 0 || len(env.RssFilter.Exclusion) > 0
}

func init() {
	flaginit()
	if Testing {
		// skip
		return
	}
	b, err := os.ReadFile(EnvPath)
	errs.PanicErr(err)
	errs.PanicErr(yaml.Unmarshal(b, &Env))
	if Env.Qbt.Timeout <= 0 {
		Env.Qbt.Timeout = 3000
	}
	log.Info(log.Struct{"ver", Ver, "github", projlk}, "AniCat")
	Env.Print()
}

func loginit(debug bool) {
	var level = "info"
	if debug {
		level = "debug"
	}
	output := os.Stderr
	if runtime.GOOS == "windows" && !IDEdebugging {
		var err error
		output, err = os.OpenFile("./output.log", os.O_TRUNC|os.O_CREATE, 0777)
		if err != nil {
			defer log.Error(log.Struct{"err", err}, "create logfile failed")
		}
	}
	log.Init("text", level, "2006-01-02T15:04:05", debug, output)
	log.Debug(log.Struct{"os", runtime.GOOS}, "debug mode")
}
