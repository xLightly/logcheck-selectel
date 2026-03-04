package rules

import (
	"strings"
)

var defaultSensitiveKeywords = []string{
	"password", "passwd", "pass_word",
	"secret", "api_key", "apikey", "api_secret",
	"token", "access_token", "refresh_token", "auth_token",
	"private_key", "privatekey",
	"credential", "credentials",
	"ssn", "social_security",
	"credit_card", "creditcard", "card_number",
	"cvv", "cvc",
	"pin_code", "pincode",
	"session_id", "sessionid",
	"cookie",
	"authorization",
}

func CheckNoSensitiveData(msg string, customPatterns []string) string {
	lower := strings.ToLower(msg)
	all := append(defaultSensitiveKeywords, customPatterns...)
	for _, kw := range all {
		if strings.Contains(lower, strings.ToLower(kw)) {
			return "log message may contain sensitive data (keyword: \"" + kw + "\")"
		}
	}
	return ""
}
