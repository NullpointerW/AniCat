package errs

import (
	"errors"
	"fmt"
)

var (
	ErrCrawlNotFound         = errors.New("content not crawled")
	ErrTorrnetNotFound       = errors.New("torrnet not found")
	ErrSubjectAlreadyExisted = errors.New("subject already existed")
	ErrSubjectNotFound       = errors.New("subject not found")
	ErrBgmUrlNotFoundOnMikan = errors.New("bgm url not found on mikanani")
)

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
