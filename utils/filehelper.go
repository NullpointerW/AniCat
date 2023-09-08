package util

import (
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

var (
	videoExt = []string{".mp4", ".rmvb", ".avi", ".flv", ".m2v", ".mkv", ".wmv", ".mp3", ".wav", ".mov"}

	subtitleExt = []string{".srt", ".ass", ".sub"}
)

func IsVideofile(fn string) bool {
	for _, ext := range videoExt {
		cmpFn := strings.ToLower(fn)
		if strings.HasSuffix(cmpFn, ext) {
			return true
		}
	}
	return false
}

func IsSubtitleFile(fn string) bool {
	for _, ext := range subtitleExt {
		cmpFn := strings.ToLower(fn)
		if strings.HasSuffix(cmpFn, ext) {
			return true
		}
	}
	return false
}

func IsJsonFile(fn string) bool {
	sep := strings.Split(fn, ".")
	if len(sep) < 2 {
		return false
	}
	return strings.ToLower(sep[1]) == "json"
}

func FileSeparatorConv(path string) string {
	return strings.ReplaceAll(path, "\\", "/")
}

// TrimExtensionAndGetEpi trim the rename file ext and name
// eg：
// example xxxS01E02 => S01E02
func TrimExtensionAndGetEpi(fn string) string {
	sep := "."
	sp := strings.Split(fn, sep)
	o := sep + sp[len(sp)-1]
	trimed := strings.ReplaceAll(fn, o, "")
	return trimed[len(trimed)-6:]
}

func IsRegexp(str string) bool {
	_, err := regexp.Compile(str)
	return err == nil
}

func GetExecutePath() (string, error) {
	executePath, err := os.Executable()
	if err != nil {
		return "", err
	}
	executePath, err = filepath.EvalSymlinks(filepath.Dir(executePath))
	if err != nil {
		return "", err
	}
	return executePath, nil
}
