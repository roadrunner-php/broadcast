package broadcast

import (
	"errors"
	"github.com/spiral/roadrunner/service"
)

// Config configures the broadcast extension.
type Config struct {
	// Path defines on this URL the middleware must be activated. Same path must be handled by underlying
	// application kernel to authorize the consumption. Optional.
	Path string

	// RedisConfig configures redis broker.
	Redis *RedisConfig
}

// Hydrate reads the configuration values from the source configuration.
func (c *Config) Hydrate(cfg service.Config) error {
	if err := cfg.Unmarshal(c); err != nil {
		return err
	}

	if c.Redis != nil {
		return c.Redis.isValid()
	}

	return nil
}

// InitDefaults enables in memory broadcast configuration.
func (c *Config) InitDefaults() error {
	return nil
}

// RedisConfig configures redis broker.
type RedisConfig struct {
	Addr     string
	Password string
	DB       int
}

func (r *RedisConfig) isValid() error {
	if r.Addr == "" {
		return errors.New("redis addr must be specified")
	}

	return nil
}
