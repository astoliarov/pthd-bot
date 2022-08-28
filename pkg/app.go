package pkg

import (
	"context"
	"fmt"
	"github.com/getsentry/sentry-go"
	"github.com/jmoiron/sqlx"
	"log"
	"teamkillbot/pkg/connectors/telegram"
	"teamkillbot/pkg/dao"
	"teamkillbot/pkg/services"
	"time"
)

type Application struct {
	Config          *Config
	db              *sqlx.DB
	teamKillLogDAO  *dao.TeamKillLogDAO
	teamKillService *services.TeamKillService
}

func NewApplication() *Application {
	config, configLoadErr := ConfitaConfigLoader()
	if configLoadErr != nil {
		log.Fatalf(fmt.Sprintf("Issue while loading config: %s", configLoadErr))
	}

	db, openErr := dao.OpenSQLite(config.PathToSQLite)
	if openErr != nil {
		log.Fatalf("Cannot open sqlite: %s", openErr)
	}

	teamKillDAO := dao.NewTeamKillLogDAO(db)

	responseSelector := &services.ResponseSelectorService{}
	teamKillService := services.NewTeamKillService(
		teamKillDAO,
		responseSelector,
	)

	app := &Application{
		Config:          config,
		db:              db,
		teamKillLogDAO:  teamKillDAO,
		teamKillService: teamKillService,
	}

	return app
}

func (app *Application) RunBot(ctx context.Context) {
	bot, connectionErr := telegram.InitBot(app.Config.BotTGToken)
	if connectionErr != nil {
		log.Fatalf("Cannot open telegram bot: %s", connectionErr)
	}

	router := telegram.NewMessageRouter(
		bot,
		app.teamKillService,
		app.Config.ChatID,
	)

	router.ListenToUpdates(ctx)

	defer sentry.Flush(2 * time.Second)
}

func (app *Application) MigrateUp() error {
	return dao.MigrateUp(app.db)
}
