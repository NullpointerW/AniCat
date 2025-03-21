package util

import (
	"fmt"
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

func ConvertZhCnNumbToa(cnn string) (string,error) {
	runes := []rune(cnn)
	nl := len(runes)
	if nl == 1 {
		if cnn == "十" {
			return "10",nil
		}
		return string(zh_cn_numb[runes[0]]),nil
	} else if nl == 2 {
		i, err := strconv.Atoi(string(zh_cn_numb[runes[1]]))
		if err != nil {
			return "1",fmt.Errorf("convert cn_number to int failed: %w",err)
		}
		return strconv.Itoa(10 + i),nil
	} else if nl == 3 {
		i, err := strconv.Atoi(string(zh_cn_numb[runes[0]]))
		if err != nil {
			return "1",fmt.Errorf("convert cn_number to int failed: %w",err)
		}
		i *= 10
		e, err := strconv.Atoi(string(zh_cn_numb[runes[2]]))
		if err != nil {
			return "1",fmt.Errorf("convert cn_number to int failed: %w",err)
		}
		i += e
		return strconv.Itoa(i),nil
	} else {
		return "1",fmt.Errorf("convert cn_number to int failed: cannot convert zh-cn numbers with more than 3 digits")
	}
}

func CheckZhCn(s string) bool {
	ok, _ := regexp.MatchString("[\u4e00-\u9fa5]", s)
	return ok
}
