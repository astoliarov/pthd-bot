package services

import (
	"teamkillbot/pkg/entities"
	"teamkillbot/pkg/interfaces"
	"time"
)

type BotKillService struct {
	dao              interfaces.IBotKillLogDAO
	responseSelector *ResponseSelectorService
}

func NewBotKillService(dao interfaces.IBotKillLogDAO, responseSelector *ResponseSelectorService) *BotKillService {
	return &BotKillService{
		dao:              dao,
		responseSelector: responseSelector,
	}
}

func (s *BotKillService) ProcessBotKill(request *BotKillRequest, source string) (string, error) {
	botKill := &entities.BotKill{
		Victim:     normalizeName(request.Victim),
		Source:     source,
		HappenedAt: time.Now(),
	}
	saveErr := s.dao.Save(botKill)
	if saveErr != nil {
		return "", saveErr
	}

	response, getResponseErr := s.responseSelector.GetResponse()
	if getResponseErr != nil {
		return "", getResponseErr
	}

	return response, nil
}

func (s *BotKillService) ProcessGetTopVictims(source string) ([]*entities.TopVictimLog, error) {
	return s.dao.GetTopVictims(source)
}
