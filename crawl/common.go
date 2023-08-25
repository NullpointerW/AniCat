package crawl

import (
	"github.com/NullpointerW/anicat/log"
	"io"
	"net/url"
	"os"
	"strings"

	CFG "github.com/NullpointerW/anicat/conf"
	"github.com/NullpointerW/anicat/errs"
	"github.com/gocolly/colly"
	"github.com/gocolly/colly/proxy"
)

func NewCollector() *colly.Collector {
	c := colly.NewCollector()
	c.Limit(&colly.LimitRule{Parallelism: 1})
	SetProxy(c)
	return c
}

func SetProxy(c *colly.Collector) {
	if len(CFG.Env.Crawl.Proxies) != 0 {
		if p, err := proxy.RoundRobinProxySwitcher(
			CFG.Env.Crawl.Proxies...,
		); err == nil {
			c.SetProxyFunc(p)
		}
	}
}

func UrlEncode(s string) (utoa string) {
	a := url.QueryEscape(strings.ReplaceAll(s, " ", "+"))
	utoa = strings.ReplaceAll(a, "%2B", "+")
	return
}

func Downloadfile(filepath string, remote io.ReadCloser) error {
	f, err := os.OpenFile(filepath, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0777)
	if err != nil {
		return err
	}
	defer remote.Close()
	defer f.Close()
	wn, err := io.Copy(f, remote)
	log.Info(log.Struct{"size", wn}, "cover downloaded")
	if err != nil {
		return err
	}
	if wn == 0 {
		return errs.ErrCoverDownLoadZeroSize
	}
	return nil
}
