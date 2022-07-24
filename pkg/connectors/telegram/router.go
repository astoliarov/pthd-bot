package telegram

import (
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"log"
	"strings"
	"teamkillbot/pkg/services"
)

func getMessageType(text string) int {
	if strings.Contains(text, teamKillMessageKey) {
		return messageTypeKillLog
	}
	if strings.HasPrefix(text, echoKey) {
		return messageTypeEcho
	}

	return unprocessableMessage
}

func replyToMessage(bot *tgbotapi.BotAPI, message *tgbotapi.Message, responseText string) error {
	msg := tgbotapi.NewMessage(message.Chat.ID, responseText)
	msg.ReplyToMessageID = message.MessageID
	_, sendErr := bot.Send(msg)
	if sendErr != nil {
		return sendErr
	}
	return nil
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
	return replyToMessage(r.bot, message, message.Text)
}

func (r MessageRouter) processMessageTeamKill(message *tgbotapi.Message) error {
	log.Printf("processing message team kill")
	request := services.NewTeamKillFromText(message.Text)
	if request == nil {
		sendErr := replyToMessage(r.bot, message, cannotParseTeamKillMessage)
		if sendErr != nil {
			return sendErr
		}
	}

	response, processErr := r.teamKillService.ProcessTeamKill(request)
	if processErr != nil {
		return processErr
	}

	if response != "" {
		sendErr := replyToMessage(r.bot, message, response)
		if sendErr != nil {
			return sendErr
		}
	}

	return nil
}
