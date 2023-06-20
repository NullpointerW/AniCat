package util

import (
	"log"
	"regexp"
	"strconv"
	"strings"
	"time"
	"unicode"
)

var Asia_Shanghai, _ = time.LoadLocation("Asia/Shanghai")

const (
	YMDParseLayout  = "2006年1月2日"
	ShortDateLayout = "2006-01"
	Day             = 24 * time.Hour // 24h0m0s
	Week            = 7 * Day
)

var zh_cn_numb = map[rune]byte{
	'一': '1',
	'二': '2',
	'三': '3',
	'四': '4',
	'五': '5',
	'六': '6',
	'七': '7',
	'八': '8',
	'九': '9',
}

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

func FileSeparatorConv(path string) string {
	return strings.ReplaceAll(path, "\\", "/")
}

func IsRegexp(str string) bool {
	_, err := regexp.Compile(str)
	return err == nil
}

func IsNumber(s string) bool {
	for _, r := range s {
		if !unicode.IsNumber(r) {
			return false
		}
	}
	return true
}

func ConvertZhCnNumbToa(cnn string) string {
	runes := []rune(cnn)
	nl := len(runes)
	if nl == 1 {
		return string(zh_cn_numb[runes[0]])
	} else if nl == 2 {
		i, err := strconv.Atoi(string(zh_cn_numb[runes[1]]))
		if err != nil {
			log.Println(err)
			return "1"
		}
		return strconv.Itoa(10 + i)
	} else if nl == 3 {
		i, err := strconv.Atoi(string(zh_cn_numb[runes[0]]))
		if err != nil {
			log.Println(err)
			return "1"
		}
		i *= 10
		e, err := strconv.Atoi(string(zh_cn_numb[runes[2]]))
		if err != nil {
			log.Println(err)
			return "1"
		}
		i += e
		return strconv.Itoa(i)
	} else {
		log.Println("convert fail:cannot convert zh-cn numbers with more than 3 digits")
		return "1"
	}
}

func CheckZhCn(s string) bool {
	ok, _ := regexp.MatchString("[\u4e00-\u9fa5]", s)
	return ok
}
