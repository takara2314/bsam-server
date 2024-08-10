package logging

import (
	"log/slog"
	"os"

	"github.com/m-mizutani/clog"
)

func InitSlog() {
	// 本番環境はJSON形式でログを出力する
	if os.Getenv("ENVIRONMENT") == "production" {
		logger := slog.New(
			slog.NewJSONHandler(os.Stdout, nil),
		)
		slog.SetDefault(logger)
		return
	}

	// 開発環境はカラー付きのログを出力する
	handler := clog.New(
		clog.WithColor(true),
		clog.WithSource(true),
	)
	logger := slog.New(handler)
	slog.SetDefault(logger)
}
