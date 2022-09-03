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

	if len(parts) < 2 {
		return nil
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

type BotKillRequest struct {
	Victim string
}

func NewBotKillRequest(text string) *BotKillRequest {

	parts := strings.Split(text, "#botkill")

	convert := func(name string) string {
		name = strings.ToLower(name)
		name = strings.Trim(name, " ")
		return name
	}

	if len(parts) < 1 {
		return nil
	}

	victim := convert(parts[0])

	if victim == "" {
		return nil
	}

	return &BotKillRequest{
		Victim: convert(parts[0]),
	}
}
