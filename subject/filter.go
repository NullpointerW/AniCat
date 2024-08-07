package subject

import (
	"fmt"
	"github.com/NullpointerW/anicat/log"
	"regexp"
	"strings"
)

func BuildFilterPerlReg(vbs []string) string {
	var reg string
	const tmp = `(?=.*?%s)`
	if len(vbs) != 0 {
		reg += "(?i)"
		for _, ct := range vbs {
			vb := strings.ReplaceAll(ct, ",", "|")
			vb = "(" + vb + ")"
			reg += fmt.Sprintf(tmp, vb)
		}
		return reg
	} else {
		return ""
	}
}

func BuildFilterRegs(vbs []string) []string {
	if len(vbs) != 0 {
		regs := make([]string, 0, len(vbs))
		for _, ct := range vbs {
			vb := strings.ReplaceAll(ct, ",", "|")
			vb = "(?i)" + vb
			regs = append(regs, vb)
		}
		return regs
	} else {
		return nil
	}
}

func FilterWithRegs(s string, contains, exclusions []string) bool {
	var (
		containOk, exclusionOk bool
	)
	if len(contains) == 0 {
		containOk = true
	}
	if len(exclusions) == 0 {
		exclusionOk = true
	}
	if !containOk {
		//containOks := make([]bool, 0, len(contains))
		for _, reg := range contains {
			var ok bool
			csre, err := regexp.Compile(reg)
			if err != nil {
				log.Error(log.Struct{"err", err}, "globalFilter: contains regexp compile failed")
				ok = true
			} else {
				ok = csre.MatchString(s)
				log.Debug(log.Struct{"containRegexp", csre.String(), "matchingString", s, "matched", ok})
				//improve: break immediately
				if !ok {
					return false
				}
			}
			//containOks = append(containOks, ok)
		}
		containOk = true
		//for _, ok := range containOks {
		//	if !ok {
		//		containOk = false
		//		break
		//	}
		//}
	}

	if !exclusionOk {
		//exclusionOks := make([]bool, 0, len(exclusions))
		for _, reg := range exclusions {
			var ok bool
			clsre, err := regexp.Compile(reg)
			if err != nil {
				log.Error(log.Struct{"err", err}, "globalFilter: exclusions regexp compile failed")
				ok = true
			} else {
				ok = !clsre.MatchString(s)
				log.Debug(log.Struct{"exclusionRegexp", clsre.String(), "matchingString", s, "matched", ok})
				if !ok {
					return false
				}
			}
			//exclusionOks = append(exclusionOks, ok)
		}
		exclusionOk = true
		//for _, ok := range exclusionOks {
		//	if !ok {
		//		exclusionOk = false
		//		break
		//	}
		//}
	}
	return containOk && exclusionOk
}

func FilterWithCustomReg(s string, e Extra) bool {
	clsOk, exlOk := true, true
	if cls := e.RssOption.MustContain; cls != "" {
		clsReg, err := regexp.Compile(cls)
		if err != nil {
			log.Error(log.Struct{"err", err}, "customFilter: contains regexp compile failed")
		} else {
			clsOk = clsReg.MatchString(s)
		}
	}
	if exl := e.RssOption.MustNotContain; exl != "" {
		exlReg, err := regexp.Compile(exl)
		if err != nil {
			log.Error(log.Struct{"err", err}, "customFilter: contains regexp compile failed")
		} else {
			exlOk = !exlReg.MatchString(s)
		}
	}
	return clsOk && exlOk
}

func FilterWithCustom(s string, e Extra) bool {
	cls, ext := strings.Fields(e.RssOption.MustContain), strings.Fields(e.RssOption.MustNotContain)
	return FilterWithRegs(s, cls, ext)
}
