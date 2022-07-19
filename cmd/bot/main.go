package main

import (
	"fmt"
	"log"
	"teamkillbot/pkg"
	"teamkillbot/pkg/connectors"
	"teamkillbot/pkg/dao"
	"teamkillbot/pkg/services"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func main() {
	config, configLoadErr := pkg.ConfitaConfigLoader()
	if configLoadErr != nil {
		log.Fatalf(fmt.Sprintf("Issue while loading config: %s", configLoadErr))
	}

	db, openErr := dao.OpenSQLite(config.PathToSQLite)
	if openErr != nil {
		log.Fatalf(fmt.Sprintf("%s", openErr))
	}

	teamKillDAO := dao.NewTeamKillLogDAO(db)
	createErr := teamKillDAO.EnsureTable()
	if createErr != nil {
		log.Fatalf(fmt.Sprintf("%s", createErr))
	}

	responseSelector := &services.ResponseSelectorService{}
	teamKillService := services.NewTeamKillService(
		teamKillDAO,
		responseSelector,
	)

	bot, err := tgbotapi.NewBotAPI(config.BotTGToken)
	if err != nil {
		log.Panic(err)
	}

	bot.Debug = true

	log.Printf("Authorized on account %s", bot.Self.UserName)

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	router := connectors.NewMessageRouter(
		bot,
		teamKillService,
	)

	router.ListenToUpdates()
}
