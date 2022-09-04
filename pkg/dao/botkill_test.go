package dao

import (
	"github.com/jmoiron/sqlx"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"log"
	"pthd-bot/pkg/entities"
	"testing"
	"time"
)

type BotKillLogDAOTestCase struct {
	suite.Suite

	db  *sqlx.DB
	dao *BotKillLogDAO
}

func (s *BotKillLogDAOTestCase) SetupSuite() {
	db, openErr := sqlx.Open("sqlite3", ":memory:")
	if openErr != nil {
		log.Fatalf("Cannot open test sqlite: %s", openErr)
	}

	s.db = db

	migrationErr := MigrateUp(s.db)
	if migrationErr != nil {
		log.Fatalf("Cannot migrate db: %s", migrationErr)
	}

	s.dao = NewBotKillLogDAO(s.db)
}

func (s *BotKillLogDAOTestCase) TearDownSuite() {
	s.db.Close()
}

func (s *BotKillLogDAOTestCase) BeforeTest(suiteName, testName string) {
	s.db.Exec(`DELETE FROM bot_kill_log`)
}

func (s *BotKillLogDAOTestCase) Test__InsertBotKillLog__VictimsCalculatedCorrectly() {
	botKill := &entities.BotKill{
		Victim:     "victim",
		Source:     "source",
		HappenedAt: time.Now(),
	}

	saveErr := s.dao.Save(botKill)
	assert.Nil(s.T(), saveErr)

	saveErr = s.dao.Save(botKill)
	assert.Nil(s.T(), saveErr)

	topVictim, getErr := s.dao.GetTopVictims("source")

	assert.Nil(s.T(), getErr)

	assert.Equal(s.T(), 1, len(topVictim))
	assert.Equal(s.T(), "victim", topVictim[0].Name)
	assert.Equal(s.T(), 2, topVictim[0].DeathsCount)
}

func (s *BotKillLogDAOTestCase) Test__InsertBotKillLogWithDifferentSources__VictimsCalculatedCorrectly() {
	botKillOne := &entities.BotKill{
		Victim:     "victim",
		Source:     "source",
		HappenedAt: time.Now(),
	}

	botKillTwo := &entities.BotKill{
		Victim:     "victim",
		Source:     "source1",
		HappenedAt: time.Now(),
	}

	saveErr := s.dao.Save(botKillOne)
	assert.Nil(s.T(), saveErr)

	saveErr = s.dao.Save(botKillTwo)
	assert.Nil(s.T(), saveErr)

	topVictim, getErr := s.dao.GetTopVictims("source")

	assert.Nil(s.T(), getErr)

	assert.Equal(s.T(), 1, len(topVictim))
	assert.Equal(s.T(), "victim", topVictim[0].Name)
	assert.Equal(s.T(), 1, topVictim[0].DeathsCount)
}

func TestBotKillLogDAO(t *testing.T) {
	testSuite := BotKillLogDAOTestCase{}
	suite.Run(t, &testSuite)
}
