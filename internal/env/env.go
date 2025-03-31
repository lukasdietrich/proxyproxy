package env

import (
	"os"
)

const (
	KEY_PAC_URL = "PROXYPROXY_PAC_URL"
	KEY_ADDR    = "PROXYPROXY_ADDR"
	KEY_VERBOSE = "PROXYPROXY_VERBOSE"
)

func String(key string) string {
	return StringOrDefault(key, "")
}

func StringOrDefault(key, defaultValue string) string {
	value, ok := os.LookupEnv(key)
	if !ok {
		return defaultValue
	}

	return value
}
