package errs

import (
	"errors"
	"fmt"
	"github.com/NullpointerW/anicat/log"
	"runtime"
)

var (
	ErrCrawlNotFound            = errors.New("crawl not found")
	ErrCoverDownLoadZeroSize    = errors.New("cover file download zero size after multiple attempts")
	ErrUnknownResCrawlLsType    = errors.New("unknown resource ls crawl type")
	ErrLsGroupUnavailableOnTorr = errors.New("command `ls group` is unavailableOnTorrent type")

	ErrBgmTVApiPrefix = errors.New("bgmTV api")

	ErrTorrentNotFound           = errors.New("there are no torrent files on qbt")
	ErrTorrentOnSavePathNotFound = errors.New("there are no torrent files on savepath of a subject")
	ErrQbtDataNotFound           = errors.New("no data found from the api request (should be found)")
	ErrSubjectAlreadyExisted     = errors.New("subject already existed")
	ErrSubjectNotFound           = errors.New("subject not found")
	ErrBgmUrlNotFoundOnMikan     = errors.New("bgm url not found on anicat")
	ErrUndefinedCrawlListType    = errors.New("undefined crawl list type")
	WarnRssRuleNotMatched        = errors.New("there is no any series mached,check your auto-download rule!")
	// command error
	ErrUnknownCommand           = errors.New("unknown command")
	ErrMissingCommandArgument   = errors.New("missing command argument")
	ErrAddCommandMissingHelping = errors.New("")
	WarnReservedCommand_lsg     = errors.New("the command 'lsg' is currently unavailable. use 'lsi' to  view the list of subtitle groups ")

	ErrCannotCaptureEpisNum = errors.New("can not capture episode num from text")

	ErrItemAlreadyPushed = errors.New("item was already pushed")

	ErrNoLinkFoundOnRssFeed = errors.New("can not found link on feed")
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

func PanicErr(err error) {
	if err != nil {
		if runtime.GOOS == "windows" {
			log.Error(log.Struct{"err", err}, "PANIC! process crashed")
		}
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
