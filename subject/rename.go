package subject

import (
	"log"
	"regexp"
	"strings"

	
	"github.com/NullpointerW/anicat/errs"

	qbt "github.com/NullpointerW/go-qbittorrent-apiv2"
)

const (
	reg0v_e = `]\[(\d{2})[vV]` // [02v1]
	reg1_e = `\[(\d+)\]`      // [02]
	reg2_e = `\b-\s*(\d+)`    // - 02

)

var regs = []string{reg0v_e, reg1_e, reg2_e}

func CaptureEpisNum(text string) (string, error) {
	for _, reg := range regs {
		// fmt.Println(reg)
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
	return "", errs.Custom("%w:%s", errs.ErrCannotCaptureEpisNum, text)
}

func Rename(s *Subject, torr qbt.Torrent) (string, error) {
	sp := strings.Split(torr.Name, ".")
	extension := sp[len(sp)-1]
	basename := s.Name
	season := "S"
	episode := "E"

	epin, err := CaptureEpisNum(reg0)
	if err != nil {
		return "", err
	}
	sean := s.Season
	r := []rune(sean)
	if len(r) == 1 {
		sean = "0" + sean
	}
	season += sean
	episode += epin
	rename := basename + " " + season + episode + extension
	log.Println("rename file", `"`, torr.Name, `"`, "to", `"`, rename, `"`)
	return rename,nil

}
