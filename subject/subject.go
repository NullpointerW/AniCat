package subject

import (
	"errors"
	"fmt"
	"regexp"
	"strconv"
	"strings"
	"time"

	CFG "github.com/NullpointerW/anicat/conf"
	CC "github.com/NullpointerW/anicat/crawl/cover"
	IC "github.com/NullpointerW/anicat/crawl/information"
	RC "github.com/NullpointerW/anicat/crawl/resource"
	DL "github.com/NullpointerW/anicat/downloader"
	"github.com/NullpointerW/anicat/downloader/builtin"
	"github.com/NullpointerW/anicat/downloader/rss"
	"github.com/NullpointerW/anicat/downloader/torrent"
	"github.com/NullpointerW/anicat/errs"
	"github.com/NullpointerW/anicat/log"
	util "github.com/NullpointerW/anicat/utils"
	qbt "github.com/NullpointerW/go-qbittorrent-apiv2"
)

// Subject is the persisted record for a single bangumi subscription.
// Only JSON-tagged fields are written to disk; runtime state lives in the
// embedded RuntimeState which is initialised by runtimeInit and never serialised.
type Subject struct {
	// --- identity & metadata ---
	SubjId     int     `json:"subjId"`
	Name       string  `json:"name"`
	OriginName string  `json:"originName"`
	Alias      string  `json:"alias"`
	Typ        BgmiTyp `json:"typ"`

	// --- folder & timing ---
	FolderName string `json:"folderName"` // from TMDB
	FolderTime string `json:"folderTime"` // from TMDB
	Path       string `json:"path"`
	Season     string `json:"season"`
	Part       string `json:"part"` // e.g. pt1, pt2
	StartTime  string `json:"startTime"`
	EndTime    string `json:"endTime"`

	// --- state flags ---
	Finished  bool `json:"finished"`
	Terminate bool `json:"terminate"`
	Episode   int  `json:"episode"`

	// --- download config ---
	ResourceTyp     ResourceTyp `json:"resourceTyp"`
	ResourceUrl     string      `json:"resourceUrl"`
	BuiltinDownload bool        `json:"builtinDownload"`

	// --- qBittorrent downloader state ---
	TorrentHash string              `json:"torrentHash"`
	Pushed      map[string]string   `json:"pushed"`
	RssTorrents map[string]struct{} `json:"rssTorrents"`

	// --- builtin downloader state ---
	RssTorrentsName     map[string]struct{}         `json:"rssTorrentsName"`
	RssGuids            map[string]struct{}         `json:"rssGuids"`
	Filter              *FilterVerb                 `json:"filter,omitempty"`
	TorrentUrls         map[string]RssFileOptStrage `json:"torrentUrls"`
	TorrentFinishedUrls map[string]struct{}         `json:"torrentFinishedUrls"`

	// RuntimeState holds all goroutine, channel, and in-memory-only fields.
	// It is never serialised; runtimeInit populates it on startup.
	RuntimeState `json:"-"`
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
	if s.FinishedTorrentNameList == nil && s.ResourceTyp == RSS {
		f := make([]builtin.TorrentProgress, 0, len(s.TorrentFinishedUrls))
		for u := range s.TorrentFinishedUrls {
			f = append(f, builtin.TorrentProgress{
				Percentage: 100,
				Name:       s.TorrentUrls[u].Renamed,
			})
		}
		s.FinishedTorrentNameList = util.NewListView(f)
	} else if s.FinishedTorrentNameList == nil {
		f := make([]builtin.TorrentProgress, 0, len(s.TorrentFinishedUrls))
		for u := range s.TorrentFinishedUrls {
			f = append(f, builtin.TorrentProgress{
				Percentage: 100,
				Name:       u,
			})
		}
		s.FinishedTorrentNameList = util.NewListView(f)
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

// finalizeSubject handles the common steps after basic info is fetched:
// loading fields, folder init, cover scrape, download setup, and runtime start.
func finalizeSubject(subject *Subject, ext *Extra) (int, error) {
	GetSeason(subject)
	subject.GetPart()
	subject.trimName()
	if err := initFolder(subject); err != nil {
		return 0, err
	}
	lastS, err := FindLastSeason(subject.Path)
	if err != nil {
		return 0, err
	}
	if err = subject.scrapeCover(lastS); err != nil {
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
	if err = subject.writeJson(); err != nil {
		return 0, err
	}
	subject.runtimeInit(false)
	log.Info(log.Struct{"sid", subject.SubjId}, "create subject succeeded")
	return subject.SubjId, nil
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
	if err = subject.Loadfields(tips); err != nil {
		return 0, err
	}
	return finalizeSubject(subject, ext)
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
		var bgmurl, title string
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
	if err = subject.Loadfields(tips); err != nil {
		return 0, err
	}
	return finalizeSubject(subject, ext)
}

// folderSearch tries TMDB folder lookup in order: Name → OriginName → each Alias → strip season suffix.
func (s *Subject) folderSearch(tmdbTyp string) error {
	var err error
	s.FolderName, s.FolderTime, err = IC.FloderSearch(tmdbTyp, s.Name)
	if err == nil {
		return nil
	}
	if !errors.Is(err, errs.ErrCrawlNotFound) {
		return err
	}
	s.FolderName, s.FolderTime, err = IC.FloderSearch(tmdbTyp, s.OriginName)
	if err == nil {
		return nil
	}
	if !errors.Is(err, errs.ErrCrawlNotFound) {
		return err
	}
	for _, alias := range strings.Split(s.Alias, "|") {
		s.FolderName, s.FolderTime, err = IC.FloderSearch(tmdbTyp, alias)
		if err == nil {
			return nil
		}
		if !errors.Is(err, errs.ErrCrawlNotFound) {
			return err
		}
	}
	re := regexp.MustCompile(`第(.)季`)
	if match := re.FindStringSubmatch(s.Name); len(match) > 1 {
		trimmed := strings.TrimRight(strings.ReplaceAll(s.Name, fmt.Sprintf("第%s季", match[1]), ""), " ")
		s.FolderName, s.FolderTime, err = IC.FloderSearch(tmdbTyp, trimmed)
	}
	return err
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
	} else {
		s.StartTime = tips[IC.SubjMoveStartTime]
		s.Finished = true
	}
	s.Alias = tips[IC.Alias]

	tmdbTyp := IC.TMDB_TYP_TV
	if s.Typ == MOVIE {
		tmdbTyp = IC.TMDB_TYP_MOVIE
	}
	return s.folderSearch(tmdbTyp)
}

func (s *Subject) FetchInfo() error {
	tips, err := IC.BgmTVInfoScrape(s.SubjId)
	if err != nil {
		log.Error(log.Struct{"sid", s.SubjId, "err", err}, "FetchInfo: scrape failed, skipping update")
		return nil // non-fatal: keep running with existing info
	}
	wrap := errs.ErrWrapper{}
	wrap.Handle(func() error { return s.Loadfields(tips) })
	wrap.Handle(func() error { Mgr.Sync(); return nil })
	wrap.Handle(func() error { return s.writeJson() })
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

// applyFilter checks whether desc passes the filter defined by ext and global config.
// Returns true if the item should be downloaded, false if it should be skipped.
func applyFilter(sid int, desc string, enaFl bool, ext *Extra) bool {
	if enaFl {
		if !FilterWithRegs(desc, BuildFilterRegs(CFG.Env.RssFilter.Contain), BuildFilterRegs(CFG.Env.RssFilter.Exclusion)) {
			log.Info(log.Struct{"sid", sid, "filtered", desc}, "global filtered")
			return false
		}
		return true
	}
	if ext != nil && !ext.NoArgs() {
		if ext.RssOption.UseRegex {
			if !FilterWithCustomReg(desc, *ext) {
				log.Info(log.Struct{"sid", sid, "filtered", desc}, "custom filtered")
				return false
			}
		} else {
			if !FilterWithCustom(desc, *ext) {
				log.Info(log.Struct{"sid", sid, "filtered", desc}, "custom filtered")
				return false
			}
		}
	}
	return true
}

func download(subj *Subject, ext *Extra) error {
	if subj.ResourceTyp == Torrent {
		h, err := torrent.Add(subj.ResourceUrl, subj.Path, subj.QbtTag())
		subj.TorrentHash = h
		return err
	}

	if subj.Finished {
		it, err := rss.AddAndGetItems(subj.ResourceUrl, subj.RssPath())
		log.Debug(log.Struct{"sid", subj.SubjId, "rss path", subj.RssPath()}, "add RssResource")
		if err != nil {
			return err
		}
		enaFl := CFG.Env.EnabledFilter() && (ext == nil || ext.NoArgs())
		for _, a := range it.Articles {
			desc := a.Description
			log.Debug(log.Struct{"rssDesc", desc}, "traverse rssItems")
			isCollOrColl := subj.isCollection(desc)
			if !isCollOrColl {
				for _, reg := range coll_regs {
					re, err := regexp.Compile(reg)
					if err != nil {
						return err
					}
					if re.MatchString(desc) {
						isCollOrColl = true
						break
					}
				}
			}
			if !isCollOrColl {
				continue
			}
			if !applyFilter(subj.SubjId, desc, enaFl, ext) {
				continue
			}
			log.Info(log.Struct{"sid", subj.SubjId, "name", subj.Name, "matched", desc, "rss path", subj.RssPath()}, "matched collection")
			return subj.rssToTorr(a.TorrentURL)
		}

		log.Info(log.Struct{"sid", subj.SubjId, "name", subj.Name, "rss path", subj.RssPath()}, "not matched any collection")
		if err = torrent.AddCategroy(subj.QbtCateg()); err != nil {
			return err
		}
		if ext != nil && !ext.NoArgs() {
			return rss.SetAutoDLRule(subj.ResourceUrl, subj.QbtCateg(), subj.Path, subj.RssPath(),
				ext.RssOption.UseRegex, ext.RssOption.MustContain, ext.RssOption.MustNotContain)
		}
		return rss.SetAutoDLRule(subj.ResourceUrl, subj.QbtCateg(), subj.Path, subj.RssPath(),
			enaFl, BuildFilterPerlReg(CFG.Env.RssFilter.Contain), BuildFilterPerlReg(CFG.Env.RssFilter.Exclusion))
	}

	// ongoing series: set up auto-download rule
	if err := torrent.AddCategroy(subj.QbtCateg()); err != nil {
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
		if ext.NoArgs() && CFG.Env.EnabledFilter() {
			r.UseRegex = true
			r.MustContain = BuildFilterPerlReg(CFG.Env.RssFilter.Contain)
			r.MustNotContain = BuildFilterPerlReg(CFG.Env.RssFilter.Exclusion)
		} else {
			r.MustContain = ext.RssOption.MustContain
			r.MustNotContain = ext.RssOption.MustNotContain
		}
	}
	return rss.Download(r, subj.RssPath())
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
