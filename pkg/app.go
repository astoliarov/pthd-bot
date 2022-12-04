package pkg

import (
	"context"
	"github.com/getsentry/sentry-go"
	"github.com/jmoiron/sqlx"
	"github.com/rs/zerolog/log"
	"pthd-bot/pkg/connectors/telegram"
	"pthd-bot/pkg/dao"
	"pthd-bot/pkg/services"
	"time"
)

type Application struct {
	Config          *Config
	db              *sqlx.DB
	teamKillLogDAO  *dao.TeamKillLogDAO
	teamKillService *services.TeamKillService
	botKillService  *services.BotKillService
}

func NewApplication() *Application {
	setupLogs()

	config, configLoadErr := ConfitaConfigLoader()
	if configLoadErr != nil {
		log.Fatal().Err(configLoadErr).Msg("Can't load config")
	}

	db, openErr := dao.OpenSQLite(config.PathToSQLite)
	if openErr != nil {
		log.Fatal().
			Err(openErr).
			Str("path", config.PathToSQLite).
			Msg("Can't open sqlite")
	}

	teamKillDAO := dao.NewTeamKillLogDAO(db)
	botKillDAO := dao.NewBotKillLogDAO(db)

	responseSelector := &services.ResponseSelectorService{}
	teamKillService := services.NewTeamKillService(
		teamKillDAO,
		responseSelector,
	)

	botKillService := services.NewBotKillService(botKillDAO, responseSelector)

	app := &Application{
		Config:          config,
		db:              db,
		teamKillLogDAO:  teamKillDAO,
		teamKillService: teamKillService,
		botKillService:  botKillService,
	}

	return app
}

func (app *Application) RunBot(ctx context.Context) {
	bot, connectionErr := telegram.InitBot(app.Config.BotTGToken)
	if connectionErr != nil {
		log.Fatal().
			Err(connectionErr).
			Msg("Can't open telegram bot")
	}

	router := telegram.NewMessageRouter(
		bot,
		app.teamKillService,
		app.botKillService,
		app.Config.ChatID,
	)

	router.ListenToUpdates(ctx)

	defer sentry.Flush(2 * time.Second)
}

func (app *Application) MigrateUp() error {
	return dao.MigrateUp(app.db)
}
