package services

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestNewTeamKillFromText_PassedCorrectMessage_ParsedToRequest(t *testing.T) {
	text := "Лёша #teamkill Игорь"

	parsedMessage := NewTeamKillFromText(text)

	assert.Equal(t, "лёша", parsedMessage.Killer)
	assert.Equal(t, "игорь", parsedMessage.Victim)
}

func TestNewTeamKillFromText_PassedMessageWithoutKiller_NotParsed(t *testing.T) {
	text := "#teamkill Игорь"

	parsedMessage := NewTeamKillFromText(text)

	assert.Nil(t, parsedMessage)
}

func TestNewTeamKillFromText_PassedMessageWithoutVictim_NotParsed(t *testing.T) {
	text := "Лёша #teamkill"

	parsedMessage := NewTeamKillFromText(text)

	assert.Nil(t, parsedMessage)
}

func TestNewTeamKillFromText_PassedMessageWithoutKillerWithSpace_NotParsed(t *testing.T) {
	text := "#teamkill Игорь"

	parsedMessage := NewTeamKillFromText(text)

	assert.Nil(t, parsedMessage)
}

func TestNewTeamKillFromText_PassedMessageWithoutVictimWithSpace_NotParsed(t *testing.T) {
	text := "Лёша #teamkill"

	parsedMessage := NewTeamKillFromText(text)

	assert.Nil(t, parsedMessage)
}

func TestNewTeamKillFromText_PassedMessageWithoutVictimWithoutKiller_NotParsed(t *testing.T) {
	text := "#teamkill"

	parsedMessage := NewTeamKillFromText(text)

	assert.Nil(t, parsedMessage)
}
