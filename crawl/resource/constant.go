package resource

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

const (
	MikanRssLiXpath = `/html/body[@class='main']/
		div[@id='sk-container']/
		div[@class='central-container']/
		ul[@class='list-inline an-ul']/li`

	BgmXpathExp = `/html/body[@class='main']/div[@id='sk-container']/
		div[@class='pull-left leftbar-container']/
		p[@class='bangumi-info'][last()]/
		a/@href`
)

var ResourceAPIs = map[string]string{
	"search": "/Home/Search?searchstr=",
}
