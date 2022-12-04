package services

import (
	"pthd-bot/pkg/entities"
	"pthd-bot/pkg/interfaces"
	"time"
	"unicode"
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

func (s *TeamKillService) validateName(name string) error {
	isCyrSymbolsPresent := false
	isLatSymbolsPresent := false

	for _, symbol := range name {
		if unicode.Is(unicode.Cyrillic, symbol) {
			isCyrSymbolsPresent = true
		}
		if unicode.Is(unicode.Latin, symbol) {
			isLatSymbolsPresent = true
		}
	}

	if isCyrSymbolsPresent && isLatSymbolsPresent {
		return NewErrMixedCharactersInName(name)
	}

	return nil
}

func (s *TeamKillService) AddTeamKill(request *TeamKillRequest, source string) (string, error) {
	if err := s.validateName(request.Killer); err != nil {
		return "", err
	}

	if err := s.validateName(request.Victim); err != nil {
		return "", err
	}

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

func (s *TeamKillService) GetTopKillers(source string) ([]*entities.TopKillerLog, error) {
	return s.dao.GetTopKillers(source)
}

func (s *TeamKillService) GetTopVictims(source string) ([]*entities.TopVictimLog, error) {
	return s.dao.GetTopVictims(source)
}
