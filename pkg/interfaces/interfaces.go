package interfaces

import (
	"teamkillbot/pkg/entities"

	_ "github.com/golang/mock/mockgen/model"
)

//go:generate mockgen -destination=../../tests/mocks/iteamkilldao_mock.go -package=mocks teamkillbot/pkg/interfaces ITeamKillLogDAO
type ITeamKillLogDAO interface {
	Save(kill *entities.TeamKill) error
	GetTopKillers(source string) ([]*entities.TopKillerLog, error)
	GetTopVictims(source string) ([]*entities.TopVictimLog, error)
}
