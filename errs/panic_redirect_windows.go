package errs

import (
	"golang.org/x/sys/windows"
	"os"
	"runtime"
)

var _g *os.File

func PanicRedirect(file *os.File) error {
	_g = file
	err := windows.SetStdHandle(windows.STD_ERROR_HANDLE, windows.Handle(file.Fd()))
	if err != nil {
		return err
	}
	// 内存回收前关闭文件描述符
	runtime.SetFinalizer(_g, func(fd *os.File) {
		_ = fd.Close()
	})
	return nil
}
