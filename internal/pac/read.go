package pac

import (
	"crypto/tls"
	"fmt"
	"io"
	"net/http"
	"os"
)

func read(url string) ([]byte, error) {
	r, err := client().Get(url)
	if err != nil {
		return nil, err
	}

	//nolint:errcheck
	defer r.Body.Close()

	if r.StatusCode != 200 {
		return nil, fmt.Errorf("could not read pac url: %s", r.Status)
	}

	return io.ReadAll(r.Body)
}

func client() *http.Client {
	transport := http.Transport{
		Proxy: nil,
		TLSClientConfig: &tls.Config{
			InsecureSkipVerify: true,
		},
	}
	transport.RegisterProtocol("file", http.NewFileTransportFS(os.DirFS("/")))

	return &http.Client{
		Transport: &transport,
	}
}
