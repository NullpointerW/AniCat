package torrnet

import (
	"fmt"
	"testing"
)

func TestDL(t *testing.T) {
	h, err := Add("magnet:?xt=urn:btih:3522edcc5e979347bf1bc6a99cf12c15b5e66170&tr=http://open.acgtracker.com:1096/announce", "./dl", "subj333")
	if err != nil {
		t.Error(err)
		t.FailNow()
	}
	fmt.Println(h)
}