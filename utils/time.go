package util

import (
	"time"
)

var Asia_Shanghai, _ = time.LoadLocation("Asia/Shanghai")

const (
	YMDParseLayout   = "2006年1月2日"
	YMD02ParseLayout = "2006年01月02日"
	ShortDateLayout  = "200601"
	Day              = 24 * time.Hour // 24h0m0s
	Week             = 7 * Day
)

func ParseShort02Time(strd string) (string, error) {
	t, err := ParseTime(strd, YMD02ParseLayout)
	if err != nil {
		return "", err
	}
	return t.Format(ShortDateLayout), nil
}

func ParseShortTime(strd string) (string, error) {
	t, err := ParseTime(strd, YMDParseLayout)
	if err != nil {
		return "", err
	}
	return t.Format(ShortDateLayout), nil
}

func ParseTime(strd, layout string) (time.Time, error) {
	t, err := time.ParseInLocation(layout, strd, time.Local)
	if err != nil {
		return time.Time{}, err
	}
	return t, nil
}

func ParseTimeStr(t time.Time) string {
	return t.Format(YMDParseLayout)
}
