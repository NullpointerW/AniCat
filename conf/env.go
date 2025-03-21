package conf

import (
	"fmt"
	"github.com/NullpointerW/anicat/errs"
	"github.com/NullpointerW/anicat/log"
	eslog "github.com/NullpointerW/anicat/pkg/log"
	util "github.com/NullpointerW/anicat/utils"
	"gopkg.in/yaml.v3"
	"os"
	"path/filepath"
	"runtime"
)

var (
	Ver    = "unknown"
	projlk = "https://github.com/NullpointerW/AniCat"
	SrvCTL = len(os.Args) > 1 && (os.Args[1] == "install" || os.Args[1] == "uninstall" || os.Args[1] == "start")
)

var Env Environment
var BgmiLogger *eslog.EnhanceLogger

type Environment struct {
	Port            int    `yaml:"port"`
	SubjPath        string `yaml:"path"`
	DropOnDuplicate bool   `yaml:"drop-duplicate"`
	RssFilter       struct {
		Contain   []string `yaml:"contain"`
		Exclusion []string `yaml:"exclusion"`
	} `yaml:"rss-filter"`
	Crawl struct {
		Proxies []string `yaml:"proxies"`
	} `yaml:"crawl"`
	Qbt struct {
		Url             string `yaml:"url"`
		Username        string `yaml:"username"`
		Password        string `yaml:"password"`
		LocalConnect    bool   `yaml:"localed"`
		Timeout         int    `yaml:"timeout"`
		TrackerProvider string `yaml:"tracker-provider"`
		Proxy           struct {
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
	BgmiLog           bool `yaml:"bangumi-log"`
	BuiltinDownloader bool `yaml:"builtin-downloader"`
}

func (env *Environment) Print() {
	logStruct := log.Struct{"port", env.Port, "subjectPath", env.SubjPath}
	if env.DropOnDuplicate {
		logStruct.Append("drop-onDuplicate", "yes")
	}
	log.Info(logStruct, "basicSetting")
	logStruct.Clear()
	if env.EnabledFilter() {
		logStruct.Append("globalFilter", "enable", "containWords", env.RssFilter.Contain, "exclusionWords",
			env.RssFilter.Exclusion)
		log.Info(logStruct, "global filter setting")
		logStruct.Clear()
	}

	if len(env.Crawl.Proxies) != 0 {
		logStruct.Append("scraperProxies", env.Crawl.Proxies)
		log.Info(logStruct, "crawling setting")
		logStruct.Clear()
	}
	if !env.BuiltinDownloader {
		logStruct.Append("qbt-webUrl", env.Qbt.Url, "qbt-apiRequestTimeout(ms)", env.Qbt.Timeout)
		log.Info(logStruct, "qbt setting")
	}
}

func (env *Environment) EmailPrint() {
	log.Info(log.Struct{"host", env.Pusher.Email.Host, "port", env.Pusher.Email.Port, "username", env.Pusher.Email.Username}, "SMTP setting")
}

func (env *Environment) EnabledFilter() bool {
	return len(env.RssFilter.Contain) > 0 || len(env.RssFilter.Exclusion) > 0
}

func init() {
	if SrvCTL {
		return
	}
	flagInit()
	if Testing {
		// skip
		return
	}
	b, err := os.ReadFile(EnvPath)
	errCallbackFunc := func() {
		if runtime.GOOS == "windows" {
			log.Error(log.Struct{"err", err}, "PANIC! process crashed")
		}
	}
	errs.PanicErr(err, errCallbackFunc)
	errs.PanicErr(yaml.Unmarshal(b, &Env), errCallbackFunc)
	if Env.Qbt.Timeout <= 0 {
		Env.Qbt.Timeout = 3000
	}
	log.Info(log.Struct{"ver", Ver, "github", projlk}, "AniCat")
	Env.Print()
	if Env.BgmiLog {
		executePath, err := util.GetExecutePath()
		if err != nil {
			log.Error(log.Struct{"err", err}, "open bgmi-log failed")
			Env.BgmiLog = false
			return
		}
		executePath += "/bangumi.log"
		output, err := os.OpenFile(executePath, os.O_APPEND|os.O_CREATE, 0777)
		if err != nil {
			log.Error(log.Struct{"err", err}, "open bgmi-log failed")
			Env.BgmiLog = false
			return
		}
		BgmiLogger = eslog.New("text", "info", "2006-01-02T15:04:05", false, 0, output)
	}
}

func logInit(debug bool) {
	var level = "info"
	if debug {
		level = "debug"
	}
	output := os.Stdout
	if runtime.GOOS == "windows" && !IdeDebugging {
		var err error
		executePath, err := util.GetExecutePath()
		executePath += "/output.log"
		if err != nil {
			fmt.Println(err)
			executePath = "." + executePath
		}
		executePath = util.FileSeparatorConv(executePath)
		output, err = os.OpenFile(executePath, os.O_TRUNC|os.O_CREATE, 0777)
		if err != nil {
			output = os.Stdout
			defer log.Error(log.Struct{"err", err}, "create logfile failed")
		}
		// receive PANIC
		dirP := filepath.Dir(executePath)
		PanicP := filepath.Join(dirP, "panic.log")
		panicOp, err := os.OpenFile(PanicP, os.O_APPEND|os.O_CREATE, 0777)
		if err == nil {
			err = errs.PanicRedirect(panicOp)
			if err != nil {
				defer log.Error(log.Struct{"err", err}, "redirect stderr failed")
				_ = panicOp.Close()
			}
		}

	}
	log.Init("text", level, "2006-01-02T15:04:05", debug, output)
	log.Debug(log.Struct{"os", runtime.GOOS}, "debug mode")

}
