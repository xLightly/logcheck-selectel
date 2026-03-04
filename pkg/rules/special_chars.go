package rules

import (
	"strings"
	"unicode"
)

var allowedPunctuation = map[rune]bool{
	'.': true, ',': true, ':': true, ';': true,
	'-': true, '_': true, '/': true, '\\': true,
	'=': true, '(': true, ')': true, '[': true, ']': true,
	'{': true, '}': true, '\'': true, '"': true,
	'%': true, '@': true, '#': true, '&': true, '+': true,
	'<': true, '>': true, '|': true, '^': true, '~': true,
	'`': true, '$': true,
}
var forbiddenRepeated = map[rune]bool{
	'!': true, '?': true, '.': true,
}

func CheckNoSpecialChars(msg string, extraAllowed string) string {
	extraSet := make(map[rune]bool)
	for _, r := range extraAllowed {
		extraSet[r] = true
	}
	for _, r := range msg {
		if extraSet[r] {
			continue
		}
		if isEmoji(r) {
			return "log message should not contain emoji"
		}
	}
	for _, ch := range []rune{'!', '?', '.'} {
		pattern := string([]rune{ch, ch})
		if strings.Contains(msg, pattern) {
			return "log message should not contain repeated special characters"
		}
	}
	for _, r := range msg {
		if extraSet[r] {
			continue
		}
		if r == '!' || r == '?' {
			return "log message should not contain special characters like '!' or '?'"
		}
		if !unicode.IsLetter(r) && !unicode.IsDigit(r) && !unicode.IsSpace(r) && !allowedPunctuation[r] && r > 31 {
			if !isEmoji(r) {
				return "log message should not contain special characters"
			}
		}
	}

	return ""
}

func isEmoji(r rune) bool {
	if r < 128 {
		return false
	}
	return unicode.Is(unicode.So, r) ||
		unicode.Is(unicode.Sk, r) ||
		(r >= 0x1F600 && r <= 0x1F64F) ||
		(r >= 0x1F300 && r <= 0x1F5FF) ||
		(r >= 0x1F680 && r <= 0x1F6FF) ||
		(r >= 0x1F900 && r <= 0x1F9FF) ||
		(r >= 0x2600 && r <= 0x26FF) ||
		(r >= 0x2700 && r <= 0x27BF) ||
		(r >= 0xFE00 && r <= 0xFE0F) ||
		(r >= 0x200D && r <= 0x200D)
}

func FixNoSpecialChars(msg string) string {
	var b strings.Builder
	prev := rune(0)
	for _, r := range msg {
		if isEmoji(r) {
			continue
		}
		if r == '!' || r == '?' {
			continue
		}
		if forbiddenRepeated[r] && r == prev {
			continue
		}
		b.WriteRune(r)
		prev = r
	}
	return strings.TrimSpace(b.String())
}
