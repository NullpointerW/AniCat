package rename

import (
	"fmt"
	"github.com/NullpointerW/anicat/log"
	"os"
	"testing"
)

func TestCaptureEpisNum(t *testing.T) {
	e, err := CaptureEpisNum("[桜都字幕组] 狼与辛香料 / Ookami to Koushinryou：Merchant Meets the Wise Wolf [18][1080p][繁体内嵌]")
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(e)
}

func TestTv(t *testing.T) {
	log.Init("text", "info", "2006-01-02T15:04:05", false, os.Stderr)
	tv, err := Tv("Ookami to Koushinryou", "01", "[UHA-WINGS][Bocchi the Rock!][12 END][x264 1080p][CHS].mp4")
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(tv)
}
func TestSubtitleFileLang(t *testing.T) {
	SubtitleFileLang("")
}
func TestSubtitleFileLang2(t *testing.T) {
	fmt.Println("a&" == "a\u0026")    
	fmt.Println("\n")
}
