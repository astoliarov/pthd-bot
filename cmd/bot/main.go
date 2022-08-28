package main

import (
	"context"
	"github.com/getsentry/sentry-go"
	"log"
	"os"
	"os/signal"
	"syscall"
	"teamkillbot/pkg"
)

func main() {
	app := pkg.NewApplication()

	if app.Config.SentryDSN != "" {
		sentry.Init(sentry.ClientOptions{
			Dsn: app.Config.SentryDSN,
		})
	}

	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)
	exitChan := make(chan os.Signal, 1)
	signal.Notify(exitChan, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		<-exitChan
		cancel()
		log.Println("Received break signal")
	}()

	migrateErr := app.MigrateUp()
	if migrateErr != nil {
		log.Fatalln(migrateErr)
	}
	app.RunBot(ctx)
}
