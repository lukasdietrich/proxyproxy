package proxy

import (
	"io"
	"log/slog"
	"net/http"
)

var (
	_ http.Handler = &Handler{}
)

var (
	proxyHeaders = []string{
		"Proxy-Connection",
		"Proxy-Authorization",
	}
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
	slog.Debug("serve http", slog.String("url", r.RequestURI), slog.String("method", r.Method))

	clearProxyHeaders(r)

	res, err := h.RoundTripper.RoundTrip(r)
	if err != nil {
		return err
	}

	return copyResponse(w, res)
}

func clearProxyHeaders(r *http.Request) {
	for _, header := range proxyHeaders {
		r.Header.Del(header)
	}
}

func copyResponse(w http.ResponseWriter, r *http.Response) error {
	clearHeader(w.Header())
	copyHeader(w.Header(), r.Header)
	w.WriteHeader(r.StatusCode)

	defer r.Body.Close()
	n, err := io.Copy(w, r.Body)

	slog.Debug("copied body", slog.Int64("bytes", n))

	return err
}

func copyHeader(dst, src http.Header) {
	for key, values := range src {
		dst[key] = append(values[:0:0], values...)
	}
}

func clearHeader(header http.Header) {
	for key := range header {
		header.Del(key)
	}
}
