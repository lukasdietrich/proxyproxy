package proxy

import (
	"log/slog"
	"net/http"
	"net/url"

	"github.com/rs/xid"

	"github.com/lukasdietrich/proxyproxy/internal/cache"
	"github.com/lukasdietrich/proxyproxy/internal/pac"
)

var (
	_ http.Handler = &Handler{}
)

type resolveRequestProxyFunc func(*http.Request) (*url.URL, error)
type resolveUrlProxyFunc func(*url.URL) (*url.URL, error)

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
			Proxy: wrapResolveRequestProxyFunc(cache.Func(upstream.Resolve)),
		},
	}
}

func wrapResolveRequestProxyFunc(resolve resolveUrlProxyFunc) resolveRequestProxyFunc {
	return func(r *http.Request) (*url.URL, error) {
		strippedUrl := stripUrl(r.URL)
		return resolve(strippedUrl)
	}
}

func stripUrl(u *url.URL) *url.URL {
	return &url.URL{
		Scheme: u.Scheme,
		Host:   u.Host,
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
