package pac

import (
	"crypto/tls"
	"io"
	"net/http"
	"os"
)

func read(url string) ([]byte, error) {
	r, err := client().Get(url)
	if err != nil {
		return nil, err
	}

	defer r.Body.Close()
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
