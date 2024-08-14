package subject

import (
	"github.com/NullpointerW/anicat/downloader/builtin"
	"github.com/NullpointerW/anicat/rename"
	util "github.com/NullpointerW/anicat/utils"
	"github.com/anacrolix/torrent/metainfo"
	"github.com/anacrolix/torrent/storage"
	"path/filepath"
	"strings"
)

type FilePath struct {
	builtin.FileName
	DirPath string
}

func (d FilePath) Dir() storage.TorrentDirFilePathMaker {
	return func(baseDir string, info *metainfo.Info, infoHash metainfo.Hash) string {
		return d.DirPath
	}
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
