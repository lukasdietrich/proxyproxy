package main

import (
	"flag"
	"log/slog"
	"net/http"

	"github.com/lukasdietrich/proxyproxy/internal/pac"
	"github.com/lukasdietrich/proxyproxy/internal/proxy"
)

func main() {
	if err := run(); err != nil {
		slog.Error("fatal", slog.Any("err", err))
	}
}

func run() error {
	var (
		addr    string
		pacUrl  string
		verbose bool
	)

	flag.StringVar(&addr, "addr", ":8080", "Address to listen on")
	flag.StringVar(&pacUrl, "pac", "file:/absolute/url/pac.js", "Absolute url to pac file")
	flag.BoolVar(&verbose, "verbose", false, "Enable verbose logging")
	flag.Parse()

	if verbose {
		slog.SetLogLoggerLevel(slog.LevelDebug)
	}

	upstream, err := pac.New(pacUrl)
	if err != nil {
		return err
	}

	handler := proxy.New(upstream)
	return http.ListenAndServe(addr, handler)
}
