// Package main ...
package main

import (
	"flag"

	"github.com/SerjRamone/reposter-bot/config"
	"github.com/SerjRamone/reposter-bot/internal/app"
)

func main() {
	configFile := flag.String("p", ".yml", "Path to config.yml file")
	flag.Parse()

	cfg := config.Get(*configFile)

	app := app.New(cfg)
	app.Run()
}
