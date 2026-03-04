package a

import "log/slog"

func example() {
	slog.Info("Starting server on port 8080") // want `log message should start with a lowercase letter`
	slog.Info("starting server on port 8080")
	slog.Error("Failed to connect") // want `log message should start with a lowercase letter`
	slog.Error("failed to connect")
}
