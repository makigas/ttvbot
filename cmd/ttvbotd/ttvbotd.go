package main

import (
	"errors"
	"flag"

	"go.uber.org/fx"
	"gopkg.makigas.es/ttvbot/domain/config"
	"gopkg.makigas.es/ttvbot/hooks"
	"gopkg.makigas.es/ttvbot/httpd/server"
	"gopkg.makigas.es/ttvbot/integrations/discord"
	"gopkg.makigas.es/ttvbot/integrations/ttv/eventsub"
	"gopkg.makigas.es/ttvbot/integrations/ttv/helixapi"
	"gopkg.makigas.es/ttvbot/persistence/redis"
)

var (
	ErrLoadConfig = errors.New("cannot load configuration")

	mainConfig config.Config

	flagConfig string
)

func init() {
	// Parse command line flags.
	flag.StringVar(&flagConfig, "config", "", "the config file to read")
	flag.Parse()
}

func main() {
	// Load the configuration and provide it as an FX module.
	if err := mainConfig.ReadFromFile(flagConfig); err != nil {
		panic(errors.Join(ErrLoadConfig, err))
	}
	configFx := fx.Module("Config", fx.Provide(
		func() *config.Config {
			return &mainConfig
		},
	))

	fx.New(
		configFx,

		// Main dependencies
		server.Module,
		redis.Module,
		discord.Module,
		eventsub.Module,
		helixapi.Module,

		// Register hooks
		fx.Invoke(hooks.StreamNotifications),
	).Run()
}
