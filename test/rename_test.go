package test

import (
	"testing"

	// "github.com/NullpointerW/anicat/download/detection"
	"github.com/NullpointerW/anicat/subject"
)

func TestCaptureEpisNum(t *testing.T) {
	n, err := subject.CaptureEpisNum("[ANi] 我內心的糟糕念頭 - 09 [1080P][Baha][WEB-DL][AAC AVC][CHT].mp4	")
	if err != nil {
		t.Error(err)
		t.FailNow()
	}
	t.Log(n)
}
