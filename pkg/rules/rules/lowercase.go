package rules

import (
	"unicode"
	"unicode/utf8"
)

func CheckLowercaseStart(msg string) string {
	if len(msg) == 0 {
		return ""
	}
	r, _ := utf8.DecodeRuneInString(msg)
	if unicode.IsUpper(r) {
		return "log message should start with a lowercase letter"
	}
	return ""
}

func FixLowercaseStart(msg string) string {
	if len(msg) == 0 {
		return msg
	}
	r, size := utf8.DecodeRuneInString(msg)
	if unicode.IsUpper(r) {
		return string(unicode.ToLower(r)) + msg[size:]
	}
	return msg
}
