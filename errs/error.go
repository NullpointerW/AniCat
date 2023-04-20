package errs

import (
	"errors"
	"fmt"
)

var (
	ErrCrawlNotFound  = errors.New("content not crawled")
	
)

func Custom(format string, a ...any) error {
	return fmt.Errorf(format, a...)
}

func RequireNonErr(err error) bool {
	return err != nil
}

func ErrTransfer(src error, dst *error) {
	*dst = src
}
