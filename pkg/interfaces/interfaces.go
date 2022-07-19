package interfaces

import (
	"teamkillbot/pkg/entities"
)

//go:generate mockgen -destination=../tests/mocks/iteamkilldao_mock.go -package=mocks teamkillbot/pkg ITeamKillLogDAO
type ITeamKillLogDAO interface {
	Save(kill *entities.TeamKill) error
}
