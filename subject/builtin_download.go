package subject

import (
	"path/filepath"
	"strings"

	CFG "github.com/NullpointerW/anicat/conf"
	"github.com/NullpointerW/anicat/downloader/builtin"
	"github.com/NullpointerW/anicat/downloader/rss"
	"github.com/NullpointerW/anicat/rename"
	util "github.com/NullpointerW/anicat/utils"
	"github.com/anacrolix/torrent"
	"github.com/anacrolix/torrent/metainfo"
	"github.com/anacrolix/torrent/storage"
)

type DownloadedInfo struct {
	Size int
	Name int
}
type FilePath struct {
	builtin.FileName
	DirPath string
}

func (d FilePath) Dir() storage.TorrentDirFilePathMaker {
	return func(baseDir string, info *metainfo.Info, infoHash metainfo.Hash) string {
		return d.DirPath
	}
}

func (s *Subject) builtinDownload(mt builtin.MonitoredTorrent) {
	s.DetctchanBuiltin <- mt
}

func BuildFilter(s *Subject, ex *Extra) {
	if ex == nil || ex.NoArgs() {
		if CFG.Env.EnabledFilter() {
			c := BuildFilterRegs(CFG.Env.RssFilter.Contain)
			e := BuildFilterRegs(CFG.Env.RssFilter.Exclusion)
			s.Filter = &FilterVerb{
				false,
				c,
				e,
			}
		}
	} else {
		if ex.RssOption.UseRegex {
			s.Filter = &FilterVerb{
				true,
				ex.RssOption.MustContain,
				ex.RssOption.MustNotContain,
			}
			return
		}
		s.Filter = &FilterVerb{
			false,
			strings.Fields(ex.RssOption.MustContain),
			strings.Fields(ex.RssOption.MustNotContain),
		}
	}
}
func RssReader(s *Subject) error {
	if s.ResourceTyp == Torrent {
		return nil
	}
	var ff rss.FilterFunc
	if s.Filter != nil {
		ff = s.Filter.Filter()
	}
	r := rss.NewReader(s.ResourceUrl, s.RssGuids, ff)
	if s.Typ == TV {
		if s.Finished {
			its, ok, err := r.Seek()
			if err != nil {
				return err
			}
			if ok {
				for _, it := range its {
					if s.isCollection(it.Desc) {
						s.ResourceTyp = Torrent
						s.ResourceUrl = it.TorrUrl
						return nil
					}
				}
			}
		}
	} else {
		it, ok, err := r.ReadOne()
		if err != nil {
			return err
		}
		if ok {
			s.ResourceTyp = Torrent
			s.ResourceUrl = it.TorrUrl
			return nil
		}
	}
	s.RssReader = r
	return nil
}

type RssFileOpt struct {
	Renamed string
}

//	return func(opts storage.FilePathMakerOpts) string {
//		var parts []string
//		if opts.Info.Name != metainfo.NoName {
//			parts = append(parts, opts.Info.Name)
//		}
//		return filepath.Join(append(parts, opts.File.Path...)...)
//	}
func (r *RssFileOpt) Name() storage.FilePathMaker {
	return func(opts storage.FilePathMakerOpts) string {
		if len(opts.File.Path) != 0 {
			p := opts.File.Path[len(opts.File.Path)]
			if util.IsSubtitleFile(p) {
				ss := new(util.StringAppender)
				ss.Append(r.Renamed, " ", rename.SubtitleFileLang(p), filepath.Ext(p))
				p = ss.String()
			} else if util.IsVideofile(p) {
				p = r.Renamed + filepath.Ext(p)
			}
			return p
		}
		return r.Renamed + filepath.Ext(opts.Info.Name)
	}
}

type TorrFileOpt struct {
	subj *Subject
}

func (t *TorrFileOpt) Name() storage.FilePathMaker {
	return func(opts storage.FilePathMakerOpts) string {
		if len(opts.File.Path) != 0 {
			p := opts.File.Path[len(opts.File.Path)]
			if util.IsSubtitleFile(p) {
				r, err := renameTV(t.subj, p)
				if err != nil {
					return p
				}
				ext := filepath.Ext(r)
				r = strings.TrimSuffix(r, ext)
				ss := new(util.StringAppender)
				ss.Append(r, " ", rename.SubtitleFileLang(p), ext)
				return ss.String()
			} else if util.IsVideofile(p) {
				r, err := renameTV(t.subj, p)
				if err != nil {
					return p
				}
				return r
			}
		}
		tv, err := renameTV(t.subj, opts.Info.Name)
		if err != nil {
			return opts.Info.Name
		}
		return tv
	}
}

type MovieFileOpt struct{}

func (m *MovieFileOpt) Name() storage.FilePathMaker {
	return func(opts storage.FilePathMakerOpts) string {
		if len(opts.File.Path) != 0 {
			p := opts.File.Path[len(opts.File.Path)]
			return p
		}
		return opts.Info.Name
	}
}

type MagnetUrlSeeker struct {
}

func (_ *MagnetUrlSeeker) Seek(n string) (*torrent.TorrentSpec, error) {
	return torrent.TorrentSpecFromMagnetUri(n)
}

type RssFileOptStrage struct {
	Renamed string `json:"renamed"`
}

func (s *Subject) resumeRssDownload() error {
	if s.ResourceTyp != RSS {
		return nil
	}
	sr := util.SetSubtract(s.TorrentUrls, s.TorrentFinishedUrls)
	for u, v := range sr {
		ropt := RssFileOpt{
			v.Renamed,
		}
		fop := FilePath{FileName: &ropt, DirPath: s.Path}
		t, err := builtin.DefaultDownLoader.Download(u, fop, nil)
		if err != nil {
			return err
		}
		s.builtinDownload(builtin.MonitoredTorrent{Url: u,Rename: v.Renamed,Torrent: t})
	}
	return nil
}
