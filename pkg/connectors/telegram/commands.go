package telegram

import (
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"log"
	"strings"
	"teamkillbot/pkg/services"
)

type ICommand interface {
	Process(bot *tgbotapi.BotAPI, message *tgbotapi.Message) error
	IsCommand(string) bool
	GetHelp() string
}

type CommandEcho struct{}

func (c *CommandEcho) Process(bot *tgbotapi.BotAPI, message *tgbotapi.Message) error {
	return replyToMessage(bot, message, message.Text)
}

func (c *CommandEcho) IsCommand(text string) bool {
	return strings.HasPrefix(text, echoKey)
}

func (c *CommandEcho) GetHelp() string {
	return "PTHD:echo - simple echo text"
}

type TeamKillCommand struct {
	teamKillService *services.TeamKillService
}

func (c *TeamKillCommand) Process(bot *tgbotapi.BotAPI, message *tgbotapi.Message) error {
	log.Printf("processing message team kill")
	request := services.NewTeamKillFromText(message.Text)
	if request == nil {
		sendErr := replyToMessage(bot, message, cannotParseTeamKillMessage)
		if sendErr != nil {
			return sendErr
		}
		return nil
	}

	response, processErr := c.teamKillService.ProcessTeamKill(request)
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

func (c *TeamKillCommand) IsCommand(text string) bool {
	return strings.Contains(text, teamKillMessageKey)
}

func (c *TeamKillCommand) GetHelp() string {
	return "<name> #teamkill <name> - log team kill"
}

type ShowKillersCommand struct {
	teamKillService *services.TeamKillService
}

func (c *ShowKillersCommand) Process(bot *tgbotapi.BotAPI, message *tgbotapi.Message) error {
	topKillersLog, err := c.teamKillService.ProcessGetTopKillers()
	if err != nil {
		replyErr := replyToMessage(bot, message, cannotProcessMessage)
		log.Printf("Error while trying to send message about error, while processing command: %s", replyErr)
		return err
	}

	var rows []string
	for _, killLog := range topKillersLog {
		row := fmt.Sprintf("Убивца: %s; Намолотил: %d", killLog.Name, killLog.KillCount)
		rows = append(rows, row)
	}
	msg := strings.Join(rows, "\n")

	sendErr := sendMessage(bot, message.Chat.ID, msg)
	if sendErr != nil {
		log.Printf("Error while trying to send message: %s", sendErr)
	}

	return nil
}

func (c *ShowKillersCommand) IsCommand(text string) bool {
	return strings.HasPrefix(text, showKillers)
}

func (c *ShowKillersCommand) GetHelp() string {
	return "PTHD:покажи убивц - show killers log"
}

type ShowVictimsCommand struct {
	teamKillService *services.TeamKillService
}

func (c *ShowVictimsCommand) Process(bot *tgbotapi.BotAPI, message *tgbotapi.Message) error {
	topVictimsLog, err := c.teamKillService.ProcessGetTopVictims()
	if err != nil {
		replyErr := replyToMessage(bot, message, cannotProcessMessage)
		log.Printf("Error while trying to send message about error, while processing command: %s", replyErr)
		return err
	}

	var rows []string
	for _, victimLog := range topVictimsLog {
		row := fmt.Sprintf("Неудачник: %s; Наумирал: %d", victimLog.Name, victimLog.DeathsCount)
		rows = append(rows, row)
	}
	msg := strings.Join(rows, "\n")

	sendErr := sendMessage(bot, message.Chat.ID, msg)
	if sendErr != nil {
		log.Printf("Error while trying to send message: %s", sendErr)
	}

	return nil
}

func (c *ShowVictimsCommand) IsCommand(text string) bool {
	return strings.HasPrefix(text, showVictims)
}

func (c *ShowVictimsCommand) GetHelp() string {
	return "PTHD:покажи неудачников - show victims log"
}

type RepeatCommand struct {
	teamKillService *services.TeamKillService
}

func (c *RepeatCommand) Process(bot *tgbotapi.BotAPI, message *tgbotapi.Message) error {
	log.Printf("processing message repeat")

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

	response, processErr := c.teamKillService.ProcessTeamKill(request)
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
	return "PTHD:repeat - show victims log"
}

type HelpCommand struct {
	commands []ICommand
}

func (c *HelpCommand) IsCommand(text string) bool {
	return strings.HasPrefix(text, "PTHD:help")
}

func (c *HelpCommand) GetHelp() string {
	return "PTHD:help - help list"
}

func (c *HelpCommand) Process(bot *tgbotapi.BotAPI, message *tgbotapi.Message) error {
	var rows []string
	for _, command := range c.commands {
		rows = append(rows, command.GetHelp())
	}
	msg := strings.Join(rows, "\n")

	sendErr := sendMessage(bot, message.Chat.ID, msg)
	if sendErr != nil {
		log.Printf("Error while trying to send message: %s", sendErr)
	}

	return nil
}

func (c *HelpCommand) SetCommands(commands []ICommand) {
	c.commands = commands
}
