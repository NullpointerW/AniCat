package rename

import (
	"errors"
	"fmt"
	"path/filepath"

	"github.com/NullpointerW/anicat/errs"
	"github.com/NullpointerW/anicat/llmparser"
	"github.com/NullpointerW/anicat/log"
)

func CaptureEpisNum(text string) (string, error) {
	for _, re := range compiledEpiRegs {
		match := re.FindStringSubmatch(text)
		if len(match) > 1 {
			episNum := match[1]
			if len([]byte(episNum)) == 1 {
				return "0" + episNum, nil
			}
			return episNum, nil
		}
	}
	matchs := compiledSpecialReg.FindAllStringSubmatch(text, -1)
	if matchs != nil {
		if l := len(matchs); l == 1 {
			episNum := matchs[0][1]
			if len([]byte(episNum)) == 1 {
				return "0" + episNum, nil
			}
			return episNum, nil
		}
		episNum := matchs[1][1]
		if len([]byte(episNum)) == 1 {
			return "0" + episNum, nil
		}
		return episNum, nil
	}
	if llmparser.Fallback != nil {
		return llmparser.Fallback.Parse(text)
	}
	return "", fmt.Errorf("%w:%s", errs.ErrCannotCaptureEpisNum, text)
}

func Tv(base, sean, fn string) (string, error) {
	extension := filepath.Ext(fn)
	if extension == "" {
		return "", errors.New("rename: " + fn + " is not a file")
	}
	basename := base
	season := "S"
	episode := "E"
	epin, err := CaptureEpisNum(fn)
	if err != nil {
		return "", err
	}
	r := []rune(sean)
	if len(r) == 1 {
		sean = "0" + sean
	}
	season += sean
	episode += epin
	rename := basename + " " + season + episode + extension
	log.Info(log.Struct{"from", fn, "to", rename}, "rename file")
	return rename, nil
}
func SubtitleFileLang(fn string) string {
	if compiledChsSub.MatchString(fn) {
		return "chs"
	}
	if compiledChtSub.MatchString(fn) {
		return "cht"
	}
	return ""
}
