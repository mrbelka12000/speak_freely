package usecase

import "log/slog"

type opt func(uc *UseCase)

// WithLogger set custom logger
func WithLogger(log *slog.Logger) opt {
	return func(uc *UseCase) {
		uc.log = log
	}
}
