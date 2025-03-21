package errs

import (
	"errors"
	"fmt"
)

var (
	ErrCrawlNotFound            = errors.New("crawl notFound")
	ErrCoverDownLoadZeroSize    = errors.New("coverFile downloader zeroSize after multiple attempts")
	ErrUnknownResCrawlLsType    = errors.New("unknownResource ls crawl type")
	ErrLsGroupUnavailableOnTorr = errors.New("command `ls group` is unavailable On TorrentType")

	ErrBgmTVApiPrefix = errors.New("bgmTV api")

	ErrTorrentNotFound           = errors.New("there are no torrentFiles on qbt")
	ErrTorrentOnSavePathNotFound = errors.New("there are no torrentFiles on savepath of a subject")
	ErrQbtDataNotFound           = errors.New("no data found from the apiRequest (should be found)")
	ErrSubjectAlreadyExisted     = errors.New("subject already existed")
	ErrSubjectNotFound           = errors.New("subject notFound")
	ErrBgmUrlNotFoundOnMikan     = errors.New("bgm url notFound on anicat")
	ErrUndefinedCrawlListType    = errors.New("undefined crawlList type")
	WarnRssRuleNotMatched        = errors.New("there is no any series mached,check your auto-downloader rule")

	ErrCannotCaptureEpisNum = errors.New("can not capture episodeNum from text")
	ErrItemAlreadyPushed    = errors.New("item was already pushed")
	ErrNoLinkFoundOnRssFeed = errors.New("can not found link on feed")
	ErrConnHajcked          = errors.New("conn hajacked")
)

type MultiErr struct {
	errs []error
}

func (me *MultiErr) Add(e error) {
	if e != nil {
		me.errs = append(me.errs, e)
	}
}

func (me *MultiErr) Err() error {
	if len(me.errs) == 0 {
		return nil
	}
	var errstr string
	for _, e := range me.errs {
		errstr += e.Error() + "\n"
	}
	return errors.New(errstr)
}

func Custom(format string, a ...any) error {
	return fmt.Errorf(format, a...)
}

func PanicErr(err error, ef func()) {
	if err != nil {
		ef()
		panic(err)
	}
}

func RequireNonErr(err error) bool {
	return err == nil
}

func ErrTransfer(src error, dst *error) {
	*dst = src
}

type Thrower func() error
type ErrWrapper struct {
	e error
}

func (wp *ErrWrapper) Handle(t Thrower) {
	if wp.e == nil {
		wp.e = t()
	}
}

func (wp *ErrWrapper) Error() error {
	return wp.e
}

func (wp *ErrWrapper) Reset() {
	wp.e = nil
}
