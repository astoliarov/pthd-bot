package main

import (
	"fmt"
	"log"
	"teamkillbot/pkg"
	"teamkillbot/pkg/connectors/telegram"
	"teamkillbot/pkg/dao"
	"teamkillbot/pkg/services"
)

func main() {
	config, configLoadErr := pkg.ConfitaConfigLoader()
	if configLoadErr != nil {
		log.Fatalf(fmt.Sprintf("Issue while loading config: %s", configLoadErr))
	}

	db, openErr := dao.OpenSQLite(config.PathToSQLite)
	if openErr != nil {
		log.Fatalf("Cannot open sqlite: %s", openErr)
	}

	teamKillDAO := dao.NewTeamKillLogDAO(db)
	createErr := teamKillDAO.EnsureTable()
	if createErr != nil {
		log.Fatalf("Cannot create table: %s", createErr)
	}

	responseSelector := &services.ResponseSelectorService{}
	teamKillService := services.NewTeamKillService(
		teamKillDAO,
		responseSelector,
	)

	bot, connectionErr := telegram.InitBot(config.BotTGToken)
	if connectionErr != nil {
		log.Fatalf("Cannot open telegram bot: %s", connectionErr)
	}

	router := telegram.NewMessageRouter(
		bot,
		teamKillService,
	)

	router.ListenToUpdates()
}
