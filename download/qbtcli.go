package download

import (
	"log"
	"strconv"
	"strings"
	"time"

	CFG "github.com/NullpointerW/anicat/conf"
	"github.com/NullpointerW/anicat/errs"
	qbt "github.com/NullpointerW/go-qbittorrent-apiv2"
)

var Qbt *qbt.Client

func setProxyTyp(qbtCfg *qbt.Config) error {
	ptyp := strings.ToLower(CFG.Env.Qbt.Proxy.Type)
	switch ptyp {
	case "http":
		qbtCfg.ProxyType = qbt.Http
	case "httpa", "http-auth":
		qbtCfg.ProxyType = qbt.HttpA
	case "socks5":
		qbtCfg.ProxyType = qbt.Socks5
	case "socks5a", "socks5-auth":
		qbtCfg.ProxyType = qbt.Socks5A
	default:
		return errs.Custom("qbt:unknown proxy type %s", ptyp)
	}
	return nil
}

func setProxy(qbtCfg *qbt.Config) error {
	if CFG.Env.Qbt.Proxy.Type == "" {
		return nil
	}
	addr := strings.Split(CFG.Env.Qbt.Proxy.Addr, ":")
	if len(addr) != 2 {
		return errs.Custom("qbt:invalid proxy addr %s", CFG.Env.Qbt.Proxy.Addr)
	}
	host := addr[0]
	port, err := strconv.Atoi(addr[1])
	if err != nil {
		return err
	}
	err = setProxyTyp(qbtCfg)
	if err != nil {
		return err
	}
	qbtCfg.ProxyIP = host
	qbtCfg.ProxyPort = port
	qbtCfg.ProxyUsername = CFG.Env.Qbt.Proxy.Username
	qbtCfg.ProxyPassword = CFG.Env.Qbt.Proxy.Password
	qbtCfg.ProxyHostnameLookup = CFG.Env.Qbt.Proxy.Hslookup
	qbtCfg.ProxyPeerConnections = CFG.Env.Qbt.Proxy.Peer
	qbtCfg.ProxyTorrentsOnly = CFG.Env.Qbt.Proxy.TorrOnly
	return nil
}

func init() {
	var (
		cli *qbt.Client
		err error
	)
	if CFG.Env.Qbt.LocalConnect {
		cli, err = qbt.NewCli(CFG.Env.Qbt.Url)
	} else {
		cli, err = qbt.NewCli(CFG.Env.Qbt.Url, CFG.Env.Qbt.Username, CFG.Env.Qbt.Password)
	}
	errs.PanicErr(err)
	Qbt = cli
	cfg, err := Qbt.GetPreferences()
	errs.PanicErr(err)

	rssEnable(&cfg)
	err = setProxy(&cfg)
	if err != nil {
		log.Println(err)
	}else if  CFG.Env.Qbt.Proxy.Type != ""{
		log.Printf("qbt proxy addr:%s type:%s",CFG.Env.Qbt.Proxy.Addr, CFG.Env.Qbt.Proxy.Type)
	}
	errs.PanicErr(Qbt.SetPreferences(cfg))
	log.Println("qBittorrent connected")
	
}

func rssEnable(cfg *qbt.Config) {
	cfg.RSSAutoDownloadingEnabled = true
	cfg.RSSProcessingEnabled = true
	cfg.RSSMaxArticlesPerFeed = 50
	cfg.RSSRefreshInterval = 25
}

func Wait(ms int) {
	time.Sleep(time.Duration(ms) * time.Millisecond)
}
