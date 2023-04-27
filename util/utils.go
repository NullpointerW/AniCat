package util

import (
	"time"
)

var Asia_Shanghai, _ = time.LoadLocation("Asia/Shanghai")

const YMDParseLayout = "2006年1月2日"

func ParseTime(strd string) (time.Time, error) {
	t, err := time.ParseInLocation(YMDParseLayout, strd, time.Local)
	if err != nil {
		return time.Time{}, err
	}
	return t, nil
}
