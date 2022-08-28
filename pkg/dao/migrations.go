package dao

import (
	"embed"
	"github.com/jmoiron/sqlx"
	_ "github.com/mattn/go-sqlite3"
	"github.com/pressly/goose/v3"
)

//go:embed migrations/*.sql
var embeddedMigrations embed.FS //

func MigrateUp(db *sqlx.DB) error {
	setDialectErr := goose.SetDialect("sqlite3")
	if setDialectErr != nil {
		return setDialectErr
	}
	goose.SetBaseFS(embeddedMigrations)

	if err := goose.Up(db.DB, "migrations"); err != nil { //
		return err
	}

	return nil
}
