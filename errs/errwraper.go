package errs

import "errors"

var (
	ErrCrawlNotFound = errors.New("content not crawled")
)

func RequireNonErrr(err error) bool {
	return err != nil
}

func ErrTransfer(src error, dst *error) {
	*dst = src
}
