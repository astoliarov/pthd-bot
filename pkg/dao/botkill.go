package dao

import (
	"github.com/jmoiron/sqlx"
	"pthd-bot/pkg/entities"
)

type BotKillLogDAO struct {
	db *sqlx.DB
}

func NewBotKillLogDAO(db *sqlx.DB) *BotKillLogDAO {
	return &BotKillLogDAO{
		db: db,
	}
}

func (dao *BotKillLogDAO) Save(kill *entities.BotKill) error {
	query := `
		INSERT INTO bot_kill_log(victim, happened_at, source) values(?,?,?)
	`
	stmt, prepareErr := dao.db.Prepare(query)
	if prepareErr != nil {
		return prepareErr
	}

	_, execErr := stmt.Exec(kill.Victim, kill.HappenedAt, kill.Source)
	if execErr != nil {
		return execErr
	}

	return nil
}

func (dao *BotKillLogDAO) GetTopVictims(source string) ([]*entities.TopVictimLog, error) {
	var logs []entities.TopVictimLog
	query := `
		SELECT victim as name, count(*) as deaths_count
		FROM bot_kill_log
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
