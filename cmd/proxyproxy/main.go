package main

import (
	"log/slog"
	"strings"

	"github.com/spf13/viper"

	"github.com/lukasdietrich/proxyproxy/internal/auto"
	"github.com/lukasdietrich/proxyproxy/internal/proxy"
	"github.com/lukasdietrich/proxyproxy/internal/server"
)

func init() {
	viper.SetDefault("verbose", false)
}

func main() {
	setupConfig()

	if err := run(); err != nil {
		slog.Error("fatal", slog.Any("err", err))
	}
}

func run() error {
	if viper.GetBool("verbose") {
		slog.SetLogLoggerLevel(slog.LevelDebug)
		slog.Debug("enabling debug logging")
	}

	if err := auto.ConfigureFromEnv(); err != nil {
		return err
	}

	handler, err := proxy.FromEnv()
	if err != nil {
		return err
	}

	listener := server.FromEnv(handler)

	slog.Info("starting http server", slog.String("addr", listener.Addr))
	return listener.ListenAndServe()
}

func setupConfig() {
	viper.SetTypeByDefaultValue(true)
	viper.AutomaticEnv()
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	viper.SetEnvPrefix("PROXYPROXY")
}
