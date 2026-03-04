package c

import "log/slog"

func example() {
	slog.Info("server started! 🚀")                // want `log message should not contain`
	slog.Error("connection failed!!!")            // want `log message should not contain`
	slog.Warn("warning: something went wrong...") // want `log message should not contain`
	slog.Info("server started")
	slog.Error("connection failed")
}
