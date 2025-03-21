package util

import (
	"strings"
)

type StringAppender struct {
	strings.Builder
}

func (sa *StringAppender) Append(ss ...string) {
	for _, s := range ss {
		sa.WriteString(s)
	}
}
