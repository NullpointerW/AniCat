package subject

import (
	"log"
	"regexp"
	"strings"

	"github.com/NullpointerW/anicat/errs"

	qbt "github.com/NullpointerW/go-qbittorrent-apiv2"
)

func CaptureEpisNum(text string) (string, error) {
	for _, reg := range epi_regs {
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
	sep:="."
	sp := strings.Split(torr.Name, sep)
	extension := sp[len(sp)-1]
	extension=sep+extension
	basename := s.Name
	season := "S"
	episode := "E"

	epin, err := CaptureEpisNum(torr.Name)
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
	return rename, nil

}
