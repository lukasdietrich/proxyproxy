package server

import (
	"net/http"

	"github.com/spf13/viper"
)

func init() {
	viper.SetDefault("http.addr", ":8080")
	viper.SetDefault("http.timeout.read", "30s")
	viper.SetDefault("http.timeout.read.header", "10s")
	viper.SetDefault("http.timeout.write", "600s")
	viper.SetDefault("http.timeout.idle", "30s")
	viper.SetDefault("http.limit.header.bytes", "640k")
}

func FromEnv(handler http.Handler) *http.Server {
	return &http.Server{
		Addr:              viper.GetString("http.addr"),
		Handler:           handler,
		ReadTimeout:       viper.GetDuration("http.timeout.read"),
		ReadHeaderTimeout: viper.GetDuration("http.timeout.read.header"),
		WriteTimeout:      viper.GetDuration("http.timeout.write"),
		IdleTimeout:       viper.GetDuration("http.timeout.idle"),
		MaxHeaderBytes:    int(viper.GetSizeInBytes("http.limit.header.bytes")),
	}
}
