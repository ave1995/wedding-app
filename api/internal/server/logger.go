package server

import (
	"log/slog"
	"os"
)

func InitLogger() {
	handler := slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.LevelInfo,
		// AddSource: true,
	})
	slog.SetDefault(slog.New(handler))
}

func ErrAttr(err error) slog.Attr {
	return slog.Any("error", err)
}
