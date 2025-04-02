package proxy

import (
	"io"
	"log/slog"
	"net/http"
	"slices"
)

var (
	proxyHeaders = []string{
		"Proxy-Connection",
		"Proxy-Authorization",
	}
)

func (h *Handler) proxyHttp(log *slog.Logger, w http.ResponseWriter, r *http.Request) error {
	log.Debug("clearing proxy headers")
	clearProxyHeaders(r)

	log.Debug("forwarding request via http")
	res, err := h.rt.RoundTrip(r)
	if err != nil {
		return err
	}

	//nolint:errcheck
	defer res.Body.Close()

	log.Debug("copying response")
	return copyResponse(log, w, res)
}

func clearProxyHeaders(r *http.Request) {
	for _, header := range proxyHeaders {
		r.Header.Del(header)
	}
}

func copyResponse(log *slog.Logger, w http.ResponseWriter, r *http.Response) error {
	clearHeader(w.Header())
	copyHeader(w.Header(), r.Header)

	w.WriteHeader(r.StatusCode)

	n, err := io.Copy(w, r.Body)
	log.Debug("copied body", slog.Int64("bytes", n))

	return err
}

func copyHeader(dst, src http.Header) {
	for key, values := range src {
		dst[key] = slices.Clone(values)
	}
}

func clearHeader(header http.Header) {
	for key := range header {
		header.Del(key)
	}
}
