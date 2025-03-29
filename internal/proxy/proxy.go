package proxy

import (
	"log/slog"
	"net/http"
)

var (
	_ http.Handler = &Handler{}
)

type Handler struct {
}

func (h *Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	slog.Debug("serve http", slog.String("url", r.RequestURI), slog.String("method", r.Method))
}
