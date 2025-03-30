package proxy

import (
	"fmt"
	"io"
	"log/slog"
	"net"
	"net/http"
	"sync"
)

func (h *Handler) proxyHttps(w http.ResponseWriter, r *http.Request) error {
	hijacker, ok := w.(http.Hijacker)
	if !ok {
		return Error{Status: http.StatusInternalServerError, Text: "cannot hijack response writer"}
	}

	client, _, err := hijacker.Hijack()
	if err != nil {
		return err
	}

	defer client.Close()
	return h.establishTunnel(client, r)
}

func (h *Handler) establishTunnel(client net.Conn, r *http.Request) error {
	slog.Debug("establishing tunnel to target", slog.String("host", r.URL.Host))

	upstream, err := h.rt.Proxy(r)
	if err != nil {
		return err
	}

	host := r.URL.Host
	if upstream != nil {
		slog.Debug("establishing tunnel through another proxy", slog.Any("upstream", upstream))
		host = upstream.Host
	}

	target, err := net.Dial("tcp", host)
	if err != nil {
		return err
	}

	defer target.Close()

	if upstream != nil {
		// forward request to the upstream proxy
		if err := r.Write(target); err != nil {
			return err
		}
	} else {
		// only send a response if we connect to the target directly, otherwise the upstream proxy
		// sends the respose
		if _, err := fmt.Fprint(client, "HTTP/1.0 200 Connection established\r\n\r\n"); err != nil {
			return err
		}
	}

	var wg sync.WaitGroup
	copyAndClose(target, client, &wg)
	copyAndClose(client, target, &wg)

	wg.Wait()
	return nil
}

func copyAndClose(dst, src net.Conn, wg *sync.WaitGroup) {
	wg.Add(1)

	go func() {
		defer dst.(interface{ CloseWrite() error }).CloseWrite()
		defer src.(interface{ CloseRead() error }).CloseRead()

		io.Copy(dst, src)
		wg.Done()
	}()
}
