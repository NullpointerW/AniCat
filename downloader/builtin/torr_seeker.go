package builtin

import (
	"fmt"
	"io"
	"net/http"
	"net/url"
	"time"
)

type TorrentSeeker interface {
	Seek(n string) (io.Reader, error)
}

type HttpUrlReader struct {
	*http.Client
}

func (hr *HttpUrlReader) Seek(n string) (io.Reader, error) {
	p, err := url.Parse(n)
	if err != nil || !(p.Scheme == "http" || p.Scheme == "https") {
		return nil, fmt.Errorf("invalid url: %s", n)
	}
	resp, err := hr.Get(n)
	if err != nil {
		return nil, err
	}
	return resp.Body, nil
}
func NewHttpReader() *HttpUrlReader {
	transport := &http.Transport{
		MaxIdleConns:        10,
		MaxIdleConnsPerHost: 10,
		IdleConnTimeout:     90 * time.Second,
		MaxConnsPerHost:     10,
	}
	return &HttpUrlReader{
		&http.Client{Transport: transport},
	}

}
