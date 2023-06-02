package errs

import (
	"fmt"
	"testing"
)

func TestErrCustom(t *testing.T) {
	err := Custom("custom error:%s", "my custom err")
	fmt.Println(err.Error())

	err = Custom("%w torr hash:%s", ErrTorrnetNotFound, "3522edcc5e979347bf1bc6a99cf12c15b5e66170")
	fmt.Println(err.Error())
}
