package pkg

import (
	"context"
	"github.com/heetch/confita"
	"github.com/heetch/confita/backend/env"
)

type Config struct {
	PathToSQLite string
	BotTGToken   string
}

type ConfitaConfig struct {
	PathToSQLite string `config:"teamkillbot_pathtosqlite"`
	BotTGToken   string `config:"teamkillbot_bottgtoken"`
}

func (cfg *ConfitaConfig) ToConfig() *Config {
	return &Config{
		PathToSQLite: cfg.PathToSQLite,
		BotTGToken:   cfg.BotTGToken,
	}
}

func ConfitaConfigLoader() (*Config, error) {

	cfg := ConfitaConfig{
		PathToSQLite: "./teamkillbot.sqlite",
	}

	loader := confita.NewLoader(
		env.NewBackend(),
	)

	err := loader.Load(context.Background(), &cfg)
	if err != nil {
		return nil, err
	}

	return cfg.ToConfig(), nil
}
