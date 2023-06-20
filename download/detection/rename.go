package detection

import (

	"regexp"

	"github.com/NullpointerW/anicat/errs"
	"github.com/NullpointerW/anicat/subject"
	qbt "github.com/NullpointerW/go-qbittorrent-apiv2"
)

const (
	reg0 = `]\[(\d{2})[vV]` // [02v1]
	reg1 = `\[(\d+)\]`         // [02]
	reg2 = `\b-\s*(\d+)`       // - 02

)

var regs = []string{reg0, reg1, reg2}

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

func Rename(s *subject.Subject, torr qbt.Torrent) error {
	basename := s.Name
	episode := "E"
	n, err := CaptureEpisNum(reg0)
	if err != nil {
		return err
	}
	episode += n
	basename += " " + episode
	return nil
}
