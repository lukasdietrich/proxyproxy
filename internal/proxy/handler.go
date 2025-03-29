package proxy

import (
	"log/slog"
	"net/http"
)

var (
	_ http.Handler = &Handler{}
)

type Handler struct {
	RoundTripper http.RoundTripper
}

func (h *Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if err := h.handle(w, r); err != nil {
		slog.Warn("could not proxy request", slog.Any("err", err))

		if err, ok := err.(Error); ok {
			err.write(w)
			return
		}

		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
	}
}

func (h *Handler) handle(w http.ResponseWriter, r *http.Request) error {
	if r.Method == http.MethodConnect {
		return h.proxyHttps(w, r)
	}

	return h.proxyHttp(w, r)
}
