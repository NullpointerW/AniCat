package cover

const (
	DouBancoverSearchUrl = `https://movie.douban.com/j/subject_suggest?q=%s`
	DouBancoverXpathExp  = `/html/body/div[@id='wrapper']/div[@id='content']/div[@class='grid-16-8 clearfix']/div[@class='article']/ul[@class='poster-col3 clearfix']/li[1]/div[@class='cover']/a/img/@src`
	BangumiImageUrl      = "http://api.bgm.tv/v0/subjects/%d/image?type=%s"
)

const (
	Small  = "small"
	Grid   = "grid"
	Large  = "large"
	Medium = "medium"
	Common = "common"
)
