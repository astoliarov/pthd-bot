package telegram

import (
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/rs/zerolog/log"
	"pthd-bot/pkg/services"
	"strings"
)

type ICommand interface {
	Process(bot *tgbotapi.BotAPI, message *tgbotapi.Message) error
	IsCommand(string) bool
	GetHelp() string
}

type CommandEcho struct{}

func (c *CommandEcho) Process(bot *tgbotapi.BotAPI, message *tgbotapi.Message) error {
	log.Info().Msg("processing message echo")
	return replyToMessage(bot, message, message.Text)
}

func (c *CommandEcho) IsCommand(text string) bool {
	return strings.HasPrefix(text, echoKey)
}

func (c *CommandEcho) GetHelp() string {
	return fmt.Sprintf("%s - повторить сообщение", echoKey)
}

type TeamKillCommand struct {
	teamKillService *services.TeamKillService
}

func (c *TeamKillCommand) Process(bot *tgbotapi.BotAPI, message *tgbotapi.Message) error {
	log.Info().Msg("processing message team kill")

	request := services.NewTeamKillFromText(message.Text)
	if request == nil {
		sendErr := replyToMessage(bot, message, cannotParseTeamKillMessage)
		if sendErr != nil {
			return sendErr
		}
		return nil
	}

	response, processErr := c.teamKillService.AddTeamKill(request, sourceFromMessage(message))
	if processErr != nil {
		switch processErr.(type) {
		case *services.ErrMixedCharactersInName:
			castedErr := processErr.(*services.ErrMixedCharactersInName)
			msg := fmt.Sprintf("В имени \"%s\" есть кириллица и латиница", castedErr.Name)
			sendErr := replyToMessage(bot, message, msg)
			if sendErr != nil {
				return sendErr
			}
		default:
			return processErr
		}

		return nil
	}

	if response != "" {
		sendErr := replyToMessage(bot, message, response)
		log.Error().Err(sendErr).Msg("Error while trying to send message")
		if sendErr != nil {
			return sendErr
		}
	}

	return nil
}

func (c *TeamKillCommand) IsCommand(text string) bool {
	return strings.Contains(text, teamKillMessageKey)
}

func (c *TeamKillCommand) GetHelp() string {
	return "<name> #teamkill <name> - записать тимкилл"
}

type ShowKillersCommand struct {
	teamKillService *services.TeamKillService
}

func (c *ShowKillersCommand) Process(bot *tgbotapi.BotAPI, message *tgbotapi.Message) error {
	log.Info().Msg("processing message show killers")

	topKillersLog, err := c.teamKillService.GetTopKillers(sourceFromMessage(message))
	if err != nil {
		replyToMessage(bot, message, cannotProcessMessage)
		return err
	}

	var rows []string
	for _, killLog := range topKillersLog {
		row := fmt.Sprintf("Убивца: %s; Намолотил: %d", killLog.Name, killLog.KillCount)
		rows = append(rows, row)
	}
	msg := strings.Join(rows, "\n")
	if msg == "" {
		msg = "Пока никто никого не убил"
	}

	sendErr := sendMessage(bot, message.Chat.ID, msg)
	if sendErr != nil {
		return sendErr
	}

	return nil
}

func (c *ShowKillersCommand) IsCommand(text string) bool {
	return strings.HasPrefix(text, showKillers)
}

func (c *ShowKillersCommand) GetHelp() string {
	return fmt.Sprintf("%s - показать список тимкиллов", showKillers)
}

type ShowVictimsCommand struct {
	teamKillService *services.TeamKillService
}

func (c *ShowVictimsCommand) Process(bot *tgbotapi.BotAPI, message *tgbotapi.Message) error {
	log.Info().Msg("processing message show victims")

	topVictimsLog, err := c.teamKillService.GetTopVictims(sourceFromMessage(message))
	if err != nil {
		replyErr := replyToMessage(bot, message, cannotProcessMessage)
		log.Error().Err(replyErr).Msg("Error while trying to send message about error, while processing command")
		return err
	}

	var rows []string
	for _, victimLog := range topVictimsLog {
		row := fmt.Sprintf("Неудачник: %s; Наумирал: %d", victimLog.Name, victimLog.DeathsCount)
		rows = append(rows, row)
	}
	msg := strings.Join(rows, "\n")
	if msg == "" {
		msg = "Пока никто никого не убил"
	}

	sendErr := sendMessage(bot, message.Chat.ID, msg)
	if sendErr != nil {
		return sendErr
	}

	return nil
}

func (c *ShowVictimsCommand) IsCommand(text string) bool {
	return strings.HasPrefix(text, showVictims)
}

func (c *ShowVictimsCommand) GetHelp() string {
	return fmt.Sprintf("%s - вывести список жертв", showVictims)
}

type RepeatCommand struct {
	teamKillService *services.TeamKillService
}

func (c *RepeatCommand) Process(bot *tgbotapi.BotAPI, message *tgbotapi.Message) error {
	log.Info().Msg("processing message repeat")

	if message.ReplyToMessage == nil {
		sendErr := replyToMessage(bot, message, "Нечего повторять, ну")
		if sendErr != nil {
			return sendErr
		}
		return nil
	}

	request := services.NewTeamKillFromText(message.ReplyToMessage.Text)
	if request == nil {
		sendErr := replyToMessage(bot, message.ReplyToMessage, cannotParseTeamKillMessage)
		if sendErr != nil {
			return sendErr
		}
		return nil
	}

	response, processErr := c.teamKillService.AddTeamKill(request, sourceFromMessage(message))
	if processErr != nil {
		return processErr
	}

	if response != "" {
		sendErr := replyToMessage(bot, message.ReplyToMessage, response)
		if sendErr != nil {
			return sendErr
		}
	}

	return nil
}

func (c *RepeatCommand) IsCommand(text string) bool {
	return strings.HasPrefix(text, repeat)
}

func (c *RepeatCommand) GetHelp() string {
	return fmt.Sprintf("%s - обработать предыдущее сообщение об убийстве", repeat)
}

type HelpCommand struct {
	commands []ICommand
}

func (c *HelpCommand) IsCommand(text string) bool {
	return strings.HasPrefix(text, help)
}

func (c *HelpCommand) GetHelp() string {
	return fmt.Sprintf("%s - вывести помощь", help)
}

func (c *HelpCommand) Process(bot *tgbotapi.BotAPI, message *tgbotapi.Message) error {
	log.Info().Msg("processing message help")

	var rows []string
	for _, command := range c.commands {
		rows = append(rows, command.GetHelp())
	}
	msg := strings.Join(rows, "\n")

	sendErr := sendMessage(bot, message.Chat.ID, msg)
	if sendErr != nil {
		return sendErr
	}

	return nil
}

func (c *HelpCommand) SetCommands(commands []ICommand) {
	c.commands = commands
}

type BotKillCommand struct {
	service *services.BotKillService
}

func (c *BotKillCommand) Process(bot *tgbotapi.BotAPI, message *tgbotapi.Message) error {
	log.Info().Msg("processing message bot kill")
	request := services.NewBotKillRequest(message.Text)
	if request == nil {
		sendErr := replyToMessage(bot, message, cannotParseTeamKillMessage)
		if sendErr != nil {
			return sendErr
		}
		return nil
	}

	response, processErr := c.service.AddBotKill(request, sourceFromMessage(message))
	if processErr != nil {
		return processErr
	}

	if response != "" {
		sendErr := replyToMessage(bot, message, response)
		if sendErr != nil {
			return sendErr
		}
	}

	return nil
}

func (c *BotKillCommand) IsCommand(text string) bool {
	return strings.Contains(text, botKillMessageKey)
}

func (c *BotKillCommand) GetHelp() string {
	return "<name> #botkill -  записать убийство от бота"
}

type ShowBotVictimsCommand struct {
	service *services.BotKillService
}

func (c *ShowBotVictimsCommand) Process(bot *tgbotapi.BotAPI, message *tgbotapi.Message) error {
	log.Info().Msg("processing message show bot victims")

	topVictimsLog, err := c.service.GetTopBotVictims(sourceFromMessage(message))
	if err != nil {
		replyErr := replyToMessage(bot, message, cannotProcessMessage)
		log.Printf("Error while trying to send message about error, while processing command: %s", replyErr)
		return err
	}

	var rows []string
	for _, victimLog := range topVictimsLog {
		row := fmt.Sprintf("Мишень: %s; Наумирал от ботов: %d", victimLog.Name, victimLog.DeathsCount)
		rows = append(rows, row)
	}
	msg := strings.Join(rows, "\n")
	if msg == "" {
		msg = "Пока никто не умер от ботов"
	}

	sendErr := sendMessage(bot, message.Chat.ID, msg)
	if sendErr != nil {
		return sendErr
	}

	return nil
}

func (c *ShowBotVictimsCommand) IsCommand(text string) bool {
	return strings.HasPrefix(text, showBotVictims)
}

func (c *ShowBotVictimsCommand) GetHelp() string {
	return fmt.Sprintf("%s - вывести список жертв от бота", showBotVictims)
}
