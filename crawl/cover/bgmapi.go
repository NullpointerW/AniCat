package cover

import (
	"fmt"

	"net/http"

	"sync"

	"github.com/NullpointerW/mikanani/errs"
)

var (
	c *http.Client
	o sync.Once
)

func Bgmimage(subjid int, typ, filepath string) (err error) {
	o.Do(func() {
		c = &http.Client{}
	})
	url := fmt.Sprintf(BangumiImageUrl, subjid, typ)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return
	}
	req.Header.Set("User-Agent", "github.com/NullpointerW/mikanani")
	resp, err := c.Do(req)
	if err != nil {
		return
	}
	if sc := resp.StatusCode; sc != 200 {
		err = errs.Custom("bad request statusCdoe:%d", sc)
		return
	}
	return downloadfile(filepath, resp.Body)
}

// shorthand for Bgmimage(id,"large",filepath)
func TouchbgmCoverImg(id int, filepath string) error {
	return Bgmimage(id, Large, filepath)
}
