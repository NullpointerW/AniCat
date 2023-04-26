package download

import (
	qbt "github.com/NullpointerW/go-qbittorrent-apiv2"
	CFG "github.com/NullpointerW/mikanani/conf"
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
	if err != nil {
		panic(err)
	}
	Qbt = cli
}
