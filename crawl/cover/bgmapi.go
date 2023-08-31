package cover

import (
	"fmt"

	"net/http"

	CR "github.com/NullpointerW/anicat/crawl"
)

func Bgmimage(subjid int, typ, filepath string) (err error) {
	url := fmt.Sprintf(BangumiImageUrl, subjid, typ)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return
	}
	resp, err := CR.BgmiRequest(req)
	if err != nil {
		return
	}
	return CR.DownloadFile(filepath, resp.Body)
}

// shorthand for Bgmimage(id,"large",filepath)
func TouchbgmCoverImg(id int, filepath string) error {
	return Bgmimage(id, Large, filepath)
}
