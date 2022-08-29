package dao

import "github.com/jmoiron/sqlx"

func OpenSQLite(pathToSQLite string) (*sqlx.DB, error) {
	return sqlx.Open("sqlite3", pathToSQLite)
}
