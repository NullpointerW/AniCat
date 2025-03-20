package debug

import (
	// "fmt"
	"fmt"
	"net/http"
	_ "net/http/pprof"

	"github.com/NullpointerW/anicat/conf"
)

func init() {
	if conf.IdeDebugging {
		go func() {
			err := http.ListenAndServe("127.0.0.1:6879", nil)
			if err != nil {
				panic(err)
			}
		}()
		fmt.Println("pprof enable,listen on localhost:6879/debug/pprof/")
	}

}
