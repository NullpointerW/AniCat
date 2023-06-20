package test

import (
	"testing"

	// "github.com/NullpointerW/anicat/download/detection"
	"github.com/NullpointerW/anicat/subject"
)

func TestCaptureEpisNum(t *testing.T) {
	n, err := subject.CaptureEpisNum("[星空字幕组][小林家的龙女仆S / Kobayashi-san Chi no Maid Dragon S][09v2][简体内嵌][1080P][WebRip][MP4] [复制磁连]	")
	if err != nil {
		t.Error(err)
		t.FailNow()
	}
	t.Log(n)
}
