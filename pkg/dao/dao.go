package dao

import (
	"github.com/jmoiron/sqlx"
	_ "github.com/mattn/go-sqlite3"
	"teamkillbot/pkg/entities"
)

func OpenSQLite(pathToSQLite string) (*sqlx.DB, error) {
	return sqlx.Open("sqlite3", pathToSQLite)
}

type TeamKillLogDAO struct {
	db *sqlx.DB
}

func NewTeamKillLogDAO(db *sqlx.DB) *TeamKillLogDAO {
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

func (dao *TeamKillLogDAO) GetTopKillers() ([]*entities.TopKillerLog, error) {
	query := `
		SELECT killer as name, count(*) as kill_count
		FROM team_kill_log
		GROUP BY killer
		ORDER BY kill_count desc;
	`
	rows, queryErr := dao.db.Queryx(query)
	if queryErr != nil {
		return nil, queryErr
	}

	var logs []*entities.TopKillerLog

	for rows.Next() {
		log := &entities.TopKillerLog{}
		scanErr := rows.StructScan(log)
		if scanErr != nil {
			return nil, scanErr
		}
		logs = append(logs, log)
	}

	return logs, nil
}

func (dao *TeamKillLogDAO) GetTopVictims() ([]*entities.TopVictimLog, error) {
	query := `
		SELECT victim as name, count(*) as deaths_count
		FROM team_kill_log
		GROUP BY victim
		ORDER BY deaths_count desc;
	`
	rows, queryErr := dao.db.Queryx(query)
	if queryErr != nil {
		return nil, queryErr
	}

	var logs []*entities.TopVictimLog

	for rows.Next() {
		log := &entities.TopVictimLog{}
		scanErr := rows.StructScan(log)
		if scanErr != nil {
			return nil, scanErr
		}
		logs = append(logs, log)
	}

	return logs, nil
}
