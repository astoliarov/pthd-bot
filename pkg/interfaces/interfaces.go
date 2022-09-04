package interfaces

import (
	"pthd-bot/pkg/entities"

	_ "github.com/golang/mock/mockgen/model"
)

//go:generate mockgen -destination=../../tests/mocks/iteamkilldao_mock.go -package=mocks pthd-bot/pkg/interfaces ITeamKillLogDAO
type ITeamKillLogDAO interface {
	Save(kill *entities.TeamKill) error
	GetTopKillers(source string) ([]*entities.TopKillerLog, error)
	GetTopVictims(source string) ([]*entities.TopVictimLog, error)
}

//go:generate mockgen -destination=../../tests/mocks/ibotkilldao_mock.go -package=mocks pthd-bot/pkg/interfaces IBotKillLogDAO
type IBotKillLogDAO interface {
	Save(kill *entities.BotKill) error
	GetTopVictims(source string) ([]*entities.TopVictimLog, error)
}
