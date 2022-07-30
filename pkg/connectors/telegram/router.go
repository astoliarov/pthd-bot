package telegram

import (
	"fmt"
	"github.com/getsentry/sentry-go"
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
	if strings.HasPrefix(text, showKillers) {
		return messageTypeShowKillers
	}
	if strings.HasPrefix(text, showVictims) {
		return messageTypeShowVictims
	}
	if strings.HasPrefix(text, repeat) {
		return messageTypeRepeat
	}

	return unprocessableMessage
}

func sendMessage(bot *tgbotapi.BotAPI, chatId int64, responseText string) error {
	msg := tgbotapi.NewMessage(chatId, responseText)
	_, sendErr := bot.Send(msg)
	if sendErr != nil {
		return sendErr
	}
	return nil
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
	bot               *tgbotapi.BotAPI
	teamKillService   *services.TeamKillService
	messageProcessors map[int]func(message *tgbotapi.Message) error
	chatID            int64
}

func NewMessageRouter(bot *tgbotapi.BotAPI, service *services.TeamKillService, chatID int64) *MessageRouter {
	router := &MessageRouter{
		bot:               bot,
		teamKillService:   service,
		messageProcessors: map[int]func(message *tgbotapi.Message) error{},
		chatID:            chatID,
	}

	router.register(messageTypeEcho, router.processMessageEcho)
	router.register(messageTypeKillLog, router.processMessageTeamKill)
	router.register(messageTypeShowKillers, router.processMessageShowKillers)
	router.register(messageTypeShowVictims, router.processMessageShowVictims)
	router.register(messageTypeRepeat, router.processMessageRepeat)

	return router
}

func (r *MessageRouter) register(code int, processor func(message *tgbotapi.Message) error) {
	r.messageProcessors[code] = processor
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
				sentry.CaptureException(processErr)
				log.Printf("Error while processing message: %s", processErr)
			}
		}
	}
}

func (r MessageRouter) processMessage(message *tgbotapi.Message) error {
	messageType := getMessageType(message.Text)

	if message.Chat.ID != r.chatID {
		return nil
	}

	processor, exists := r.messageProcessors[messageType]
	if !exists {
		return nil
	}

	err := processor(message)
	if err != nil {
		return err
	}

	return nil
}

func (r MessageRouter) processMessageEcho(message *tgbotapi.Message) error {
	return replyToMessage(r.bot, message, message.Text)
}

func (r MessageRouter) processMessageShowKillers(message *tgbotapi.Message) error {

	topKillersLog, err := r.teamKillService.ProcessGetTopKillers()
	if err != nil {
		replyErr := replyToMessage(r.bot, message, cannotProcessMessage)
		log.Printf("Error while trying to send message about error, while processing command: %s", replyErr)
		return err
	}

	var rows []string
	for _, killLog := range topKillersLog {
		row := fmt.Sprintf("Убивца: %s; Намолотил: %d", killLog.Name, killLog.KillCount)
		rows = append(rows, row)
	}
	msg := strings.Join(rows, "\n")

	sendErr := sendMessage(r.bot, message.Chat.ID, msg)
	if sendErr != nil {
		log.Printf("Error while trying to send message: %s", sendErr)
	}

	return nil
}

func (r MessageRouter) processMessageShowVictims(message *tgbotapi.Message) error {

	topVictimsLog, err := r.teamKillService.ProcessGetTopVictims()
	if err != nil {
		replyErr := replyToMessage(r.bot, message, cannotProcessMessage)
		log.Printf("Error while trying to send message about error, while processing command: %s", replyErr)
		return err
	}

	var rows []string
	for _, victimLog := range topVictimsLog {
		row := fmt.Sprintf("Неудачник: %s; Наумирал: %d", victimLog.Name, victimLog.DeathsCount)
		rows = append(rows, row)
	}
	msg := strings.Join(rows, "\n")

	sendErr := sendMessage(r.bot, message.Chat.ID, msg)
	if sendErr != nil {
		log.Printf("Error while trying to send message: %s", sendErr)
	}

	return nil
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

func (r MessageRouter) processMessageRepeat(message *tgbotapi.Message) error {
	log.Printf("processing message repeat")

	if message.ReplyToMessage == nil {
		sendErr := replyToMessage(r.bot, message, "Нечего повторять, ну")
		if sendErr != nil {
			return sendErr
		}
		return nil
	}

	request := services.NewTeamKillFromText(message.ReplyToMessage.Text)
	if request == nil {
		sendErr := replyToMessage(r.bot, message.ReplyToMessage, cannotParseTeamKillMessage)
		if sendErr != nil {
			return sendErr
		}
	}

	response, processErr := r.teamKillService.ProcessTeamKill(request)
	if processErr != nil {
		return processErr
	}

	if response != "" {
		sendErr := replyToMessage(r.bot, message.ReplyToMessage, response)
		if sendErr != nil {
			return sendErr
		}
	}

	return nil
}
