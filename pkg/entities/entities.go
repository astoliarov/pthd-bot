package entities

import "time"

type TeamKill struct {
	Killer     string
	Victim     string
	HappenedAt time.Time
}
