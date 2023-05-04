package util

import (
	"regexp"
	"strings"
	"time"
)

var Asia_Shanghai, _ = time.LoadLocation("Asia/Shanghai")

const (
	YMDParseLayout = "2006年1月2日"
	Day            = 24 * time.Hour // 24h0m0s
	Week           = 7 * Day
)

func ParseTime(strd string) (time.Time, error) {
	t, err := time.ParseInLocation(YMDParseLayout, strd, time.Local)
	if err != nil {
		return time.Time{}, err
	}
	return t, nil
}

func FileSeparatorConv(path string) string {
	return strings.ReplaceAll(path, "\\", "/")
}

func IsRegexp(str string) bool {
	_, err := regexp.Compile(str)
	return err == nil
}
