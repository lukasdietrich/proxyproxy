package pac

import (
	"fmt"
	"iter"
	"log/slog"
	"net/http"
	"net/url"
	"slices"
	"strings"
)

var (
	supportedUpstreamProxies = []string{"http", "https", "proxy"}
)

type Config struct {
	resolve resolveFunc
}

func New(url string) (*Config, error) {
	source, err := read(url)
	if err != nil {
		return nil, err
	}

	resolve, err := compile(source)
	if err != nil {
		return nil, err
	}

	config := Config{
		resolve: resolve,
	}

	return &config, nil
}

func (c *Config) Resolve(r *http.Request) (*url.URL, error) {
	requestUrl := stripUrl(r.URL)
	target := c.resolve(requestUrl.String(), requestUrl.Hostname())
	proxies := parseTargetWithFallback(target)

	for proxy := range proxies {
		slog.Debug("resolved upstead proxy", slog.String("uri", r.RequestURI), slog.Any("target", proxy))

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

func parseTargetWithFallback(targets *string) iter.Seq[*url.URL] {
	return func(yield func(*url.URL) bool) {
		if targets == nil {
			yield(nil)
			return
		}

		for target := range strings.SplitSeq(*targets, ";") {
			if !yield(parseTarget(strings.TrimSpace(target))) {
				return
			}
		}
	}
}

func parseTarget(target string) *url.URL {
	switch fields := strings.Fields(target); strings.ToUpper(fields[0]) {
	case "DIRECT":
		return nil

	default:
		return &url.URL{
			Scheme: strings.ToLower(fields[0]),
			Host:   fields[1],
		}
	}
}
