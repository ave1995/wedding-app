package factory

import (
	"log/slog"
	"os"
	"sync"
	"time"
)

var (
	loggerOnce     sync.Once
	loggerInstance *slog.Logger
)

func getLogger() *slog.Logger {
	loggerOnce.Do(func() {
		handler := slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
			// AddSource: true,
			Level: slog.LevelInfo,
			ReplaceAttr: func(groups []string, a slog.Attr) slog.Attr {
				if a.Key == slog.TimeKey {
					if t, ok := a.Value.Any().(time.Time); ok {
						a.Value = slog.StringValue(t.Format(time.DateTime))
					}
				}
				return a
			},
		})
		loggerInstance = slog.New(handler)
	})

	return loggerInstance
}
