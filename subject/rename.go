package subject

import (
	"log"
	"os"
	"regexp"
	"strings"
	"time"

	CFG "github.com/NullpointerW/anicat/conf"
	"github.com/NullpointerW/anicat/errs"
	"github.com/NullpointerW/anicat/util"

	DL "github.com/NullpointerW/anicat/download"
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

func RenameTV(s *Subject, torr qbt.Torrent) (string, error) {
	return renameTV(s, torr.Name)
}

func renameTV(s *Subject, fn string) (string, error) {
	sep := "."
	sp := strings.Split(fn, sep)
	extension := sp[len(sp)-1]
	extension = sep + extension
	basename := s.Name
	season := "S"
	episode := "E"

	epin, err := CaptureEpisNum(fn)
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
	log.Println("rename file", `"`, fn, `"`, "to", `"`, rename, `"`)
	return rename, nil
}

// func RenameMovie(s *Subject, torr qbt.Torrent) string {
// 	sep := "."
// 	sp := strings.Split(torr.Name, sep)
// 	extension := sp[len(sp)-1]
// 	extension = sep + extension
// 	basename := s.Name
// 	return basename + " " + extension
// }

func renameTorr(s *Subject, torr qbt.Torrent) error {
	epis := make(map[string]struct{})
	fs, err := DL.Qbt.Files(torr.Hash)
	if err != nil {
		return err
	}
	merr := errs.MultiErr{}
	for _, f := range fs {
		if fn := f.Name; util.IsVideofile(fn) {
			fn = util.FileSeparatorConv(fn)
			sep := strings.Split(fn, "/")
			fn = sep[len(sep)-1]
			rn, err := renameTV(s, fn)
			if err != nil {
				merr.Add(err)
				merr.Add(DL.Qbt.RenameFile(torr.Hash, f.Name, fn)) // mabye drop
				continue
			}
			se := util.TrimExtensionAndGetEpi(rn)
			if _, e := epis[se]; !e {
				epis[se] = struct{}{}
				merr.Add(DL.Qbt.RenameFile(torr.Hash, f.Name, rn))
				continue
			}
			if !CFG.Env.DropOnDumplicate {
				merr.Add(DL.Qbt.RenameFile(torr.Hash, f.Name, fn))
			}
		}
	}
	time.Sleep(1000 * time.Millisecond) // wati for qbt moving files
	merr.Add(os.RemoveAll(torr.ContentPath))
	return merr.Err()
}
