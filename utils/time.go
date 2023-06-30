package util

import (
	
	"time"
)

var Asia_Shanghai, _ = time.LoadLocation("Asia/Shanghai")

const (
	YMDParseLayout  = "2006年1月2日"
	ShortDateLayout = "200601"
	Day             = 24 * time.Hour // 24h0m0s
	Week            = 7 * Day
)

func ParseShortTime(strd string) (string, error) {
	t, err := ParseTime(strd)
	if err != nil {
		return "", err
	}
	return t.Format(ShortDateLayout), nil
}

func ParseTime(strd string) (time.Time, error) {
	t, err := time.ParseInLocation(YMDParseLayout, strd, time.Local)
	if err != nil {
		return time.Time{}, err
	}
	return t, nil
}

func ParseTimeStr(t time.Time) string {
	return t.Format(YMDParseLayout)
}





