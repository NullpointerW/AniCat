package util

import (
	"regexp"
	"strings"
	"time"
)

var Asia_Shanghai, _ = time.LoadLocation("Asia/Shanghai")

const (
	YMDParseLayout  = "2006年1月2日"
	ShortDateLayout = "2006-01"
	Day             = 24 * time.Hour // 24h0m0s
	Week            = 7 * Day
)

type Tuple[F, S any] struct {
	slot1 F
	slot2 S
}

func (t *Tuple[F, S]) Get0() F {
	return t.slot1
}

func (t *Tuple[F, S]) Get1() S {
	return t.slot2
}

func NewTuple[F, S any](a1 F, a2 S) Tuple[F, S] {
	return Tuple[F, S]{
		a1,
		a2,
	}
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

func FileSeparatorConv(path string) string {
	return strings.ReplaceAll(path, "\\", "/")
}

func IsRegexp(str string) bool {
	_, err := regexp.Compile(str)
	return err == nil
}
