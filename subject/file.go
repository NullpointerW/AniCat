package subject

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
	"time"

	CFG "github.com/NullpointerW/anicat/conf"
	DL "github.com/NullpointerW/anicat/downloader"
	"github.com/NullpointerW/anicat/downloader/rss"
	TORR "github.com/NullpointerW/anicat/downloader/torrent"
	"github.com/NullpointerW/anicat/errs"
	"github.com/NullpointerW/anicat/log"
	util "github.com/NullpointerW/anicat/utils"
)

var HOME = CFG.Env.SubjPath

func Scan() {
	builtin := CFG.Env.BuiltinDownloader
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
							if builtin != s.BuiltinDownload {
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
		os.Exit(1)
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
	if s.BuiltinDownload {
		return nil
	}
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
	return strings.TrimSuffix(strings.TrimSuffix(n, "\\"), "/")
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

var errXmlDocNotMatch = errors.New("modify tvshow.nfo: cannot found <title></title>")

func JellyfinMetaDataHelper(dp, name string, exited chan struct{}) {
	path := dp + "/" + "tvshow.nfo"
	for {
		select {
		case <-exited:
			return
		default:
			_, err := os.Stat(dp + "/" + "tvshow.nfo")
			if !os.IsNotExist(err) {
				err = InitTvNfo(path, name)
				if err != nil {
					log.Errorf(log.Struct{"err", err}, "JellyfinMetaDataHelper: modify %s failed", path)
					if err == errXmlDocNotMatch {
						goto sleep
					}
				} else {
					log.Info(nil, "JellyfinMetaData is ok")
				}
				return
			}
		sleep:
			time.Sleep(20 * time.Second)
		}
	}
}

func InitTvNfo(p, t string) error {
	xmlFile, err := os.OpenFile(p, os.O_RDWR, 0777)
	if err != nil {
		return err
	}
	defer xmlFile.Close()
	byteValue, _ := io.ReadAll(xmlFile)
	xmldata := string(byteValue)
	const doc = `<title>%s</title>`
	exp := fmt.Sprintf(doc, "(.*?)")
	re := regexp.MustCompile(exp)
	match := re.FindStringSubmatch(xmldata)
	newTitle, oldTitle := fmt.Sprintf(doc, t), ""
	if len(match) > 1 {
		if match[1] == t {
			return nil
		}
		oldTitle = fmt.Sprintf(doc, match[1])
	} else {
		return errXmlDocNotMatch
	}
	newXmlData := strings.ReplaceAll(xmldata, oldTitle, newTitle)
	err = xmlFile.Truncate(0)
	if err != nil {
		return err
	}
	_, err = xmlFile.Seek(0, 0)
	if err != nil {
		return err
	}
	_, err = xmlFile.WriteString(newXmlData)
	return err
}
