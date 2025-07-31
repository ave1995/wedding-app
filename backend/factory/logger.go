package factory

import (
	"log/slog"
	"os"
	"sync"
)

var (
	loggerOnce     sync.Once
	loggerInstance *slog.Logger
)

func getLogger() *slog.Logger {
	loggerOnce.Do(func() {
		handler := slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
			Level: slog.LevelInfo,
			// AddSource: true,
		})
		loggerInstance = slog.New(handler)
	})

	return loggerInstance
}
