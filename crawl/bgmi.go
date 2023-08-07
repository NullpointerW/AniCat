package crawl

import (
	"fmt"
	"net/http"
	"time"

	"sync"

	"github.com/NullpointerW/anicat/errs"
)

var (
	c        *http.Client
	o        sync.Once
	BgmiRoot = "http://api.bgm.tv/"
)

type BgmiSubjIntro struct {
	Id         int               `json:"id"`
	Url        string            `json:"url"`
	Type       int               `json:"type"`
	Name       string            `json:"name"`
	NameCN     string            `json:"name_cn"`
	Summary    string            `json:"summary"`
	AirDate    string            `json:"air_date"`
	AirWeekday int               `json:"air_weekday"`
	Images     map[string]string `json:"images"`
}

func BgmiRequest(req *http.Request) (*http.Response, error) {
	o.Do(func() {
		c = &http.Client{
			Timeout: 10 * time.Second,
		}
	})
	req.Header.Set("User-Agent", "github.com/NullpointerW/anicat")
	resp, err := c.Do(req)
	if err != nil {
		return nil, err
	}
	if sc := resp.StatusCode; sc != 200 {
		err = fmt.Errorf("%w: bad request statusCdoe:%d", errs.ErrBgmTVApiPrefix, sc)
		return nil, err
	}
	return resp, nil
}
