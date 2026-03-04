package rules

import (
	"testing"
)

func TestCheckLowercaseStart(t *testing.T) {
	tests := []struct {
		msg  string
		want bool
	}{
		{"starting server", false},
		{"Starting server", true},
		{"123 numeric start", false},
		{"", false},
		{"failed to connect", false},
		{"Failed to connect", true},
		{"сервер запущен", false},
		{"Сервер запущен", true},
	}
	for _, tt := range tests {
		got := CheckLowercaseStart(tt.msg)
		if (got != "") != tt.want {
			t.Errorf("CheckLowercaseStart(%q) = %q, wantDiag=%v", tt.msg, got, tt.want)
		}
	}
}

func TestFixLowercaseStart(t *testing.T) {
	tests := []struct {
		in, want string
	}{
		{"Starting server", "starting server"},
		{"starting server", "starting server"},
		{"", ""},
		{"ABC", "aBC"},
	}
	for _, tt := range tests {
		got := FixLowercaseStart(tt.in)
		if got != tt.want {
			t.Errorf("FixLowercaseStart(%q) = %q, want %q", tt.in, got, tt.want)
		}
	}
}

func TestCheckEnglishOnly(t *testing.T) {
	tests := []struct {
		msg  string
		want bool
	}{
		{"starting server", false},
		{"port 8080", false},
		{"запуск сервера", true},
		{"日本語テスト", true},
		{"hello world 123", false},
		{"", false},
		{"mixed english и русский", true},
	}
	for _, tt := range tests {
		got := CheckEnglishOnly(tt.msg)
		if (got != "") != tt.want {
			t.Errorf("CheckEnglishOnly(%q) = %q, wantDiag=%v", tt.msg, got, tt.want)
		}
	}
}

func TestCheckNoSpecialChars(t *testing.T) {
	tests := []struct {
		msg  string
		want bool
	}{
		{"server started", false},
		{"server started! 🚀", true},
		{"connection failed!!!", true},
		{"warning: something went wrong...", true},
		{"request completed", false},
		{"error code=42", false},
		{"hello!", true},
		{"value is 100%", false},
	}
	for _, tt := range tests {
		got := CheckNoSpecialChars(tt.msg, "")
		if (got != "") != tt.want {
			t.Errorf("CheckNoSpecialChars(%q) = %q, wantDiag=%v", tt.msg, got, tt.want)
		}
	}
}

func TestFixNoSpecialChars(t *testing.T) {
	tests := []struct {
		in, want string
	}{
		{"server started! 🚀", "server started"},
		{"failed!!!", "failed"},
		{"something went wrong...", "something went wrong."},
		{"clean message", "clean message"},
	}
	for _, tt := range tests {
		got := FixNoSpecialChars(tt.in)
		if got != tt.want {
			t.Errorf("FixNoSpecialChars(%q) = %q, want %q", tt.in, got, tt.want)
		}
	}
}

func TestCheckNoSensitiveData(t *testing.T) {
	tests := []struct {
		msg    string
		custom []string
		want   bool
	}{
		{"user authenticated successfully", nil, false},
		{"user password: 12345", nil, true},
		{"api_key=abc123", nil, true},
		{"token: xyz", nil, true},
		{"request completed", nil, false},
		{"credit_card number logged", nil, true},
		{"session_id=abcdef", nil, true},
		{"customer ssn stored", nil, true},
		{"my_custom_field=value", []string{"my_custom_field"}, true},
		{"my_custom_field=value", nil, false},
	}
	for _, tt := range tests {
		got := CheckNoSensitiveData(tt.msg, tt.custom)
		if (got != "") != tt.want {
			t.Errorf("CheckNoSensitiveData(%q, %v) = %q, wantDiag=%v", tt.msg, tt.custom, got, tt.want)
		}
	}
}
