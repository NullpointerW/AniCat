package rename

import (
	"errors"
	"fmt"
	"github.com/NullpointerW/anicat/errs"
	"github.com/NullpointerW/anicat/log"
	"path/filepath"
	"regexp"
)

func CaptureEpisNum(text string) (string, error) {
	for _, reg := range epiRegs {
		regexper := regexp.MustCompile(reg)
		match := regexper.FindStringSubmatch(text)
		if len(match) > 1 {
			episNum := match[1]
			if len([]byte(episNum)) == 1 {
				return "0" + episNum, nil
			}
			return episNum, nil
		}
	}
	regexper := regexp.MustCompile(specialReg)
	matchs := regexper.FindAllStringSubmatch(text, -1)
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
	reg, _ := regexp.Compile(chsSubStationReg)
	ok := reg.MatchString(fn)
	if ok {
		return "chs"
	}
	reg, _ = regexp.Compile(chtSubStationReg)
	ok = reg.MatchString(fn)
	if ok {
		return "cht"
	}
	return ""
}
