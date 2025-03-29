package main

import (
	"flag"
	"log/slog"
	"net/http"

	"github.com/lukasdietrich/proxyproxy/internal/proxy"
)

func main() {
	var (
		addr    string
		verbose bool
	)

	flag.StringVar(&addr, "addr", ":8080", "Address to listen on")
	flag.BoolVar(&verbose, "verbose", false, "Enable verbose logging")
	flag.Parse()

	if verbose {
		slog.SetLogLoggerLevel(slog.LevelDebug)
	}

	h := proxy.Handler{}
	http.ListenAndServe(addr, &h)
}
