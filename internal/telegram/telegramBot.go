package telegram

import (
	"INT-11-TASK1/internal/entity"
	gitlabutils "INT-11-TASK1/internal/gitlab"
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/xanzy/go-gitlab"
	"log"
	"sync"
)

type BotContext struct {
	sync.Mutex
	ChatID int64
}

func CreateTelegramBot(token string) (*tgbotapi.BotAPI, error) {
	bot, err := tgbotapi.NewBotAPI(token)
	if err != nil {
		return nil, fmt.Errorf("failed to create Telegram bot: %w", err)
	}
	return bot, nil
}

func SendTelegramMessage(bot *tgbotapi.BotAPI, chatID int64, message string) {
	msg := tgbotapi.NewMessage(chatID, message)
	_, err := bot.Send(msg)
	if err != nil {
		log.Printf("Failed to send message: %v", err)
	}
}

func HandleSetProjectCommand(message *tgbotapi.Message, config *entity.Config, bot *tgbotapi.BotAPI) {
	args := message.CommandArguments()
	if args == "" {
		SendTelegramMessage(bot, message.Chat.ID, "Please provide a project ID.")
		return
	}

	config.GitlabProjectID = args
	SendTelegramMessage(bot, message.Chat.ID, "Project ID updated: "+config.GitlabProjectID)
}

func HandleHelpCommand(message *tgbotapi.Message, bot *tgbotapi.BotAPI) {
	helpMessage := "Available commands:\n" +
		"/setproject <projectID> - Set the GitLab project ID\n" +
		"/check - Manually check for merge requests\n" +
		"/help - Show this help message\n" +
		"/start - Starting the bot"
	SendTelegramMessage(bot, message.Chat.ID, helpMessage)
}

func HandleStartCommand(message *tgbotapi.Message, botCtx *BotContext, bot *tgbotapi.BotAPI) {
	botCtx.Lock()
	defer botCtx.Unlock()
	botCtx.ChatID = message.Chat.ID
	SendTelegramMessage(bot, message.Chat.ID, "Hello! You have successfully started the bot.")
}

func HandleTelegramUpdates(bot *tgbotapi.BotAPI, git *gitlab.Client, config *entity.Config, botCtx *BotContext) {
	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates, err := bot.GetUpdatesChan(u)
	if err != nil {
		log.Fatalf("Failed to set up update channel: %v", err)
	}

	for update := range updates {
		if update.Message == nil || !update.Message.IsCommand() {
			continue
		}

		switch update.Message.Command() {
		case "start":
			HandleStartCommand(update.Message, botCtx, bot)
		case "setproject":
			HandleSetProjectCommand(update.Message, config, bot)
		case "check":
			gitlabutils.CheckMergeRequests(git, bot, botCtx, config.GitlabProjectID)
		case "help":
			HandleHelpCommand(update.Message, bot)
		default:
			SendTelegramMessage(bot, update.Message.Chat.ID, "I don't know that command")
		}
	}
}
