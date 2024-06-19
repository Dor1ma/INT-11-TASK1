package telegram

import (
	"INT-11-TASK1/internal/entity"
	gitlabutils "INT-11-TASK1/internal/gitlab"
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/xanzy/go-gitlab"
	"log"
)

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

func HandleTelegramUpdates(bot *tgbotapi.BotAPI, git *gitlab.Client, config *entity.Config) {
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
		case "setproject":
			HandleSetProjectCommand(update.Message, config, bot)
		case "check":
			gitlabutils.CheckMergeRequests(git, bot, config.TelegramChatID, config.GitlabProjectID)
		default:
			SendTelegramMessage(bot, update.Message.Chat.ID, "Unknown command: "+update.Message.Command())
		}
	}
}
