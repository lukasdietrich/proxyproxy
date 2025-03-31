package main

import (
	"log/slog"
	"net/http"

	"github.com/lukasdietrich/proxyproxy/internal/env"
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
		verbose = env.StringOrDefault(env.KEY_VERBOSE, "0") == "1"
		addr    = env.StringOrDefault(env.KEY_ADDR, ":8080")
	)

	if verbose {
		slog.SetLogLoggerLevel(slog.LevelDebug)
		slog.Debug("enabling debug logging")
	}

	upstream, err := pac.FromEnv()
	if err != nil {
		return err
	}

	handler := proxy.New(upstream)

	slog.Info("starting http server", slog.String("addr", addr))
	return http.ListenAndServe(addr, handler)
}
