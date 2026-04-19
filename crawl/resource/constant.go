package resource

import sel "github.com/NullpointerW/anicat/crawl/selector"

type LsTyp int

const (
	Ls      = 1
	LSGroup = 2
)

func (t LsTyp) String() string {
	if t == Ls {
		return "ls"
	}
	return "ls group"
}

const resourcesBaseUrl = `https://mikanime.tv`

// fallback XPaths used when selectors.yaml is not loaded
const (
	MikanRssLiXpathDefault = `/html/body[@class='main']/div[@id='sk-container']/div[@class='central-container']/ul[@class='list-inline an-ul']/li`
	BgmXpathExpDefault     = `/html/body[@class='main']/div[@id='sk-container']/div[@class='pull-left leftbar-container']/p[@class='bangumi-info'][last()]/a/@href`
)

func mikanXpath(key string, fallback string) string {
	if v := sel.Get("mikan", key); v != "" {
		return v
	}
	return fallback
}

var ResourceAPIs = map[string]string{
	"search": "/Home/Search?searchstr=",
}
