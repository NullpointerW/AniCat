package download

import (
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
}
