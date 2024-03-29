package subject

import (
	"fmt"
	"github.com/NullpointerW/anicat/log"
	eslog "github.com/NullpointerW/anicat/pkg/log"
	"os"
	"regexp"
	"strings"

	CFG "github.com/NullpointerW/anicat/conf"
	"github.com/NullpointerW/anicat/errs"
	util "github.com/NullpointerW/anicat/utils"

	DL "github.com/NullpointerW/anicat/downloader"
	qbt "github.com/NullpointerW/go-qbittorrent-apiv2"
)

func CaptureEpisNum(text string) (string, error) {
	for _, reg := range epi_regs {
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
	regexper := regexp.MustCompile(regSpecial)
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

func checkSingleVideo(torr qbt.Torrent) bool {
	return util.IsVideofile(torr.Name)
}

func renameTV(s *Subject, fn string) (string, error) {
	sep := "."
	sp := strings.Split(fn, sep)
	extension := sp[len(sp)-1]
	extension = sep + extension
	basename := s.FolderName
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
	log.Info(log.Struct{"from", fn, "to", rename}, "rename file")
	return rename, nil
}

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
				// even rename failed ,Remove from the subfolder
				merr.Add(DL.Qbt.RenameFile(torr.Hash, f.Name, fn)) // mabye drop?
				if CFG.Env.BgmiLog {
					CFG.BgmiLogger.Infof(eslog.Struct{"sid", s.SubjId, "name", s.Name}, "episode update(unnamed): %s", torr.Name)
				}
				continue
			}
			se := util.TrimExtensionAndGetEpi(rn)
			if _, e := epis[se]; !e {
				epis[se] = struct{}{}
				merr.Add(DL.Qbt.RenameFile(torr.Hash, f.Name, rn))
				continue
			}
			if !CFG.Env.DropOnDuplicate {
				merr.Add(DL.Qbt.RenameFile(torr.Hash, f.Name, fn))
			}

			if CFG.Env.BgmiLog {
				CFG.BgmiLogger.Infof(eslog.Struct{"sid", s.SubjId, "name", s.Name}, "episode update: %s", rn)
			}
			// feat 0.0.3b: support external subtitles
		} else if fn := f.Name; util.IsSubtitleFile(fn) {
			fn = util.FileSeparatorConv(fn)
			sep := strings.Split(fn, "/")
			fn = sep[len(sep)-1]
			// remove subtitleFile to outside
			merr.Add(DL.Qbt.RenameFile(torr.Hash, f.Name, fn))
		}
	}
	DL.Wait(1000) // wati for qbt moving files
	merr.Add(os.RemoveAll(torr.ContentPath))
	return merr.Err()
}

func renameSubRssTorr(s *Subject, torr qbt.Torrent) (videoRnOk bool, rename string, err error) {
	fs, err := DL.Qbt.Files(torr.Hash)
	if err != nil {
		return false, "", err
	}
	merr := errs.MultiErr{}
	var subFsList []string

	for _, f := range fs {
		if fn := f.Name; util.IsVideofile(fn) {
			fn = util.FileSeparatorConv(fn)
			sep := strings.Split(fn, "/")
			fn = sep[len(sep)-1]
			rn, err := renameTV(s, fn)
			if err != nil {
				merr.Add(err)
				merr.Add(DL.Qbt.RenameFile(torr.Hash, f.Name, fn))
				continue
			}
			rename = rn
			se := util.TrimExtensionAndGetEpi(rn)
			if th, e := s.Pushed[se]; e {
				dumpliErr := fmt.Errorf("%w: origin_name=%s,rename=%s", errs.ErrItemAlreadyPushed, torr.Name, rn)
				merr.Add(dumpliErr)
				if CFG.Env.DropOnDuplicate && th != torr.Hash {
					log.Warn(log.Struct{"torrfn", torr.Name, "torrHash", torr.Hash, "size", torr.Size}, "delete duplicateFile")
					merr.Add(DL.Qbt.DelTorrentsFs(torr.Hash))
					return false, rn, merr.Err()
				} else if !CFG.Env.DropOnDuplicate {
					merr.Add(DL.Qbt.RenameFile(torr.Hash, f.Name, fn))
				}
				// if we find some same episode files during the traversal of the current hash files
				// and enable `dropOnDuplicate`
				// then leave it on the current path ,and delete folder once for all at the end
			} else {
				err = DL.Qbt.RenameFile(torr.Hash, f.Name, rn)
				if err != nil {
					merr.Add(err)
					continue
				}
				s.Pushed[se] = torr.Hash
				videoRnOk = true
			}
		} else if fn := f.Name; util.IsSubtitleFile(fn) {
			subFsList = append(subFsList, fn)
		}
	}
	// subFiles rename process
	for _, fullFn := range subFsList {
		fn := fullFn
		fn = util.FileSeparatorConv(fn)
		sep := strings.Split(fn, "/")
		fn = sep[len(sep)-1]
		subrn := fn
		sublang := renameSubtitleFile(fn)
		if sublang != "" && rename != "" {
			seps := strings.Split(rename, ".")
			seps = seps[:len(seps)-1]
			extSep := strings.Split(fn, ".")
			ext := extSep[len(extSep)-1]
			subrn = strings.Join(seps, "") + "-" + sublang + "." + ext
			log.Info(log.Struct{"from", fullFn, "to", subrn}, "rename subtitleFile")
		}
		// remove subtitleFile to outside
		merr.Add(DL.Qbt.RenameFile(torr.Hash, fullFn, subrn))
	}
	DL.Wait(1000) // wati for qbt moving files
	merr.Add(os.RemoveAll(torr.ContentPath))
	return videoRnOk, rename, merr.Err()
}

func renameSubtitleFile(fn string) string {
	reg, _ := regexp.Compile(CHSSubStationReg)
	ok := reg.MatchString(fn)
	if ok {
		return "chs"
	}
	reg, _ = regexp.Compile(CHTSubStationReg)
	ok = reg.MatchString(fn)
	if ok {
		return "cht"
	}
	return ""
}
