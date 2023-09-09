package test

import (
	"testing"

	// "github.com/NullpointerW/anicat/downloader/detector"
	"github.com/NullpointerW/anicat/subject"
)

func TestCaptureEpisNum(t *testing.T) {
	n, err := subject.CaptureEpisNum("[orion origin] Tengoku Daimakyou [1-12] [1080p] [H265 AAC] [CHS＆JPN].mp4")
	if err != nil {
		t.Error(err)
		t.FailNow()
	}
	t.Log(n)
}
