package telegram

import (
	"context"
	"github.com/getsentry/sentry-go"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"log"
	"teamkillbot/pkg/services"
)

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
	commands          []ICommand
	chatID            int64
}

func NewMessageRouter(bot *tgbotapi.BotAPI, service *services.TeamKillService, chatID int64) *MessageRouter {

	echoCommand := &CommandEcho{}
	teamKillCommand := &TeamKillCommand{teamKillService: service}
	showKillersCommand := &ShowKillersCommand{teamKillService: service}
	showVictimsCommand := &ShowVictimsCommand{teamKillService: service}
	repeatCommand := &RepeatCommand{teamKillService: service}
	helpCommand := &HelpCommand{commands: []ICommand{}}
	commands := []ICommand{
		echoCommand,
		teamKillCommand,
		showKillersCommand,
		showVictimsCommand,
		repeatCommand,
		helpCommand,
	}
	helpCommand.SetCommands(commands)

	router := &MessageRouter{
		bot:             bot,
		teamKillService: service,
		commands:        commands,
		chatID:          chatID,
	}

	return router
}

func (r *MessageRouter) registerCommand(command ICommand) {
	r.commands = append(r.commands, command)
}

func (r MessageRouter) ListenToUpdates(ctx context.Context) {
	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates := r.bot.GetUpdatesChan(u)

	for {
		select {
		case update := <-updates:
			if update.Message != nil { // If we got a message
				log.Printf("[%s] %s", update.Message.From.UserName, update.Message.Text)
				processErr := r.processMessage(update.Message)
				if processErr != nil {
					sentry.CaptureException(processErr)
					log.Printf("Error while processing message: %s", processErr)
				}
			}
		case <-ctx.Done():
			log.Println("Stopping router")
			return
		}
	}
}

func (r MessageRouter) processMessage(message *tgbotapi.Message) error {
	if message.Chat.ID != r.chatID {
		return nil
	}

	for _, command := range r.commands {
		if command.IsCommand(message.Text) {
			processErr := command.Process(r.bot, message)
			if processErr != nil {
				return processErr
			}
		}
	}

	return nil
}
