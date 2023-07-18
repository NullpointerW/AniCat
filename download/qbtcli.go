package download

import (
	"time"

	CFG "github.com/NullpointerW/anicat/conf"
	"github.com/NullpointerW/anicat/errs"
	qbt "github.com/NullpointerW/go-qbittorrent-apiv2"
)

var Qbt *qbt.Client

func init() {
	var (
		cli *qbt.Client
		err error
	)
	if CFG.Env.Qbt.LocalConnect {
		cli, err = qbt.NewCli(CFG.Env.Qbt.Host)
	} else {
		cli, err = qbt.NewCli(CFG.Env.Qbt.Host, CFG.Env.Qbt.Username, CFG.Env.Qbt.Password)
	}
	errs.PanicErr(err)
	Qbt = cli
	err = qbtRssEnable()
	errs.PanicErr(err)
}

func qbtRssEnable() error {
	cfg, err := Qbt.GetPreferences()
	if err != nil {
		return err
	}
	cfg.RSSAutoDownloadingEnabled = true
	cfg.RSSProcessingEnabled = true
	cfg.RSSMaxArticlesPerFeed = 50
	cfg.RSSRefreshInterval = 25
	err = Qbt.SetPreferences(cfg)
	return err
}

func Wait(ms int) {
	time.Sleep(time.Duration(ms) * time.Millisecond)
}
