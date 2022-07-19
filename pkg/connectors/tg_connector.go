package connectors

import (
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"log"
	"strings"
	"teamkillbot/pkg/services"
)

const unprocessableMessage = 0
const messageTypeKillLog = 1
const messageTypeEcho = 2

const teamKillMessageKey = "#teamkill"
const echoKey = "PTHD:echo"

const CannotParseTeamKillMessage = "Чё ты понаписал? Не понятно ничего"

func getMessageType(text string) int {
	if strings.Contains(text, teamKillMessageKey) {
		return messageTypeKillLog
	}
	if strings.HasPrefix(text, echoKey) {
		return messageTypeEcho
	}

	return unprocessableMessage
}

type MessageRouter struct {
	bot             *tgbotapi.BotAPI
	teamKillService *services.TeamKillService
}

func NewMessageRouter(bot *tgbotapi.BotAPI, service *services.TeamKillService) *MessageRouter {
	return &MessageRouter{
		bot:             bot,
		teamKillService: service,
	}
}

func (r MessageRouter) ListenToUpdates() {
	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates := r.bot.GetUpdatesChan(u)

	for update := range updates {
		if update.Message != nil { // If we got a message
			log.Printf("[%s] %s", update.Message.From.UserName, update.Message.Text)
			processErr := r.processMessage(update.Message)
			if processErr != nil {
				log.Printf("Error while processing message: %s", processErr)
			}
		}
	}
}

func (r MessageRouter) processMessage(message *tgbotapi.Message) error {
	messageType := getMessageType(message.Text)
	fmt.Println(messageType)

	if messageType == unprocessableMessage {
		return nil
	}

	if messageType == messageTypeEcho {
		return r.processMessageEcho(message)
	}

	if messageType == messageTypeKillLog {
		return r.processMessageTeamKill(message)
	}

	return nil
}

func (r MessageRouter) processMessageEcho(message *tgbotapi.Message) error {
	msg := tgbotapi.NewMessage(message.Chat.ID, message.Text)
	msg.ReplyToMessageID = message.MessageID

	_, sendErr := r.bot.Send(msg)

	return sendErr
}

func (r MessageRouter) processMessageTeamKill(message *tgbotapi.Message) error {
	log.Printf("processing message team kill")
	request := services.NewTeamKillFromText(message.Text)
	if request == nil {
		msg := tgbotapi.NewMessage(message.Chat.ID, CannotParseTeamKillMessage)
		msg.ReplyToMessageID = message.MessageID
		_, sendErr := r.bot.Send(msg)
		if sendErr != nil {
			return sendErr
		}
	}

	response, processErr := r.teamKillService.ProcessTeamKill(request)
	if processErr != nil {
		return processErr
	}

	if response != "" {
		msg := tgbotapi.NewMessage(message.Chat.ID, response)
		msg.ReplyToMessageID = message.MessageID
		_, sendErr := r.bot.Send(msg)
		if sendErr != nil {
			return sendErr
		}
	}

	return nil
}
