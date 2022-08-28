package dao

import (
	"github.com/jmoiron/sqlx"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"log"
	"teamkillbot/pkg/entities"
	"testing"
	"time"
)

type TeamKillLogDAOTestCase struct {
	suite.Suite

	db  *sqlx.DB
	dao *TeamKillLogDAO
}

func (s *TeamKillLogDAOTestCase) SetupSuite() {
	log.Println("SetupSuite()")

	db, openErr := sqlx.Open("sqlite3", ":memory:")
	if openErr != nil {
		log.Fatalf("Cannot open test sqlite: %s", openErr)
	}

	s.db = db

	migrationErr := MigrateUp(s.db)
	if migrationErr != nil {
		log.Fatalf("Cannot migrate db: %s", migrationErr)
	}

	s.dao = NewTeamKillLogDAO(s.db)
}

func (s *TeamKillLogDAOTestCase) TearDownSuite() {
	log.Println("TearDownSuite()")

	s.db.Close()
}

func (s *TeamKillLogDAOTestCase) BeforeTest(suiteName, testName string) {
	s.db.Exec(`DELETE FROM team_kill_log`)
}

func (s *TeamKillLogDAOTestCase) Test__InsertTeamKillLog__KillersCalculatedCorrectly() {
	log.Println("TestExample")

	teamKill := &entities.TeamKill{
		Killer:     "killer",
		Victim:     "victim",
		Source:     "source",
		HappenedAt: time.Now(),
	}

	saveErr := s.dao.Save(teamKill)
	assert.Nil(s.T(), saveErr)

	saveErr = s.dao.Save(teamKill)
	assert.Nil(s.T(), saveErr)

	topKiller, getErr := s.dao.GetTopKillers("source")

	assert.Nil(s.T(), getErr)

	assert.Equal(s.T(), 1, len(topKiller))
	assert.Equal(s.T(), "killer", topKiller[0].Name)
	assert.Equal(s.T(), 2, topKiller[0].KillCount)
}

func (s *TeamKillLogDAOTestCase) Test__InsertTeamKillLog__VictimsCalculatedCorrectly() {
	log.Println("TestExample")

	teamKill := &entities.TeamKill{
		Killer:     "killer",
		Victim:     "victim",
		Source:     "source",
		HappenedAt: time.Now(),
	}

	saveErr := s.dao.Save(teamKill)
	assert.Nil(s.T(), saveErr)

	saveErr = s.dao.Save(teamKill)
	assert.Nil(s.T(), saveErr)

	topVictims, getErr := s.dao.GetTopVictims("source")

	assert.Nil(s.T(), getErr)

	assert.Equal(s.T(), 1, len(topVictims))
	assert.Equal(s.T(), "victim", topVictims[0].Name)
	assert.Equal(s.T(), 2, topVictims[0].DeathsCount)
}

func (s *TeamKillLogDAOTestCase) Test__InsertTeamKillLogWithDifferentKillers__KillersCountedCorrectly() {
	log.Println("TestExample")

	teamKillOne := &entities.TeamKill{
		Killer:     "killer",
		Victim:     "victim",
		Source:     "source",
		HappenedAt: time.Now(),
	}

	teamKillTwo := &entities.TeamKill{
		Killer:     "killer2",
		Victim:     "victim",
		Source:     "source",
		HappenedAt: time.Now(),
	}

	saveErr := s.dao.Save(teamKillOne)
	assert.Nil(s.T(), saveErr)

	saveErr = s.dao.Save(teamKillTwo)
	assert.Nil(s.T(), saveErr)

	topKillers, getErr := s.dao.GetTopKillers("source")

	assert.Nil(s.T(), getErr)

	var killerOneLog *entities.TopKillerLog
	for _, killerLog := range topKillers {
		if killerLog.Name == "killer" {
			killerOneLog = killerLog
		}
	}

	var killerTwoLog *entities.TopKillerLog
	for _, killerLog := range topKillers {
		if killerLog.Name == "killer2" {
			killerTwoLog = killerLog
		}
	}

	assert.Equal(s.T(), 2, len(topKillers))

	assert.Equal(s.T(), 1, killerOneLog.KillCount)
	assert.Equal(s.T(), 1, killerTwoLog.KillCount)
}

func (s *TeamKillLogDAOTestCase) Test__InsertTeamKillLogWithDifferentKillers__VictimsCountedCorrectly() {
	log.Println("TestExample")

	teamKillOne := &entities.TeamKill{
		Killer:     "killer",
		Victim:     "victim",
		Source:     "source",
		HappenedAt: time.Now(),
	}

	teamKillTwo := &entities.TeamKill{
		Killer:     "killer",
		Victim:     "victim2",
		Source:     "source",
		HappenedAt: time.Now(),
	}

	saveErr := s.dao.Save(teamKillOne)
	assert.Nil(s.T(), saveErr)

	saveErr = s.dao.Save(teamKillTwo)
	assert.Nil(s.T(), saveErr)

	topVictims, getErr := s.dao.GetTopVictims("source")

	assert.Nil(s.T(), getErr)

	var victimOneLog *entities.TopVictimLog
	for _, victimLog := range topVictims {
		if victimLog.Name == "victim" {
			victimOneLog = victimLog
		}
	}

	var victimTwoLog *entities.TopVictimLog
	for _, victimLog := range topVictims {
		if victimLog.Name == "victim2" {
			victimTwoLog = victimLog
		}
	}

	assert.Equal(s.T(), 2, len(topVictims))

	assert.Equal(s.T(), 1, victimOneLog.DeathsCount)
	assert.Equal(s.T(), 1, victimTwoLog.DeathsCount)
}

func (s *TeamKillLogDAOTestCase) Test__InsertTeamKillLogWithDifferentSources__KillersCountedCorrectly() {
	log.Println("TestExample")

	teamKillOne := &entities.TeamKill{
		Killer:     "killer",
		Victim:     "victim",
		Source:     "source",
		HappenedAt: time.Now(),
	}

	teamKillTwo := &entities.TeamKill{
		Killer:     "killer",
		Victim:     "victim",
		Source:     "source2",
		HappenedAt: time.Now(),
	}

	saveErr := s.dao.Save(teamKillOne)
	assert.Nil(s.T(), saveErr)

	saveErr = s.dao.Save(teamKillTwo)
	assert.Nil(s.T(), saveErr)

	topKillers, getErr := s.dao.GetTopKillers("source")

	assert.Nil(s.T(), getErr)

	assert.Equal(s.T(), 1, len(topKillers))

	assert.Equal(s.T(), 1, topKillers[0].KillCount)
}

func (s *TeamKillLogDAOTestCase) Test__InsertTeamKillLogWithDifferentSources__VictimsCountedCorrectly() {
	log.Println("TestExample")

	teamKillOne := &entities.TeamKill{
		Killer:     "killer",
		Victim:     "victim",
		Source:     "source",
		HappenedAt: time.Now(),
	}

	teamKillTwo := &entities.TeamKill{
		Killer:     "killer",
		Victim:     "victim",
		Source:     "source2",
		HappenedAt: time.Now(),
	}

	saveErr := s.dao.Save(teamKillOne)
	assert.Nil(s.T(), saveErr)

	saveErr = s.dao.Save(teamKillTwo)
	assert.Nil(s.T(), saveErr)

	topVictims, getErr := s.dao.GetTopVictims("source")

	assert.Nil(s.T(), getErr)

	assert.Equal(s.T(), 1, len(topVictims))

	assert.Equal(s.T(), 1, topVictims[0].DeathsCount)
}

func TestTeamKillLogDAO(t *testing.T) {
	testSuite := TeamKillLogDAOTestCase{}
	suite.Run(t, &testSuite)
}
