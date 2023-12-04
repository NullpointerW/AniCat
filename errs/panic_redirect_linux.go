package errs

import (
	"os"
	"runtime"
	"syscall"
)

var _g *os.File

func PanicRedirect(file *os.File) error {
	_g = file
	if err := syscall.Dup2(int(file.Fd()), int(os.Stderr.Fd())); err != nil {
		fmt.Println(err)
		return err
	}
	// 内存回收前关闭文件描述符
	runtime.SetFinalizer(_g, func(fd *os.File) {
		_ = fd.Close()
	})
	return nil
}
