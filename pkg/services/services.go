package services

import (
	"strings"
	"teamkillbot/pkg/entities"
	"teamkillbot/pkg/interfaces"
	"time"
)

type TeamKillService struct {
	dao              interfaces.ITeamKillLogDAO
	responseSelector *ResponseSelectorService
}

func NewTeamKillService(dao interfaces.ITeamKillLogDAO, responseSelector *ResponseSelectorService) *TeamKillService {
	return &TeamKillService{
		dao:              dao,
		responseSelector: responseSelector,
	}
}

func (s *TeamKillService) ProcessTeamKill(request *TeamKillRequest) (string, error) {

	teamKill := &entities.TeamKill{
		Killer:     s.normalizeName(request.Killer),
		Victim:     s.normalizeName(request.Victim),
		HappenedAt: time.Now(),
	}
	saveErr := s.dao.Save(teamKill)
	if saveErr != nil {
		return "", saveErr
	}

	response, getResponseErr := s.responseSelector.GetResponse()
	if getResponseErr != nil {
		return "", getResponseErr
	}

	return response, nil
}

func (s *TeamKillService) normalizeName(name string) string {
	name = strings.ToLower(name)
	name = strings.Replace(name, "ё", "e", -1)
	return name
}

func (s *TeamKillService) ProcessGetTopKillers() ([]*entities.TopKillerLog, error) {
	return s.dao.GetTopKillers()
}

func (s *TeamKillService) ProcessGetTopVictims() ([]*entities.TopVictimLog, error) {
	return s.dao.GetTopVictims()
}
