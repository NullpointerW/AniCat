package util

import (
	"regexp"
	"strings"
)

var videoExt = []string{".mp4", ".rmvb", ".avi", ".flv", ".m2v", ".mkv", ".wmv", ".mp3", ".wav", ".mov"}

func IsVideofile(fn string) bool {
	for _, ext := range videoExt {
		if strings.HasSuffix(fn, ext) {
			return true
		}
	}
	return false
}

func FileSeparatorConv(path string) string {
	return strings.ReplaceAll(path, "\\", "/")
}

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
