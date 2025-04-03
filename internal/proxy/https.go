package proxy

import (
	"fmt"
	"io"
	"log/slog"
	"net"
	"net/http"
	"sync"
	"time"
)

func (h *Handler) proxyHttps(log *slog.Logger, w http.ResponseWriter, r *http.Request) error {
	hijacker, ok := w.(http.Hijacker)
	if !ok {
		return fmt.Errorf("could not hijack response writer")
	}

	log.Debug("hijacking response writer")
	client, _, err := hijacker.Hijack()
	if err != nil {
		return err
	}

	//nolint:errcheck
	defer client.Close()

	return h.establishTunnel(log, client, r)
}

func (h *Handler) establishTunnel(log *slog.Logger, client net.Conn, r *http.Request) error {
	t0 := time.Now()

	upstream, err := h.rt.Proxy(r)
	if err != nil {
		return err
	}

	host := r.URL.Host
	if upstream != nil {
		log = log.With(slog.Any("upstream", upstream))
		log.Debug("establishing tunnel through another proxy")

		host = upstream.Host
	} else {
		log.Debug("establishing tunnel to target directly")
	}

	target, err := net.Dial("tcp", host)
	if err != nil {
		return err
	}

	//nolint:errcheck
	defer target.Close()

	if upstream != nil {
		log.Debug("forwarding original request")

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
	copyAndClose(log, &wg, target, client)
	copyAndClose(log, &wg, client, target)
	wg.Wait()

	log.Debug("tunnel closed", slog.Duration("duration", time.Since(t0)))
	return nil
}

func copyAndClose(log *slog.Logger, wg *sync.WaitGroup, dst, src net.Conn) {
	wg.Add(1)

	go func() {
		if _, err := io.Copy(dst, src); err != nil {
			log.Warn("error while tunneling data", slog.Any("err", err))
		}

		if err := src.(interface{ CloseRead() error }).CloseRead(); err != nil {
			log.Debug("could not close source reader", slog.Any("err", err))
		}

		if err := dst.(interface{ CloseWrite() error }).CloseWrite(); err != nil {
			log.Debug("could not close target writer", slog.Any("err", err))
		}

		wg.Done()
	}()
}
