package infomation

const (
	infoBaseUrl      = `https://bgm.tv`
	infoPageXpathExp = `/html/body[@class='bangumi']/div[@id='wrapperNeue']/div[@id='main'][2]/div[@class='columns clearit']/div[@id='columnSearchB']/ul[@id='browserItemList']/li[1]/div[@class='inner']/h3/a[@class='l']/@href`
	infoXpathExp     = `/html/body[@class='bangumi']/div[@id='wrapperNeue']/div[@class='mainWrapper']/div[@class='columns clearit']/div[@id='columnSubjectHomeA']/div[@id='bangumiInfo']/div[@class='infobox']/ul[@id='infobox']/li`
)

// keys for info map Scrape from bgm.tv
const (
	SubjId         = "sid"
	SubjName       = "中文名"
	SubjEpisode    = "话数"
	SubjStartTime  = "放送开始"
	SubjectEndTime = "播放结束"
)

var InfoAPIs = map[string]string{
	"search":  "/subject_search/%s?cat=2",
	"subject": "/subject/%d",
}
