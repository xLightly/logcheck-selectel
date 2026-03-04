package b

import "log/slog"

func example() {
	slog.Info("запуск сервера")      // want `log message should be in English`
	slog.Error("ошибка подключения") // want `log message should be in English`
	slog.Info("starting server")
	slog.Info("mixed english и русский") // want `log message should be in English`
}
