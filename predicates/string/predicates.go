package string

import (
	"regexp"
	"strings"
	"unicode/utf8"
)

func NotEmpty(s string) bool {
	return strings.TrimSpace(s) != ""
}

// Lower :: int -> Upper :: int -> string -> bool
func InBetween(lower int) func(int) func(string) bool {
	return func(upper int) func(string) bool {
		return func(s string) bool {
			l := utf8.RuneCountInString(s)

			return lower <= l && l <= upper
		}
	}
}

func ShouldMatch(re *regexp.Regexp) func(string) bool {
	return func(s string) bool {
		return re.MatchString(s)
	}
}

func ShouldBeEmail(s string) bool {
	return ShouldMatch(regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`))(s)
}
