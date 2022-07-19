package services

import (
	"strings"
)

type TeamKillRequest struct {
	Killer string
	Victim string
}

func NewTeamKillFromText(text string) *TeamKillRequest {

	parts := strings.Split(text, "#teamkill")

	convert := func(name string) string {
		name = strings.ToLower(name)
		name = strings.Trim(name, " ")
		return name
	}

	killer := convert(parts[0])
	victim := convert(parts[1])

	if killer == "" || victim == "" {
		return nil
	}

	return &TeamKillRequest{
		Killer: convert(parts[0]),
		Victim: convert(parts[1]),
	}
}
