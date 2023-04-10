package errs

func RequireNonErrr(err error) bool {
	return err != nil
}

func ErrTransfer(src error, dst *error) {
	*dst = src
}
