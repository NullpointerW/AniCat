package subject

import (
	"context"
	"errors"
	"fmt"
	"log"
	"regexp"
	"strconv"
	"strings"
	"time"

	CC "github.com/NullpointerW/anicat/crawl/cover"
	IC "github.com/NullpointerW/anicat/crawl/information"
	RC "github.com/NullpointerW/anicat/crawl/resource"

	// DL "github.com/NullpointerW/anicat/download"
	CFG "github.com/NullpointerW/anicat/conf"
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
	OriginName  string      `json:"orginName"`
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
	return opt.MustContain == "" && opt.MustNotContain == ""
}

func BuildFilterPerlReg(vbs []string) string {
	var reg string
	const tmp = `(?=.*?%s)`
	if len(vbs) != 0 {
		reg += "(?i)"
		for _, ct := range vbs {
			vb := strings.ReplaceAll(ct, ",", "|")
			vb = "(" + vb + ")"
			reg += fmt.Sprintf(tmp, vb)
		}
		return reg
	} else {
		return ""
	}
}

func BuildFilterRegs(vbs []string) []string {
	if len(vbs) != 0 {
		regs := make([]string, 0, len(vbs))
		for _, ct := range vbs {
			vb := strings.ReplaceAll(ct, ",", "|")
			vb = "(?i)" + vb
			regs = append(regs, vb)
		}
		return regs
	} else {
		return nil
	}
}

func FilterWithRegs(s string, contains, exclusions []string) bool {
	var (
		containOk, exclusionOk bool
	)
	if len(contains) == 0 {
		containOk = true
	}
	if len(exclusions) == 0 {
		exclusionOk = true
	}
	if !containOk {
		containOks := make([]bool, 0, len(contains))
		for _, reg := range contains {
			var ok bool
			csreg, err := regexp.Compile(reg)
			if err != nil {
				log.Println(fmt.Errorf("golbal filter contains regexp error: %w", err))
				ok = true
			} else {
				ok = csreg.MatchString(s)
				util.Debugln(csreg.String(), ":", ok)
			}
			containOks = append(containOks, ok)
		}
		containOk = true
		for _, ok := range containOks {
			if !ok {
				containOk = false
				break
			}
		}
	}

	if !exclusionOk {
		exclusionOks := make([]bool, 0, len(exclusions))
		for _, reg := range exclusions {
			var ok bool
			clsreg, err := regexp.Compile(reg)
			if err != nil {
				log.Println(fmt.Errorf("golbal filter exclusions regexp error: %w", err))
				ok = true
			} else {
				ok = !clsreg.MatchString(s)
				util.Debugln(clsreg.String(), ":", ok)
			}
			exclusionOks = append(exclusionOks, ok)
		}
		exclusionOk = true
		for _, ok := range exclusionOks {
			if !ok {
				exclusionOk = false
				break
			}
		}
	}
	return containOk && exclusionOk
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

func CreateSubject(n string, ext *Extra) (int, error) {
	subject := new(Subject)

	// for testing
	util.Debugf("%#+v", *subject)

	bgmurl, err := solveResource(n, subject, ext)
	if err != nil {
		return 0, err
	}

	var tips map[string]string

	if bgmurl != "" {
		tips, err = IC.DoScrape(bgmurl)
	} else {
		tips, err = IC.Scrape(n)
	}
	if err != nil {
		return 0, err
	}

	sid, err := strconv.Atoi(tips[IC.SubjId])
	if err != nil {
		return 0, err
	}
	if Manager.Get(sid) != nil {
		return 0, fmt.Errorf("%w:sid:%d", errs.ErrSubjectAlreadyExisted, sid)
	}
	subject.SubjId = sid
	err = subject.Loadfileds(tips)
	if err != nil {
		return 0, err
	}
	GetSeason(subject)
	subject.trimName()
	err = initFolder(subject)
	if err != nil {
		return 0, err
	}
	lastS, err := FindLastSeason(subject.Path)
	if err != nil {
		return 0, err
	}
	curr, err := strconv.Atoi(subject.Season)
	if err != nil {
		return 0, err
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
						return 0, err
					}
					time.Sleep(500 * time.Millisecond)
					err = CC.DOUBANCoverScraper.Scrape(cp, n)
				}
				if err != nil {
					return 0, err
				}
			}
		}
	}

	err = download(subject, ext)
	if err != nil {
		return 0, err
	}

	// create Info-Json after init completed
	subject.writeJson()

	subject.runtimeInit(false)

	// if subject.ResourceTyp == RSS {
	// 	// DL.Wait(1500) // wait for qbt
	// 	m, err := DL.DoFetch(func() (recvd bool, err error) {
	// 		a, err := rss.GetMatchedArts(subject.RssPath())
	// 		if err != nil {
	// 			return false, err
	// 		}
	// 		return len(a) > 0, nil
	// 	}, 3000)
	// 	if err != nil {
	// 		log.Println(fmt.Errorf("check rss matched item error:%w subjid:%d", err, subject.SubjId))
	// 		return nil
	// 	}
	// 	if !m {
	// 		return errs.WarnRssRuleNotMatched
	// 	}
	// }
	log.Printf("create subj%d succeeded \n", sid)
	return sid, nil
}

func (subj *Subject) Loadfileds(tips map[string]string) error {
	subj.Name = tips[IC.SubjName]
	subj.OriginName = tips[IC.SubjOriginName]
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

	// fetch folder info,source from tmdb
	var tmdbTyp = IC.TMDB_TYP_TV
	if subj.Typ == MOVIE {
		tmdbTyp = IC.TMDB_TYP_MOVIE
	}
	var err error
	// First, attempt to search for the folder using the subject's Name
	subj.FolderName, subj.FolderTime, err = IC.FloderSearch(tmdbTyp, subj.Name)
	if err != nil {
		if errors.Is(err, errs.ErrCrawlNotFound) {
			// If the search failed because the folder was not found,
			// try again using the subject's OriginName
			subj.FolderName, subj.FolderTime, err = IC.FloderSearch(tmdbTyp, subj.OriginName)
		}
		if errors.Is(err, errs.ErrCrawlNotFound) {
			// If the search still cannot found, try using each alias from the subject's Alias field
			for _, n := range strings.Split(subj.Alias, "|") {
				subj.FolderName, subj.FolderTime, err = IC.FloderSearch(tmdbTyp, n)
				if err == nil {
					return nil
				} else if !errors.Is(err, errs.ErrCrawlNotFound) {
					// If the search failed with an error other than ErrCrawlNotFound, return the error
					return err
				}
			}
			// If all aliases failed, try removing the season number from the subject's Name and search again
			re := regexp.MustCompile(`第(.)季`)
			match := re.FindStringSubmatch(subj.Name)
			if len(match) > 1 {
				season := match[1]
				n := strings.ReplaceAll(subj.Name, fmt.Sprintf("第%s季", season), "")
				n = strings.TrimRight(n, " ")
				subj.FolderName, subj.FolderTime, err = IC.FloderSearch(tmdbTyp, n)
				return err
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
	util.Debugln("resource:", u, "bgmi url:", bgm, "is rss:", isrss)
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
		enaFl := CFG.Env.EnabledFilter()
		for _, a := range it.Articles {
			util.Debugln(a.Description)
			desc := a.Description
			for _, reg := range coll_regs {
				re, err := regexp.Compile(reg)
				if err != nil {
					return err
				}

				if re.MatchString(desc) {
					if enaFl {
						contains := BuildFilterRegs(CFG.Env.RssFilter.Contain)
						exclusions := BuildFilterRegs(CFG.Env.RssFilter.Exclusion)
						if !FilterWithRegs(desc, contains, exclusions) {
							log.Printf("golbal filter: %s filtered", desc)
							continue
						}
					}
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
		err = rss.SetAutoDLRule(subj.ResourceUrl, subj.QbtCateg(), subj.Path, subj.RssPath(),
			enaFl, BuildFilterPerlReg(CFG.Env.RssFilter.Contain), BuildFilterPerlReg(CFG.Env.RssFilter.Exclusion))
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
			if ext.NoArgs() {
				if CFG.Env.EnabledFilter() {
					// use global filter
					r.UseRegex = true
					r.MustContain, r.MustNotContain = BuildFilterPerlReg(CFG.Env.RssFilter.Contain), BuildFilterPerlReg(CFG.Env.RssFilter.Exclusion)
				}
			} else {
				r.MustContain = ext.RssOption.MustContain
				r.MustNotContain = ext.RssOption.MustNotContain
			}
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
