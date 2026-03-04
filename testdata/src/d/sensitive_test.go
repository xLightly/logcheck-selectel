package d

import "log/slog"

func example() {
	password := "secret123"
	apiKey := "abc"
	token := "tok"
	_ = password
	_ = apiKey
	_ = token
	slog.Info("user password: " + password) // want `log message may contain sensitive data`
	slog.Debug("api_key=" + apiKey)         // want `log message may contain sensitive data`
	slog.Info("token: " + token)            // want `log message may contain sensitive data`
	slog.Info("user authenticated successfully")
	slog.Debug("api request completed")
}
