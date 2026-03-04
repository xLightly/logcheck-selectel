package rules

import (
	"unicode"
)

func CheckEnglishOnly(msg string) string {
	for _, r := range msg {
		if unicode.IsLetter(r) && r > 127 {
			return "log message should be in English (contains non-ASCII letters)"
		}
	}
	return ""
}
