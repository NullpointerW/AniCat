package errs

import (
	"errors"
	"fmt"
)

var (
	ErrCrawlNotFound   = errors.New("content not crawled")
	ErrTorrnetNotFound = errors.New("torrnet not found")
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
