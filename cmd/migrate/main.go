package main

import (
	"github.com/getsentry/sentry-go"
	"log"
	"teamkillbot/pkg"
)

func main() {
	app := pkg.NewApplication()

	if app.Config.SentryDSN != "" {
		sentry.Init(sentry.ClientOptions{
			Dsn: app.Config.SentryDSN,
		})
	}

	migrateErr := app.MigrateUp()
	if migrateErr != nil {
		log.Fatalln(migrateErr)
	}
}
