package telegram

import (
	"context"
	"fmt"
	"github.com/getsentry/sentry-go"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/rs/zerolog/log"
	"pthd-bot/pkg/services"
)

func sendMessage(bot *tgbotapi.BotAPI, chatId int64, responseText string) error {
	msg := tgbotapi.NewMessage(chatId, responseText)
	_, sendErr := bot.Send(msg)
	if sendErr != nil {
		log.Error().Err(sendErr).Msg("Error while trying to send message")
		return sendErr
	}
	return nil
}

func sourceFromMessage(message *tgbotapi.Message) string {
	return fmt.Sprintf("tg_%d", message.Chat.ID)
}

func replyToMessage(bot *tgbotapi.BotAPI, message *tgbotapi.Message, responseText string) error {
	msg := tgbotapi.NewMessage(message.Chat.ID, responseText)
	msg.ReplyToMessageID = message.MessageID
	_, sendErr := bot.Send(msg)
	if sendErr != nil {
		log.Error().Err(sendErr).Msg("Error while trying to reply to message")
		return sendErr
	}
	return nil
}

type MessageRouter struct {
	bot               *tgbotapi.BotAPI
	teamKillService   *services.TeamKillService
	botKillService    *services.TeamKillService
	messageProcessors map[int]func(message *tgbotapi.Message) error
	commands          []ICommand
	chatID            int64
}

func NewMessageRouter(
	bot *tgbotapi.BotAPI,
	teamKillService *services.TeamKillService,
	botKillService *services.BotKillService,
	chatID int64,
) *MessageRouter {

	echoCommand := &CommandEcho{}
	teamKillCommand := &TeamKillCommand{teamKillService: teamKillService}
	botKillCommand := &BotKillCommand{service: botKillService}
	showKillersCommand := &ShowKillersCommand{teamKillService: teamKillService}
	showVictimsCommand := &ShowVictimsCommand{teamKillService: teamKillService}
	showBotVictimsCommand := &ShowBotVictimsCommand{service: botKillService}
	repeatCommand := &RepeatCommand{teamKillService: teamKillService}
	helpCommand := &HelpCommand{commands: []ICommand{}}
	commands := []ICommand{
		echoCommand,
		botKillCommand,
		teamKillCommand,
		showKillersCommand,
		showVictimsCommand,
		showBotVictimsCommand,
		repeatCommand,
		helpCommand,
	}
	helpCommand.SetCommands(commands)

	router := &MessageRouter{
		bot:             bot,
		teamKillService: teamKillService,
		commands:        commands,
		chatID:          chatID,
	}

	return router
}

func (r *MessageRouter) registerCommand(command ICommand) {
	r.commands = append(r.commands, command)
}

func (r *MessageRouter) ListenToUpdates(ctx context.Context) {
	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates := r.bot.GetUpdatesChan(u)

	for {
		select {
		case update := <-updates:
			if update.Message != nil { // If we got a message
				log.Debug().
					Str("from", update.Message.From.UserName).
					Str("text", update.Message.Text).
					Msg("received message")

				processErr := r.processMessage(update.Message)
				if processErr != nil {
					sentry.CaptureException(processErr)
					log.Error().
						Err(processErr).
						Msg("Error while processing message")
				}
			}
		case <-ctx.Done():
			log.Info().Msg("Stopping router")
			return
		}
	}
}

func (r *MessageRouter) processMessage(message *tgbotapi.Message) error {
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
