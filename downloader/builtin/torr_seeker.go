package builtin

import (
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/anacrolix/torrent"
	"github.com/anacrolix/torrent/metainfo"
)

type TorrentSeeker interface {
	Seek(n string) (*torrent.TorrentSpec, error)
}

type HttpUrlSeeker struct {
	*http.Client
}

func (hr *HttpUrlSeeker) Seek(n string) (*torrent.TorrentSpec, error) {
	p, err := url.Parse(n)
	if scheme := strings.ToLower(p.Scheme); err != nil || !(scheme == "http" || scheme == "https") {
		return nil, fmt.Errorf("invalid url: %s", n)
	}
	resp, err := hr.Get(n)
	if err != nil {
		return nil, err
	}
	mf, err := metainfo.Load(resp.Body)
	if err != nil {
		return nil, err
	}
	ts := torrent.TorrentSpecFromMetaInfo(mf)
	return ts, nil
}
func NewHttpSeeker() *HttpUrlSeeker {
	transport := &http.Transport{
		MaxIdleConns:        10,
		MaxIdleConnsPerHost: 10,
		IdleConnTimeout:     90 * time.Second,
		MaxConnsPerHost:     10,
	}
	return &HttpUrlSeeker{
		&http.Client{Transport: transport},
	}

}
