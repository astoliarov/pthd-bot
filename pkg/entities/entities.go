package entities

import "time"

type TeamKill struct {
	Killer     string
	Victim     string
	HappenedAt time.Time
	Source     string
}

type BotKill struct {
	Victim     string
	HappenedAt time.Time
	Source     string
}

type TopKillerLog struct {
	Name      string `db:"name"`
	KillCount int    `db:"kill_count"`
}

type TopVictimLog struct {
	Name        string `db:"name"`
	DeathsCount int    `db:"deaths_count"`
}
