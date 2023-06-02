package errs

import (
	"errors"
	"fmt"
)

var (
	ErrCrawlNotFound            = errors.New("content not crawled")
	ErrCoverDownLoadZeroSize    = errors.New("cover file download zero size after multiple attempts")
	ErrUnknownResCrawlLsType    = errors.New("unknown resource ls crawl type")
	ErrLsGroupUnavailableOnTorr = errors.New("command `ls group` is unavailableOnTorrent type")

	ErrTorrnetNotFound           = errors.New("there are no torrent files on qbt")
	ErrTorrnetOnSavePathNotFound = errors.New("there are no torrent files on savepath of a subject")
	ErrSubjectAlreadyExisted     = errors.New("subject already existed")
	ErrSubjectNotFound           = errors.New("subject not found")
	ErrBgmUrlNotFoundOnMikan     = errors.New("bgm url not found on mikanani")
	ErrUndefinedCrawlListType    = errors.New("undefined crawl list type")
	WarnRssRuleNotMatched        = errors.New("there is no any series mached,check your auto-download rule!")
	// command error
	ErrUnknownCommand           = errors.New("unknown command")
	ErrMissingCommandArgument   = errors.New("missing command argument")
	ErrAddCommandMissingHelping = errors.New("")
	WarnReservedCommand_lsg     = errors.New("the command 'lsg' is currently unavailable. use 'lsi' to  view the list of subtitle groups ")

	ErrItemAlreadyPushed = errors.New("item was already pushed")
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
		panic(err)
	}
}

func RequireNonErr(err error) bool {
	return err != nil
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