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
	return establishTunnel(client, r.URL.Host)
}

func establishTunnel(client net.Conn, host string) error {
	slog.Debug("establishing tunnel to target", slog.String("host", host))

	target, err := net.Dial("tcp", host)
	if err != nil {
		return err
	}

	defer target.Close()

	fmt.Fprint(client, "HTTP/1.0 200 Connection established\r\n\r\n")

	var wg sync.WaitGroup
	copyAndClose(target.(*net.TCPConn), client.(*net.TCPConn), &wg)
	copyAndClose(client.(*net.TCPConn), target.(*net.TCPConn), &wg)

	wg.Wait()
	return nil
}

func copyAndClose(dst, src *net.TCPConn, wg *sync.WaitGroup) {
	wg.Add(1)

	go func() {
		defer dst.CloseWrite()
		defer src.CloseRead()

		io.Copy(dst, src)
		wg.Done()
	}()
}
