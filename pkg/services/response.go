package services

import "math/rand"

var defaulResponses = []string{
	"Красиво",
	"А что случилось?",
	"\"О, вкусный ЧВКшник\" (c)",
}

type ResponseSelectorService struct{}

func (s *ResponseSelectorService) GetResponse() (string, error) {
	response := defaulResponses[rand.Intn(len(defaulResponses))]
	return response, nil
}
