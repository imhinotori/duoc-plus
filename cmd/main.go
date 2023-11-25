package main

import (
	"github.com/charmbracelet/log"
	"github.com/imhinotori/duoc-plus/internal/config"
	"github.com/imhinotori/duoc-plus/internal/server"
	"github.com/knadh/koanf/parsers/toml"
	"github.com/knadh/koanf/providers/file"
	"github.com/knadh/koanf/v2"
	"os"
)

func main() {
	k := koanf.New(".")

	var basePath string

	if _, err := os.Stat("/.dockerenv"); err == nil {
		log.Debug("Running inside Docker container")
		basePath = "/data/"
	} else {
		log.Debug("Running outside Docker container")
		basePath = "data/"
	}

	if err := k.Load(file.Provider(basePath+"config.toml"), toml.Parser()); err != nil {
		log.Fatal("Failed reading configuration", "err", err)
	}

	srv, err := server.New(config.FromKoanf(k))
	if err != nil {
		log.Fatal("Failed creating Server", "err", err)
	}

	if err = srv.Run(); err != nil {
		log.Fatal("Failed running Server", "err", err)
	}

}
