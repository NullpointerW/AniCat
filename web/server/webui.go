package webserver

import (
	"encoding/json"
	"fmt"
	"io/fs"
	"net/http"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strconv"
	"strings"

	CFG "github.com/NullpointerW/anicat/conf"
	"github.com/NullpointerW/anicat/log"
	"github.com/NullpointerW/anicat/subject"
	util "github.com/NullpointerW/anicat/utils"
)

type Episode struct {
	Index    int    `json:"index"`
	Filename string `json:"filename"`
	SxxExx   string `json:"sxxexx"`
	Season   int    `json:"season"`
	Ep       int    `json:"episode"`
	Path     string `json:"path"`
}

type BangumiItem struct {
	SubjId    int    `json:"subjId"`
	Name      string `json:"name"`
	Season    string `json:"season"`
	Typ       string `json:"typ"`
	Terminate bool   `json:"terminate"`
	HasCover  bool   `json:"hasCover"`
}

type BangumiDetail struct {
	SubjId    int       `json:"subjId"`
	Name      string    `json:"name"`
	Season    string    `json:"season"`
	Typ       string    `json:"typ"`
	Terminate bool      `json:"terminate"`
	Episodes  []Episode `json:"episodes"`
	HasCover  bool      `json:"hasCover"`
}

var sxxexxRe = regexp.MustCompile(`(?i)[Ss](\d+)[Ee](\d+)`)

func scanEpisodes(s *subject.Subject) []Episode {
	var episodes []Episode
	if s.Path == "" {
		return episodes
	}
	_ = filepath.WalkDir(s.Path, func(path string, d os.DirEntry, err error) error {
		if err != nil || d.IsDir() {
			return nil
		}
		fn := d.Name()
		if !util.IsVideofile(fn) {
			return nil
		}
		rel, _ := filepath.Rel(s.Path, path)
		rel = strings.ReplaceAll(rel, "\\", "/")
		ep := Episode{Filename: fn, Path: rel}
		if m := sxxexxRe.FindStringSubmatch(fn); len(m) == 3 {
			ep.SxxExx = m[0]
			ep.Season, _ = strconv.Atoi(m[1])
			ep.Ep, _ = strconv.Atoi(m[2])
		}
		episodes = append(episodes, ep)
		return nil
	})
	sort.Slice(episodes, func(i, j int) bool {
		if episodes[i].Season != episodes[j].Season {
			return episodes[i].Season < episodes[j].Season
		}
		return episodes[i].Ep < episodes[j].Ep
	})
	for i := range episodes {
		episodes[i].Index = i + 1
	}
	return episodes
}

func hasCover(path string) bool {
	_, err := os.Stat(filepath.Join(path, subject.CoverFN))
	return err == nil
}

type statusRecorder struct {
	http.ResponseWriter
	status int
}

func (sr *statusRecorder) WriteHeader(code int) {
	sr.status = code
	sr.ResponseWriter.WriteHeader(code)
}

func logMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		rec := &statusRecorder{ResponseWriter: w, status: http.StatusOK}
		next.ServeHTTP(rec, r)
		log.Info(log.Struct{"method", r.Method, "path", r.URL.Path, "status", rec.status}, "webui: request")
	})
}

func cors(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusNoContent)
			return
		}
		next(w, r)
	}
}

func jsonResp(w http.ResponseWriter, v any) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	_ = json.NewEncoder(w).Encode(v)
}

func handleList(w http.ResponseWriter, _ *http.Request) {
	ls := subject.Mgr.List()
	log.Info(log.Struct{"count", len(ls)}, "webui: handleList queried")
	for _, s := range ls {
		log.Info(log.Struct{"sid", s.SubjId, "name", s.Name, "season", s.Season, "typ", s.Typ.String(), "terminate", s.Terminate}, "webui: handleList item")
	}
	out := make([]BangumiItem, 0, len(ls))
	for _, s := range ls {
		out = append(out, BangumiItem{
			SubjId:    s.SubjId,
			Name:      s.Name,
			Season:    s.Season,
			Typ:       s.Typ.String(),
			Terminate: s.Terminate,
			HasCover:  hasCover(s.Path),
		})
	}
	sort.Slice(out, func(i, j int) bool { return out[i].SubjId > out[j].SubjId })
	jsonResp(w, out)
}

func handleDetail(w http.ResponseWriter, r *http.Request) {
	sid, err := strconv.Atoi(strings.TrimPrefix(r.URL.Path, "/api/bangumi/"))
	if err != nil {
		http.Error(w, "invalid sid", http.StatusBadRequest)
		return
	}
	s := subject.Mgr.Get(sid)
	if s == nil {
		http.Error(w, "not found", http.StatusNotFound)
		return
	}
	jsonResp(w, BangumiDetail{
		SubjId:    s.SubjId,
		Name:      s.Name,
		Season:    s.Season,
		Typ:       s.Typ.String(),
		Terminate: s.Terminate,
		Episodes:  scanEpisodes(s),
		HasCover:  hasCover(s.Path),
	})
}

func handleCover(w http.ResponseWriter, r *http.Request) {
	sid, err := strconv.Atoi(strings.TrimPrefix(r.URL.Path, "/api/cover/"))
	if err != nil {
		http.Error(w, "invalid sid", http.StatusBadRequest)
		return
	}
	s := subject.Mgr.Get(sid)
	if s == nil {
		http.Error(w, "not found", http.StatusNotFound)
		return
	}
	http.ServeFile(w, r, filepath.Join(s.Path, subject.CoverFN))
}

func handleVideo(w http.ResponseWriter, r *http.Request) {
	rest := strings.TrimPrefix(r.URL.Path, "/api/video/")
	slash := strings.Index(rest, "/")
	if slash < 0 {
		http.Error(w, "invalid path", http.StatusBadRequest)
		return
	}
	sid, err := strconv.Atoi(rest[:slash])
	if err != nil {
		http.Error(w, "invalid sid", http.StatusBadRequest)
		return
	}
	s := subject.Mgr.Get(sid)
	if s == nil {
		http.Error(w, "not found", http.StatusNotFound)
		return
	}
	rel := filepath.Clean(rest[slash+1:])
	if strings.Contains(rel, "..") {
		http.Error(w, "forbidden", http.StatusForbidden)
		return
	}
	cleanBase := filepath.Clean(s.Path)
	full := filepath.Join(cleanBase, rel)
	if !strings.HasPrefix(full, cleanBase+string(os.PathSeparator)) {
		http.Error(w, "forbidden", http.StatusForbidden)
		return
	}
	http.ServeFile(w, r, full)
}

// Listen starts the Web UI HTTP server on webui-port (default 12315).
func Listen() {
	port := CFG.Env.WebUIPort
	if port == 0 {
		port = 12315
	}
	mux := http.NewServeMux()

	mux.HandleFunc("/api/bangumi/", cors(func(w http.ResponseWriter, r *http.Request) {
		p := r.URL.Path
		if p == "/api/bangumi/" || p == "/api/bangumi" {
			handleList(w, r)
		} else {
			handleDetail(w, r)
		}
	}))
	mux.HandleFunc("/api/cover/", cors(handleCover))
	mux.HandleFunc("/api/video/", cors(handleVideo))
	subFS, _ := fs.Sub(staticFiles, "static")
	mux.Handle("/", http.FileServer(http.FS(subFS)))

	addr := fmt.Sprintf(":%d", port)
	log.Info(log.Struct{"addr", addr}, "webui: server started")
	if err := http.ListenAndServe(addr, logMiddleware(mux)); err != nil {
		log.Error(log.Struct{"err", err}, "webui: server error")
	}
}

// scanEpisodesFromPath is a testable helper that scans a directory path directly.
func scanEpisodesFromPath(path string) []Episode {
	if path == "" {
		return nil
	}
	var episodes []Episode
	_ = filepath.WalkDir(path, func(fp string, d os.DirEntry, err error) error {
		if err != nil || d.IsDir() {
			return nil
		}
		fn := d.Name()
		if !util.IsVideofile(fn) {
			return nil
		}
		rel, _ := filepath.Rel(path, fp)
		rel = strings.ReplaceAll(rel, "\\", "/")
		ep := Episode{Filename: fn, Path: rel}
		if m := sxxexxRe.FindStringSubmatch(fn); len(m) == 3 {
			ep.SxxExx = m[0]
			ep.Season, _ = strconv.Atoi(m[1])
			ep.Ep, _ = strconv.Atoi(m[2])
		}
		episodes = append(episodes, ep)
		return nil
	})
	sort.Slice(episodes, func(i, j int) bool {
		if episodes[i].Season != episodes[j].Season {
			return episodes[i].Season < episodes[j].Season
		}
		return episodes[i].Ep < episodes[j].Ep
	})
	for i := range episodes {
		episodes[i].Index = i + 1
	}
	return episodes
}
