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

// @title Duoc Plus API
// @version 1.0
// @description Duoc Plus, is a REST API that allows you to access your grades, schedule and attendance from DuocUC.
// @termsOfService https://www.duoc.cl/politica-privacidad/
// @securityDefinitions.bearerAuth Bearer

// @contact.name Matias "Hinotori" Canovas
// @contact.url https://github.com/imhinotori/
// @contact.email hello@hinotori.moe

// @license.name MIT
// @license.url https://opensource.org/license/mit/

// @host api-duoc.hinotori.moe
// @BasePath /
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
