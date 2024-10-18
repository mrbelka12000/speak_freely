package v1

import "log/slog"

type opt func(h *Handler)

func WithLogger(log *slog.Logger) opt {
	return func(h *Handler) {
		h.log = log
	}
}
