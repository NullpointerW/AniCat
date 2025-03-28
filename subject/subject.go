package subject

import (
	"context"
	"errors"
	"fmt"
	"regexp"
	"strconv"
	"strings"

	// "sync"
	"time"

	DL "github.com/NullpointerW/anicat/downloader"
	"github.com/NullpointerW/anicat/log"

	CC "github.com/NullpointerW/anicat/crawl/cover"
	IC "github.com/NullpointerW/anicat/crawl/information"
	RC "github.com/NullpointerW/anicat/crawl/resource"

	// DL "github.com/NullpointerW/anicat/downloader"
	CFG "github.com/NullpointerW/anicat/conf"

	"github.com/NullpointerW/anicat/downloader/builtin"
	"github.com/NullpointerW/anicat/downloader/rss"
	"github.com/NullpointerW/anicat/downloader/torrent"
	"github.com/NullpointerW/anicat/errs"
	util "github.com/NullpointerW/anicat/utils"
	qbt "github.com/NullpointerW/go-qbittorrent-apiv2"
)

// Subject as a basic object of each bangumi
// program will load it from OS file and manage them
// it will be refreshed to OS file after some of the fields have updating
type Subject struct {
	SubjId      int         `json:"subjId"`
	FolderName  string      `json:"folderName"` // source from tmdb
	Name        string      `json:"name"`
	OriginName  string      `json:"originName"`
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
	Part        string      `json:"part"` // eg: pt1、pt2
	// used while `ResourceTyp` is `Torrent`
	TorrentHash string `json:"torrentHash"`
	// manager use ctx cancel func to Exit goroutine running the current subject.
	// when delete a subject manager,should run the cancel func and if a goroutine is running
	// for this subject it will Exit.
	// Context is hold by subject-running goroutine
	// while subject-running goroutine Exit actively func should be called
	Exit context.CancelFunc `json:"-"`
	// before detector-goroutine push to subject,Check if this channel is closed.
	// before exit Exited channel should be closed
	Exited chan struct{} `json:"-"`
	// While detector-goroutine detected that the resource downloader of the subject is completed
	// it will send downLoad message to subject-running goroutine
	// received and push to terminal
	PushChan chan qbt.Torrent `json:"-"`
	// The anime series of this project has already ended and all episodes have been downloaded.
	// while init,if this flag is true then there is no need to start a goroutine to run it
	// exit actively flag should be set to true
	Terminate bool `json:"terminate"`
	// a Set store all pushed renamed episodes,avoid duplicate push.
	// content will like be `xxx S01E01,xxx S01E05...`
	Pushed        map[string]string   `json:"pushed"`
	RssTorrents   map[string]struct{} `json:"rssTorrents"`
	OperationChan chan Operate        `json:"-"`
	// builtin-downloader filed
	BuiltinDownload         bool                                    `json:"builtinDownload"`
	RssTorrentsName         map[string]struct{}                     `json:"rssTorrentsName"`
	RssReader               *rss.Reader                             `json:"-"`
	RssGuids                map[string]struct{}                     `json:"rssGuids"`
	Filter                  *FilterVerb                             `json:"filter"`
	TorrentUrls             map[string]RssFileOptStrage             `json:"torrentUrls"`
	TorrentFinishedUrls     map[string]struct{}                     `json:"torrentFinishedUrls"`
	MonitorchanBuiltin      chan builtin.MonitoredTorrent           `json:"-"`
	PushChanBuiltin         chan builtin.MonitoredTorrent           `json:"-"`
	FinihsedTorrentNameList *util.ListView[builtin.TorrentProgress] `json:"-"`
	TorrentMonitor          *builtin.TorrentProgressMonitor         `json:"-"`
}
type subjOp int

const (
	Rename subjOp = iota
)

type Operate struct {
	op  subjOp
	arg any
}

func NewOperate(op subjOp, arg any) Operate {
	return Operate{op, arg}
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
		Name           string
	}
}

func (ex *Extra) NoArgs() bool {
	opt := ex.RssOption
	return opt.MustContain == "" && opt.MustNotContain == ""
}

func (s *Subject) initializeFinishedTorrentNameList() {
	if s.FinihsedTorrentNameList == nil && s.ResourceTyp == RSS {
		f := make([]builtin.TorrentProgress, 0, len(s.TorrentFinishedUrls))
		for u := range s.TorrentFinishedUrls {
			f = append(f, builtin.TorrentProgress{
				Percentage: 100,
				Name:       s.TorrentUrls[u].Renamed,
			})
		}
		s.FinihsedTorrentNameList = util.NewListView(f)
	} else if s.FinihsedTorrentNameList == nil {
		f := make([]builtin.TorrentProgress, 0, len(s.TorrentFinishedUrls))
		for u := range s.TorrentFinishedUrls {
			// torrent type will store name to finished list
			f = append(f, builtin.TorrentProgress{
				Percentage: 100,
				Name:       u})
		}
		s.FinihsedTorrentNameList = util.NewListView(f)
	}
}

// QbtTag The tag used when adding a torrent with qbt
// can be used to monitor the downloader status of resources
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
	if Mgr.Get(sid) != nil {
		return 0, fmt.Errorf("%w:sid:%d", errs.ErrSubjectAlreadyExisted, sid)
	}
	subject.SubjId = sid
	err = subject.Loadfields(tips)
	if err != nil {
		return 0, err
	}
	GetSeason(subject)
	subject.GetPart()
	subject.trimName()
	err = initFolder(subject)
	if err != nil {
		return 0, err
	}
	lastS, err := FindLastSeason(subject.Path)
	if err != nil {
		return 0, err
	}

	err = subject.scrapeCover(lastS)
	if err != nil {
		return 0, err
	}

	subject.BuiltinDownload = CFG.Env.BuiltinDownloader
	if !subject.BuiltinDownload {
		err = download(subject, ext)
	} else {
		err = BuiltinDownloadPrepare(subject, ext)
	}
	if err != nil {
		return 0, err
	}

	// create Info-Json after init completed
	err = subject.writeJson()
	if err != nil {
		return 0, err
	}

	subject.runtimeInit(false)

	log.Info(log.Struct{"sid", subject.SubjId}, "create subject succeeded")
	return sid, nil
}

// CreateSubjectViaFeed use a specified rss-feed url as the resource to create a subject,
// if arg `name` is not empty,then will use specified name to fetch info,
// otherwise parse the feed for link or title to fetch it.
// eg:
//
//		`add --feed <url>`
//	 we fetch bgmTV link first,if it doesn't exist,then get title
//
//		`add --feed <url> --name <specified-name>`
//	 use specified name only
func CreateSubjectViaFeed(feed, name string, ext *Extra) (int, error) {
	subject := &Subject{ResourceTyp: RSS, ResourceUrl: feed}
	fp := rss.Parser{Feed: feed}
	var (
		err  error
		tips map[string]string
	)
	if name != "" {
		tips, err = IC.Scrape(name)
	} else {
		var (
			bgmurl string
			title  string
		)
		if title, bgmurl, err = fp.GetTitleAndLink(); err != nil {
			return 0, err
		} else if bgmurl == "" {
			log.Warn(log.Struct{"feed", feed, "err", err}, errs.ErrNoLinkFoundOnRssFeed)
			tips, err = IC.Scrape(title)
		} else {
			tips, err = IC.DoScrape(bgmurl)
		}
	}
	if err != nil {
		return 0, err
	}
	sid, err := strconv.Atoi(tips[IC.SubjId])
	if err != nil {
		return 0, err
	}
	if Mgr.Get(sid) != nil {
		return 0, fmt.Errorf("%w: sid=%d", errs.ErrSubjectAlreadyExisted, sid)
	}
	subject.SubjId = sid
	err = subject.Loadfields(tips)
	if err != nil {
		return 0, err
	}
	GetSeason(subject)
	subject.GetPart()
	subject.trimName()
	err = initFolder(subject)
	if err != nil {
		return 0, err
	}
	lastS, err := FindLastSeason(subject.Path)
	if err != nil {
		return 0, fmt.Errorf("getLastSeason failed: %w", err)
	}
	err = subject.scrapeCover(lastS)
	if err != nil {
		return 0, err
	}

	if !subject.BuiltinDownload {
		err = download(subject, ext)
	} else {
		err = BuiltinDownloadPrepare(subject, ext)
	}
	if err != nil {
		return 0, err
	}

	err = subject.writeJson()
	if err != nil {
		return 0, err
	}
	subject.runtimeInit(false)
	log.Info(log.Struct{"sid", subject.SubjId}, "create subject succeeded")
	return sid, nil
}

func (s *Subject) Loadfields(tips map[string]string) error {
	defer func() {
		if s.FolderName != "" && strings.ContainsRune(s.FolderName, '?') {
			s.FolderName = strings.ReplaceAll(s.FolderName, "?", "？")
		}
	}()
	s.Name = tips[IC.SubjName]
	s.OriginName = tips[IC.SubjOriginName]
	if s.Name == "" {
		s.Name = s.OriginName
	}
	if _, e := tips[IC.SubjStartTime]; e {
		s.Typ = TV
	} else {
		s.Typ = MOVIE
	}
	s.Episode, _ = strconv.Atoi(tips[IC.SubjEpisode])
	if s.Typ == TV {
		s.StartTime = tips[IC.SubjStartTime]
		if et, e := tips[IC.SubjectEndTime]; e {
			n := time.Now()
			eti, err := util.ParseTime(et, util.YMDParseLayout)
			if err != nil {
				reg := regexp.MustCompile(reg0_bgmTvTime)
				submatch := reg.FindStringSubmatch(et)
				if submatch != nil {
					et = submatch[0]
					eti, err = util.ParseTime(et, util.YMDParseLayout)
					if err != nil {
						return err
					}
				} else {
					return err
				}
			}
			s.EndTime = et
			s.Finished = n.After(eti) || n.Equal(eti)
		}
	} else { // if movie finished
		s.StartTime = tips[IC.SubjMoveStartTime]
		s.Finished = true
	}
	s.Alias = tips[IC.Alias]

	// fetch folder info,source from tmdb
	var tmdbTyp = IC.TMDB_TYP_TV
	if s.Typ == MOVIE {
		tmdbTyp = IC.TMDB_TYP_MOVIE
	}
	var err error
	// First, attempt to search for the folder using the subject's Name
	s.FolderName, s.FolderTime, err = IC.FloderSearch(tmdbTyp, s.Name)
	if err != nil {
		if errors.Is(err, errs.ErrCrawlNotFound) {
			// If the search failed because the folder was not found,
			// try again using the subject's OriginName
			s.FolderName, s.FolderTime, err = IC.FloderSearch(tmdbTyp, s.OriginName)
		}
		if errors.Is(err, errs.ErrCrawlNotFound) {
			// If the search still cannot found, try using each alias from the subject's Alias field
			for _, n := range strings.Split(s.Alias, "|") {
				s.FolderName, s.FolderTime, err = IC.FloderSearch(tmdbTyp, n)
				if err == nil {
					return nil
				} else if !errors.Is(err, errs.ErrCrawlNotFound) {
					// If the search failed with an error other than ErrCrawlNotFound, return the error
					return err
				}
			}
			// If all aliases failed, try removing the season number from the subject's Name and search again
			re := regexp.MustCompile(`第(.)季`)
			match := re.FindStringSubmatch(s.Name)
			if len(match) > 1 {
				season := match[1]
				n := strings.ReplaceAll(s.Name, fmt.Sprintf("第%s季", season), "")
				n = strings.TrimRight(n, " ")
				s.FolderName, s.FolderTime, err = IC.FloderSearch(tmdbTyp, n)
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
		return s.Loadfields(tips)
	})
	wrap.Handle(func() error {
		Mgr.Sync()
		return nil
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
		log.Debug(log.Struct{"subtitleGroup", opt.Group}, "specify subtitleGroup")
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
	log.Debug(log.Struct{"resource", u, "bgmtvUrl", bgm, "isRss", isrss}, "solvedResource")
	return bgm, nil
}

func download(subj *Subject, ext *Extra) error {
	if subj.ResourceTyp == Torrent {
		h, err := torrent.Add(subj.ResourceUrl, subj.Path, subj.QbtTag())
		subj.TorrentHash = h
		return err
	} else if subj.Finished {
		it, err := rss.AddAndGetItems(subj.ResourceUrl, subj.RssPath())
		log.Debug(log.Struct{"sid", subj.SubjId, "rss path ", subj.RssPath()}, "add RssResource")
		if err != nil {
			return err
		}
		enaFl := CFG.Env.EnabledFilter() && (ext == nil || ext.NoArgs())
	rssTraverse:
		for _, a := range it.Articles {
			log.Debug(log.Struct{"rssDesc", a.Description}, "traverse rssItems")
			desc := a.Description
			if subj.isCollection(desc) {
				if enaFl {
					contains := BuildFilterRegs(CFG.Env.RssFilter.Contain)
					exclusions := BuildFilterRegs(CFG.Env.RssFilter.Exclusion)
					if !FilterWithRegs(desc, contains, exclusions) {
						log.Info(log.Struct{"sid", subj.SubjId, "filtered", desc}, "global filtered")
						continue
					}
				} else if !(ext == nil || ext.NoArgs()) {
					if ext.RssOption.UseRegex {
						if !FilterWithCustomReg(desc, *ext) {
							log.Info(log.Struct{"sid", subj.SubjId, "filtered", desc}, "custom filtered")
							continue
						}
					} else {
						if !FilterWithCustom(desc, *ext) {
							log.Info(log.Struct{"sid", subj.SubjId, "filtered", desc}, "custom filtered")
							continue
						}
					}
				}
				log.Info(log.Struct{"sid", subj.SubjId, "name", subj.Name, "matched", desc, "rss path", subj.RssPath()}, "matched collection")
				return subj.rssToTorr(a.TorrentURL)
			}
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
							log.Info(log.Struct{"sid", subj.SubjId, "filtered", desc}, "global filter")
							continue rssTraverse
						}
					} else if !(ext == nil || ext.NoArgs()) {
						if ext.RssOption.UseRegex {
							if !FilterWithCustomReg(desc, *ext) {
								log.Info(log.Struct{"sid", subj.SubjId, "filtered", desc}, "custom filtered")
								continue
							}
						} else {
							if !FilterWithCustom(desc, *ext) {
								log.Info(log.Struct{"sid", subj.SubjId, "filtered", desc}, "custom filtered")
								continue
							}
						}
					}
					log.Info(log.Struct{"sid", subj.SubjId, "name", subj.Name, "matched", desc, "rss path", subj.RssPath()}, "matched collection")
					return subj.rssToTorr(a.TorrentURL)
				}
			}
		}
		log.Info(log.Struct{"sid", subj.SubjId, "name", subj.Name, "rss path", subj.RssPath()}, "not matched any collection")
		err = torrent.AddCategroy(subj.QbtCateg())
		if err != nil {
			return err
		}
		if !(ext == nil || ext.NoArgs()) {
			err = rss.SetAutoDLRule(subj.ResourceUrl, subj.QbtCateg(), subj.Path, subj.RssPath(),
				ext.RssOption.UseRegex, ext.RssOption.MustContain, ext.RssOption.MustNotContain)
		} else {
			err = rss.SetAutoDLRule(subj.ResourceUrl, subj.QbtCateg(), subj.Path, subj.RssPath(),
				enaFl, BuildFilterPerlReg(CFG.Env.RssFilter.Contain), BuildFilterPerlReg(CFG.Env.RssFilter.Exclusion))
		}
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
	ns = append(ns, s.Name, s.OriginName)
	as := strings.Split(s.Alias, "|")
	ns = append(ns, as...)
	for _, n := range ns {
		re := regexp.MustCompile(zhreg)
		match := re.FindStringSubmatch(n)
		if len(match) > 1 {
			m := match[1]
			if iszh := util.CheckZhCn(m); iszh {
				m, _ = util.ConvertZhCnNumbToa(m)
			}
			s.Season = fmt.Sprintf("%02s", m)
			return
		}
		for _, rg := range sregs {
			re := regexp.MustCompile(rg)
			match := re.FindStringSubmatch(n)
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

func (s *Subject) GetPart() {
	var ns []string
	ns = append(ns, s.Name, s.OriginName)
	as := strings.Split(s.Alias, "|")
	ns = append(ns, as...)
	for _, n := range ns {
		for _, reg := range part_regs {
			re := regexp.MustCompile(reg)
			match := re.FindStringSubmatch(n)
			if len(match) > 1 {
				m := match[1]
				s.Part = fmt.Sprintf("pt%s", m)
				return
			}
		}
		re, _ := regexp.Compile(reg_part2)
		matched := re.MatchString(n)
		if matched {
			s.Part = "pt2"
			return
		}
	}
}

func (s *Subject) GetSeasonAndPart() string {
	return s.Season + s.Part
}

func (s *Subject) scrapeCover(lastS int) error {
	curr, err := strconv.Atoi(s.Season)
	if err != nil {
		return fmt.Errorf("getCurrSeason failed: %w", err)
	}

	if curr > lastS {
		cp := s.Path + "/" + CoverFN
		err = CC.TouchbgmCoverImg(s.SubjId, cp)
		if err != nil {
			log.Error(log.Struct{"err", err}, "scrapeCover from bgmTV failed")
			err = CC.DOUBANCoverScraper(cp, s.Name)
			if err != nil {
				retry := 0
				for err == errs.ErrCoverDownLoadZeroSize {
					retry++
					if retry >= 3 {
						return err
					}
					time.Sleep(500 * time.Millisecond)
					err = CC.DOUBANCoverScraper(cp, s.Name)
				}
				if err != nil {
					return err
				}
			}
		}
	}
	return nil
}

func (s *Subject) isCollection(desc string) bool {
	re := regexp.MustCompile(reg2_coll)
	m := re.FindStringSubmatch(desc)
	if len(m) > 1 {
		m := m[1]
		i, _ := strconv.Atoi(m)
		return i >= s.Episode
	}
	return false
}

func (s *Subject) rssToTorr(torrUrl string) (err error) {
	err = rss.RmRss(s.RssPath())
	if err != nil {
		return err
	}
	s.ResourceTyp = Torrent
	h, err := torrent.Add(torrUrl, s.Path, s.QbtTag())
	s.TorrentHash = h
	return err
}

func (s *Subject) Rename(new string) error {
	if s.Terminate && s.ResourceTyp == Torrent {
		fs, err := DL.Qbt.Files(s.TorrentHash)
		if err != nil {
			return err
		}
		for _, f := range fs {
			old := f.Name
			newFullName := strings.ReplaceAll(old, s.FolderName, new)
			err = DL.Qbt.RenameFile(s.TorrentHash, old, newFullName)
			if err != nil {
				return err
			}
		}
	}
	if s.ResourceTyp == RSS {
		for th := range s.RssTorrents {
			fs, err := DL.Qbt.Files(th)
			if err != nil {
				return err
			}
			for _, f := range fs {
				old := f.Name
				newFullName := strings.ReplaceAll(old, s.FolderName, new)
				err = DL.Qbt.RenameFile(th, old, newFullName)
				if err != nil {
					return err
				}
			}
		}
	}
	s.FolderName = new
	err := s.writeJson()
	return err
}
func BuiltinDownloadPrepare(s *Subject, ex *Extra) error {
	if s.ResourceTyp != Torrent {
		BuildFilter(s, ex)
	}
	return RssReader(s)
}

func (s *Subject) ElapsedfromFinishedTime(e time.Duration) (bool, error) {
	if !s.Finished || s.EndTime == "" {
		return false, nil
	}
	end, err := util.ParseTime(s.EndTime, util.YMDParseLayout)
	if err != nil {
		return false, err
	}
	return time.Since(end) >= e, nil
}
