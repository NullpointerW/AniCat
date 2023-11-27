package test

import (
	"github.com/NullpointerW/anicat/subject"
	"testing"
)

func TestCaptureEpisNum(t *testing.T) {
	n, err := subject.CaptureEpisNum("春原庄的管理人小姐 05.mkv")
	if err != nil {
		t.Error(err)
		t.FailNow()
	}
	t.Log(n)
}
