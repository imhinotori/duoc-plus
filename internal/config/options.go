package config

import (
	"github.com/charmbracelet/log"
	"github.com/knadh/koanf/v2"
)

type Option func(cfg *Config) *Config

func FromKoanf(k *koanf.Koanf) Option {
	return func(cfg *Config) *Config {
		err := k.Unmarshal("", cfg)
		if err != nil {
			log.Fatal("parsing failed", "err", err)
		}

		return cfg
	}

}
