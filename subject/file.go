package subject

import (
	"encoding/json"
	"fmt"
	CFG "github.com/NullpointerW/anicat/conf"
	DL "github.com/NullpointerW/anicat/downloader"
	"github.com/NullpointerW/anicat/downloader/rss"
	TORR "github.com/NullpointerW/anicat/downloader/torrent"
	"github.com/NullpointerW/anicat/errs"
	"github.com/NullpointerW/anicat/log"
	util "github.com/NullpointerW/anicat/utils"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

var HOME = CFG.Env.SubjPath

func Scan() {
	home := trimPath(HOME)
	if fs, err := os.ReadDir(home); err == nil {
		for _, f := range fs {
			if f.IsDir() {
				log.Info(log.Struct{"path", util.FileSeparatorConv(home + string(os.PathSeparator) + f.Name())},
					"scan: found folder")
				fs, err := os.ReadDir(home + `/` + f.Name())
				if err != nil {
					log.Error(log.Struct{"err", err}, "scan: open folder failed")
					continue
				}
				for _, ff := range fs {
					isf := !ff.IsDir()
					if isf && util.IsJsonFile(ff.Name()) && strings.Contains(ff.Name(), "meta-data#") {
						if jsraw, err := os.ReadFile(home + `/` + f.Name() + `/` + ff.Name()); err == nil {
							var s Subject
							err := json.Unmarshal(jsraw, &s)
							if err != nil {
								log.Error(log.Struct{"err", err}, "scan: unmarshal json failed")
								continue
							}
							s.runtimeInit(true)
						} else {
							log.Error(log.Struct{"err", err}, "scan: open file failed")
						}
					}
				}
			}
		}
	} else {
		log.Error(log.Struct{"err", err}, "scan: open home folder failed")
	}
}

// initFolder Initialize the content library in OS file for the `subject` in the memory.
// Path can be used to monitor the downloader status of resources
// apply to RSS and Torrent type.
// If initialization is successful, write the path to Subject.Path.
func initFolder(subject *Subject) (err error) {
	var folderPath string

	folderPath = trimPath(HOME)
	sd, err := util.ParseShort02Time(strings.ReplaceAll(subject.FolderTime, " ", ""))
	if err != nil {
		return err
	}

	folderPath += "/" + subject.FolderName + " (" + sd + ")"

	err = os.MkdirAll(folderPath, 0777)
	if err != nil {
		return err
	}

	if ap, err := filepath.Abs(folderPath); err == nil {
		subject.Path = util.FileSeparatorConv(ap)
	}

	return
}

func RmFolder(s *Subject) error {
	fs, err := os.ReadDir(s.Path)
	if err != nil {
		return err
	}
	var c int
	for _, f := range fs {
		isf := !f.IsDir()
		if isf && util.IsJsonFile(f.Name()) && strings.Contains(f.Name(), "meta-data#") {
			c++
		}
	}
	if c > 1 {
		return os.Remove(s.Path + "/" + fmt.Sprintf("meta-data#%s.json", s.GetSeasonAndPart()))
	} else {
		return rmFolder(s)
	}
}

func rmFolder(s *Subject) error {
	return os.RemoveAll(s.Path)
}

func (s *Subject) writeJson() (err error) {
	b, _ := json.Marshal(*s)
	fldrp := s.Path
	jsfn := fmt.Sprintf(jsonfileName, s.GetSeasonAndPart())
	err = os.WriteFile(fldrp+"/"+jsfn, b, 0777)
	return err
}

func (s *Subject) RmRes() error {
	wrap := errs.ErrWrapper{}
	if s.ResourceTyp == RSS {
		wrap.Handle(func() error {
			return TORR.DelViaCateg(s.QbtCateg())
		})

		// There is a possibility that categories cannot be completely deleted,
		// so it needs to be repeated multiple times
		// to ensure they are truly removed
		rep := 2
		for i := 0; i < rep; i++ {
			wrap.Handle(func() error {
				categ := s.QbtTag()
				return DL.Qbt.RmCategoies(categ)
			})
		}

		wrap.Handle(func() error {
			err := rss.RmRss(s.RssPath())
			if err != nil && strings.Contains(err.Error(), "409") {
				log.Error(log.Struct{"err", err}, "rm rssSubscription failed")
				return nil
			}
			return err
		})
	} else {
		wrap.Handle(func() error {
			return TORR.Del(s.TorrentHash)
		})
		wrap.Handle(func() error {
			return TORR.DelTag(s.QbtTag())
		})
	}
	return wrap.Error()
}

func trimPath(n string) string {
	return strings.TrimSuffix(strings.TrimSuffix(CFG.Env.SubjPath, "\\"), "/")
}

func FindLastSeason(p string) (int, error) {
	fs, err := os.ReadDir(p)
	if err != nil {
		return -1, err
	}
	var max int
	for _, f := range fs {
		isf := !f.IsDir()
		if isf && util.IsJsonFile(f.Name()) && strings.Contains(f.Name(), "meta-data#") {
			s := strings.Split(f.Name(), ".")[0]
			s = strings.ReplaceAll(s, "meta-data#", "")
			s = s[0:2] // todo get part extra
			cmp, err := strconv.Atoi(s)
			if err != nil {
				return -1, err
			}
			if cmp > max {
				max = cmp
			}
		}
	}
	return max, nil
}
