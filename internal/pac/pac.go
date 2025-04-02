package pac

import (
	"fmt"
	"iter"
	"log/slog"
	"net/http"
	"net/url"
	"slices"
	"strings"
	"time"

	"github.com/spf13/viper"
)

func init() {
	viper.SetDefault("pac.url", "")
}

var (
	supportedUpstreamProxies = []string{"http", "https", "proxy"}
)

type Config struct {
	resolve resolveFunc
}

func FromEnv() (*Config, error) {
	url := viper.GetString("pac.url")
	if url == "" {
		slog.Info("no pac url provided. defaulting direct connections")
		return Direct(), nil
	}

	slog.Info("configuring upstream proxies using pac", slog.String("url", url))
	return FromUrl(url)
}

func FromUrl(url string) (*Config, error) {
	source, err := read(url)
	if err != nil {
		return nil, err
	}

	return FromSource(source)
}

func FromSource(source []byte) (*Config, error) {
	resolve, err := compile(source)
	if err != nil {
		return nil, err
	}

	config := Config{
		resolve: resolve,
	}

	return &config, nil
}

func Direct() *Config {
	return &Config{
		resolve: func(string, string) *string {
			return nil
		},
	}
}

func (c *Config) Resolve(r *http.Request) (*url.URL, error) {
	t0 := time.Now()

	requestUrl := stripUrl(r.URL)
	target := c.resolve(requestUrl.String(), requestUrl.Hostname())
	proxies := parseTargetWithFallback(target)

	for proxy, err := range proxies {
		if err != nil {
			return nil, err
		}

		slog.Debug("resolved upsteam proxy",
			slog.String("uri", r.RequestURI),
			slog.Any("target", proxy),
			slog.Duration("t", time.Since(t0)),
		)

		if proxy != nil && !slices.Contains(supportedUpstreamProxies, proxy.Scheme) {
			slog.Warn("skipping unsupported upstream proxy", slog.Any("target", proxy))
			continue
		}

		return proxy, nil
	}

	return nil, fmt.Errorf("could not resolve valid upstream proxy")
}

func stripUrl(u *url.URL) *url.URL {
	return &url.URL{
		Scheme: u.Scheme,
		Host:   u.Host,
	}
}

func parseTargetWithFallback(targets *string) iter.Seq2[*url.URL, error] {
	return func(yield func(*url.URL, error) bool) {
		if targets == nil {
			yield(nil, nil)
			return
		}

		for target := range strings.SplitSeq(*targets, ";") {
			if !yield(parseTarget(strings.TrimSpace(target))) {
				return
			}
		}
	}
}

func parseTarget(target string) (*url.URL, error) {
	fields := strings.Fields(target)

	switch len(fields) {
	case 1:
		if fields[0] == "DIRECT" {
			return nil, nil
		}

		return nil, fmt.Errorf("invalid target %q, expected DIRECT", fields[0])

	case 2:
		target := url.URL{
			Scheme: strings.ToLower(fields[0]),
			Host:   fields[1],
		}

		return &target, nil

	default:
		return nil, fmt.Errorf("invalid target %+v, expected DIRECT or \"${TYPE} ${HOST}\"", fields)

	}
}
