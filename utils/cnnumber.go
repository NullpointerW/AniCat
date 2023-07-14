package util

import (
	"log"
	"regexp"
	"strconv"
	"unicode"
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
		if cnn == "十" {
			return "10"
		}
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
