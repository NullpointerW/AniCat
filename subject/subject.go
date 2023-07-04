package subject

import (
	"context"
	"fmt"
	"log"
	"regexp"
	"strconv"
	"strings"
	"time"

	CC "github.com/NullpointerW/anicat/crawl/cover"
	IC "github.com/NullpointerW/anicat/crawl/information"
	RC "github.com/NullpointerW/anicat/crawl/resource"
	DL "github.com/NullpointerW/anicat/download"
	"github.com/NullpointerW/anicat/download/rss"
	"github.com/NullpointerW/anicat/download/torrent"
	"github.com/NullpointerW/anicat/errs"
	util "github.com/NullpointerW/anicat/utils"
	qbt "github.com/NullpointerW/go-qbittorrent-apiv2"
)

// Subject as basic obj of each bgmi
// programm will load it from OS file and manage thme
// while some of fileds have updating it will be refresh to OS file
type Subject struct {
	SubjId      int         `json:"subjId"`
	FolderName  string      `json:"folderName"` // source from tmdb
	Name        string      `json:"name"`
	Path        string      `json:"path"`
	Finished    bool        `json:"finished"`
	Episode     int         `json:"episode"`
	ResourceTyp ResourceTyp `json:"resourceTyp"`
	ResourceUrl string      `json:"resourceUrl"`
	Typ         BgmiTyp     `json:"typ"`
	FolderTime  string      `json:"folderTime"` // source from tmdb
	StartTime   string      `json:"startTime"`
	EndTime     string      `json:"endTime"`
	Alias       string      `json:"alias"`
	Season      string      `json:"season"`
	// used while `ResourceTyp` is `Torrent`
	TorrentHash string `json:"torrentHash"`
	// manager use ctxcancel func to Exit gorountine running the current subject.
	// when delete a subject manager,should run the cancelfunc and if a gorountine is runing
	// for this subject it will Exit.
	// Context is hold by subject-running gorountine
	// while subject-running gorountine Exit actively func should be called
	Exit context.CancelFunc `json:"-"`
	// before detection-gorountine push to subject,Check if this channel is closed.
	// before exit Exited channel should  be closed
	Exited chan struct{} `json:"-"`
	// While detection-gorountine detected that the resource download of the subject is completed
	// it will send downLoad message to subject-running gorountine
	// received and push to terminal
	PushChan chan qbt.Torrent `json:"-"`
	// The anime series of this project has already ended and all episodes have been downloaded.
	// while init,if this flag is true then there is noneed to start a gorountine to run it
	// exit actively flag should  be set to true
	Terminate bool `json:"terminate"`
	// a Set store all pushed renamed episodes,avoid duplicate push.
	// content will like be `S01E01,S01E05...`
	Pushed      map[string]string   `json:"pushed"`
	RssTorrents map[string]struct{} `json:"rssTorrents"`
}

type Extra struct {
	TorrOption struct {
		Index int
	}
	RssOption struct {
		SubtitleGroup  string
		MustContain    string
		MustNotContain string
		UseRegex       bool
	}
}

func (ex *Extra) NoArgs() bool {
	opt := ex.RssOption
	return opt.SubtitleGroup == "" && opt.MustContain == "" && opt.MustNotContain == ""
}

// The tag used when adding a torrent with qbt
// can be used to monitor the download status of resources
// related to this subject file.
func (s *Subject) QbtTag() string {
	return fmt.Sprintf(QbtTag, s.SubjId)
}

func (s *Subject) QbtCateg() string {
	return s.QbtTag()
}

func (s *Subject) RssPath() string {
	return s.QbtTag()
}

func CreateSubject(n string, ext *Extra) error {
	subject := new(Subject)

	// for testing
	util.Debugf("%#+v", *subject)

	bgmurl, err := solveResource(n, subject, ext)
	if err != nil {
		return err
	}

	var tips map[string]string

	if bgmurl != "" {
		tips, err = IC.DoScrape(bgmurl)
	} else {
		tips, err = IC.Scrape(n)
	}
	if err != nil {
		return err
	}

	sid, _ := strconv.Atoi(tips[IC.SubjId])
	if Manager.Get(sid) != nil {
		return errs.Custom("%w:sid:%d", errs.ErrSubjectAlreadyExisted, sid)
	}
	subject.SubjId = sid
	err = subject.Loadfileds(tips)
	if err != nil {
		return err
	}
	GetSeason(subject)
	subject.trimName()
	err = initFolder(subject)
	if err != nil {
		return err
	}
	lastS, err := FindLastSeason(subject.Path)
	if err != nil {
		return err
	}
	curr, err := strconv.Atoi(subject.Season)
	if err != nil {
		return err
	}

	if curr > lastS {
		cp := subject.Path + "/" + CoverFN
		err = CC.TouchbgmCoverImg(sid, cp)
		if err != nil {
			log.Println(err)
			err = CC.DOUBANCoverScraper.Scrape(cp, n)
			if err != nil {
				retry := 0
				for err == errs.ErrCoverDownLoadZeroSize {
					retry++
					if retry >= 3 {
						return err
					}
					time.Sleep(500 * time.Millisecond)
					err = CC.DOUBANCoverScraper.Scrape(cp, n)
				}
				if err != nil {
					return err
				}
			}
		}
	}

	err = download(subject, ext)
	if err != nil {
		return err
	}

	// create Info-Json after init completed
	subject.writeJson()

	subject.runtimeInit(false)

	if subject.ResourceTyp == RSS {
		DL.Wait(1500) // wait for qbt
		a, err := rss.GetMatchedArts(subject.RssPath())
		if err == nil && len(a) == 0 {
			return errs.WarnRssRuleNotMatched
		}
	}

	return nil
}

func (subj *Subject) Loadfileds(tips map[string]string) error {
	subj.Name = tips[IC.SubjName]
	if _, e := tips[IC.SubjStartTime]; e {
		subj.Typ = TV
	} else {
		subj.Typ = MOVIE
	}
	subj.Episode, _ = strconv.Atoi(tips[IC.SubjEpisode])
	if subj.Typ == TV {
		subj.StartTime = tips[IC.SubjStartTime]
		if et, e := tips[IC.SubjectEndTime]; e {
			n := time.Now()
			eti, err := util.ParseTime(et, util.YMDParseLayout)
			if err != nil {
				return err
			}
			subj.EndTime = et
			subj.Finished = n.After(eti) || n.Equal(eti)
		}
	} else { // if movie finished
		subj.StartTime = tips[IC.SubjMoveStartTime]
		subj.Finished = true
	}
	subj.Alias = tips[IC.Alias]
	// source from tmdb
	var tmdbTyp = IC.TMDB_TYP_TV
	if subj.Typ == MOVIE {
		tmdbTyp = IC.TMDB_TYP_MOVIE
	}
	var err error
	subj.FolderName, subj.FolderTime, err = IC.FloderSearch(tmdbTyp, subj.Name)
	if err != nil {
		if err == errs.ErrCrawlNotFound {
			for _, n := range strings.Split(subj.Alias, "|") {
				subj.FolderName, subj.FolderTime, err = IC.FloderSearch(tmdbTyp, n)
				if err == nil {
					log.Printf("%#+v",subj)
					return nil
				} else if err != errs.ErrCrawlNotFound {
					return err
				}
			}
		}
		return err
	}
	
	return nil
}

func (s *Subject) FetchInfo() error {
	tips, err := IC.BgmTVInfoScrape(s.SubjId)
	if err != nil {
		return nil
	}
	wrap := errs.ErrWrapper{}
	wrap.Handle(func() error {
		return s.Loadfileds(tips)
	})
	wrap.Handle(func() error {
		return s.writeJson()
	})
	return wrap.Error()
}

func solveResource(n string, subj *Subject, ext *Extra) (string, error) {
	opt := RC.Option{}
	if ext != nil {
		opt.Group = ext.RssOption.SubtitleGroup
		util.Debugln("SubtitleGroup", opt.Group)
		opt.Index = ext.TorrOption.Index
	}
	u, bgm, isrss, err := RC.Scrape(n, opt)
	if err != nil {
		return "", err
	}
	subj.ResourceUrl = u
	if isrss {
		subj.ResourceTyp = RSS
	} else {
		subj.ResourceTyp = Torrent
	}
	log.Println(u, bgm, isrss)
	return bgm, nil
}

func download(subj *Subject, ext *Extra) error {
	if subj.ResourceTyp == Torrent {
		h, err := torrent.Add(subj.ResourceUrl, subj.Path, subj.QbtTag())
		subj.TorrentHash = h
		return err
	} else if (ext == nil || ext.NoArgs()) && subj.Finished {
		it, err := rss.AddAndGetItems(subj.ResourceUrl, subj.RssPath())
		util.Debugln("rss path ", subj.RssPath())
		if err != nil {
			return err
		}
		for _, a := range it.Articles {
			log.Println(a.Description)
			desc := a.Description
			for _, reg := range coll_regs {
				re, err := regexp.Compile(reg)
				if err != nil {
					return err
				}
				if re.MatchString(desc) {
					log.Printf("%d:%s  matched collection %s \n", subj.SubjId, subj.Name, desc)
					err = rss.RmRss(subj.RssPath())
					if err != nil {
						return err
					}
					subj.ResourceTyp = Torrent
					h, err := torrent.Add(a.TorrentURL, subj.Path, subj.QbtTag())
					subj.TorrentHash = h
					return err
				}
			}
		}
		log.Printf("%d:%s not matched any collection \n", subj.SubjId, subj.Name)

		err = torrent.AddCategroy(subj.QbtCateg())
		if err != nil {
			return err
		}
		err = rss.SetAutoDLRule(subj.ResourceUrl, subj.QbtCateg(), subj.Path, subj.RssPath())
		return err
	} else {

		err := torrent.AddCategroy(subj.QbtCateg())
		if err != nil {
			return err
		}
		r := qbt.AutoDLRule{
			Enabled:          true,
			AffectedFeeds:    []string{subj.ResourceUrl},
			SavePath:         subj.Path,
			AssignedCategory: subj.QbtCateg(),
		}
		if ext != nil {
			r.UseRegex = ext.RssOption.UseRegex
			r.MustContain = ext.RssOption.MustContain
			r.MustNotContain = ext.RssOption.MustNotContain
		}
		err = rss.Download(r, subj.RssPath())
		return err
	}
}

func GetSeason(s *Subject) {
	var ns []string
	ns = append(ns, s.Name)
	as := strings.Split(s.Alias, "|")
	ns = append(ns, as...)
	for _, n := range ns {
		regexper := regexp.MustCompile(zhreg)
		match := regexper.FindStringSubmatch(n)
		if len(match) > 1 {
			m := match[1]
			if iszh := util.CheckZhCn(m); iszh {
				m = util.ConvertZhCnNumbToa(m)
			}
			s.Season = fmt.Sprintf("%02s", m)
			return
		}
		for _, rg := range sregs {
			regexper := regexp.MustCompile(rg)
			match := regexper.FindStringSubmatch(n)
			if len(match) > 1 {
				s.Season = fmt.Sprintf("%02s", match[1])
				return
			}
		}
	}
	s.Season = "01"
}

func (s *Subject) trimName() {
	s.Name = strings.ReplaceAll(util.FileSeparatorConv(s.Name), "/", " ")
}
