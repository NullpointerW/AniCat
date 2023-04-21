package crawl

import (
	"net/url"
	"strings"

	CFG "github.com/NullpointerW/mikanani/conf"
	"github.com/gocolly/colly"
	"github.com/gocolly/colly/proxy"
)

func SetProxy(c *colly.Collector) {
	if CFG.Proxy != nil {
		if p, err := proxy.RoundRobinProxySwitcher(
			CFG.Proxy...,
		); err == nil {
			c.SetProxyFunc(p)
		}
	}
}

func ConstructSearch(s string) (utoa string) {
	a := url.QueryEscape(strings.ReplaceAll(s, " ", "+"))
	utoa = strings.ReplaceAll(a, "%2B", "+")
	return
}
