package services

import (
	"pthd-bot/pkg/entities"
	"pthd-bot/pkg/interfaces"
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

func (s *TeamKillService) ProcessTeamKill(request *TeamKillRequest, source string) (string, error) {
	teamKill := &entities.TeamKill{
		Killer:     normalizeName(request.Killer),
		Victim:     normalizeName(request.Victim),
		Source:     source,
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

func (s *TeamKillService) ProcessGetTopKillers(source string) ([]*entities.TopKillerLog, error) {
	return s.dao.GetTopKillers(source)
}

func (s *TeamKillService) ProcessGetTopVictims(source string) ([]*entities.TopVictimLog, error) {
	return s.dao.GetTopVictims(source)
}
