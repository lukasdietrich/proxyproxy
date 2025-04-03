package cache

import (
	"fmt"
	"log/slog"

	"github.com/spf13/viper"
)

func init() {
	viper.SetDefault("cache.duration.item", "30m")
	viper.SetDefault("cache.interval.gc", "15m")
}

func Func[K fmt.Stringer, V any](fn func(K) (V, error)) func(K) (V, error) {
	cache := newCache[K, V](
		viper.GetDuration("cache.duration.item"),
		viper.GetDuration("cache.interval.gc"),
	)

	return func(key K) (V, error) {
		cache.mu.Lock()
		defer cache.mu.Unlock()

		if value, ok := cache.get(key); ok {
			slog.Debug("return value from cache",
				slog.Any("key", key),
				slog.Any("value", value),
			)

			return value, nil
		}

		slog.Debug("value missing from cache", slog.Any("key", key))

		value, err := fn(key)
		if err == nil {
			cache.put(key, value)
		}

		return value, err
	}
}
