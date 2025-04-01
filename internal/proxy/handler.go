package proxy

import (
	"log/slog"
	"net/http"

	"github.com/rs/xid"

	"github.com/lukasdietrich/proxyproxy/internal/pac"
)

var (
	_ http.Handler = &Handler{}
)

type Handler struct {
	rt *http.Transport
}

func FromEnv() (*Handler, error) {
	pac, err := pac.FromEnv()
	if err != nil {
		return nil, err
	}

	return New(pac), nil
}

func New(upstream *pac.Config) *Handler {
	return &Handler{
		rt: &http.Transport{
			Proxy: upstream.Resolve,
		},
	}
}

func (h *Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	log := slog.With(slog.Group("request",
		slog.Any("id", xid.New()),
		slog.String("method", r.Method),
		slog.Any("url", r.URL),
	))

	if err := h.handle(log, w, r); err != nil {
		log.Warn("could not proxy request", slog.Any("err", err))
		http.Error(w, http.StatusText(http.StatusBadGateway), http.StatusBadGateway)
	}
}

func (h *Handler) handle(log *slog.Logger, w http.ResponseWriter, r *http.Request) error {
	if r.Method == http.MethodConnect {
		return h.proxyHttps(log, w, r)
	}

	return h.proxyHttp(log, w, r)
}
