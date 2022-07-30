package pkg

import (
	"context"
	"github.com/heetch/confita"
	"github.com/heetch/confita/backend/env"
)

type Config struct {
	PathToSQLite string
	BotTGToken   string
	SentryDSN    string
	ChatID       int64
}

type ConfitaConfig struct {
	PathToSQLite string `config:"tkbot_db_path"`
	BotTGToken   string `config:"tkbot_tg_token"`
	SentryDSN    string `config:"tkbot_sentry_dsn"`
	ChatID       int64  `config:"tkbot_chat_id"`
}

func (cfg *ConfitaConfig) ToConfig() *Config {
	return &Config{
		PathToSQLite: cfg.PathToSQLite,
		BotTGToken:   cfg.BotTGToken,
		SentryDSN:    cfg.SentryDSN,
		ChatID:       cfg.ChatID,
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
