package logger

import (
	"log/slog"
	"os"

	"github.com/k1ender/task-master-go/internal/config"
)

func MustInit(cfg *config.Config) *slog.Logger {
	var logger *slog.Logger

	if cfg.ENV == config.EnvProd {
		logger = slog.New(
			slog.NewJSONHandler(
				os.Stdout,
				&slog.HandlerOptions{
					Level: slog.LevelInfo,
				},
			),
		)
	} else {
		logger = slog.New(slog.NewTextHandler(
			os.Stdout,
			&slog.HandlerOptions{
				Level: slog.LevelDebug,
			},
		))
	}

	return logger
}
