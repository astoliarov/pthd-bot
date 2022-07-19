package dao

import (
	"database/sql"
	_ "github.com/mattn/go-sqlite3"
	"teamkillbot/pkg/entities"
)

func OpenSQLite(pathToSQLite string) (*sql.DB, error) {
	return sql.Open("sqlite3", "./teamkillbot.sqlite")
}

type TeamKillLogDAO struct {
	db *sql.DB
}

func NewTeamKillLogDAO(db *sql.DB) *TeamKillLogDAO {
	return &TeamKillLogDAO{
		db: db,
	}
}

func (dao *TeamKillLogDAO) EnsureTable() error {
	createQuery := `
		CREATE TABLE IF NOT EXISTS team_kill_log
		(
			id          INTEGER  not null
				constraint id
					primary key autoincrement,
			killer      varchar  not null,
			victim      varchar  not null,
			happened_at datetime not null
		);
	`
	if _, err := dao.db.Exec(createQuery); err != nil {
		return err
	}

	return nil
}

func (dao *TeamKillLogDAO) Save(kill *entities.TeamKill) error {
	query := `
		INSERT INTO team_kill_log(killer, victim, happened_at) values(?,?,?)
	`
	stmt, prepareErr := dao.db.Prepare(query)
	if prepareErr != nil {
		return prepareErr
	}

	_, execErr := stmt.Exec(kill.Killer, kill.Victim, kill.HappenedAt)
	if execErr != nil {
		return execErr
	}

	return nil
}
