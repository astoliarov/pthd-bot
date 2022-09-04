package dao

import (
	"github.com/jmoiron/sqlx"
	_ "github.com/mattn/go-sqlite3"
	"pthd-bot/pkg/entities"
)

type TeamKillLogDAO struct {
	db *sqlx.DB
}

func NewTeamKillLogDAO(db *sqlx.DB) *TeamKillLogDAO {
	return &TeamKillLogDAO{
		db: db,
	}
}

func (dao *TeamKillLogDAO) Save(kill *entities.TeamKill) error {
	query := `
		INSERT INTO team_kill_log(killer, victim, happened_at, source) values(?,?,?,?)
	`
	stmt, prepareErr := dao.db.Prepare(query)
	if prepareErr != nil {
		return prepareErr
	}

	_, execErr := stmt.Exec(kill.Killer, kill.Victim, kill.HappenedAt, kill.Source)
	if execErr != nil {
		return execErr
	}

	return nil
}

func (dao *TeamKillLogDAO) GetTopKillers(source string) ([]*entities.TopKillerLog, error) {
	var logs []entities.TopKillerLog
	query := `
		SELECT killer as name, count(*) as kill_count
		FROM team_kill_log
		WHERE source = ?
		GROUP BY killer
		ORDER BY kill_count desc;
	`

	queryErr := dao.db.Select(&logs, query, source)
	if queryErr != nil {
		return nil, queryErr
	}

	var resultLogs []*entities.TopKillerLog
	for idx := range logs {
		resultLogs = append(resultLogs, &logs[idx])
	}

	return resultLogs, nil
}

func (dao *TeamKillLogDAO) GetTopVictims(source string) ([]*entities.TopVictimLog, error) {
	var logs []entities.TopVictimLog
	query := `
		SELECT victim as name, count(*) as deaths_count
		FROM team_kill_log
		WHERE source = ?
		GROUP BY victim
		ORDER BY deaths_count desc;
	`
	queryErr := dao.db.Select(&logs, query, source)
	if queryErr != nil {
		return nil, queryErr
	}

	var resultLogs []*entities.TopVictimLog
	for idx := range logs {
		resultLogs = append(resultLogs, &logs[idx])
	}

	return resultLogs, nil
}
