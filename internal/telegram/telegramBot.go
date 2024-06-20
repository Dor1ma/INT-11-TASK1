package telegram

import (
	"fmt"
	"github.com/Dor1ma/INT-11-TASK1/internal/entity"
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

func HandleChangeProjectCommand(message *tgbotapi.Message, config *entity.Config, bot *tgbotapi.BotAPI) {
	args := message.CommandArguments()
	if args == "" {
		SendTelegramMessage(bot, message.Chat.ID, "Пожалуйста, укажите project ID вместе с командой")
		return
	}

	config.GitlabProjectID = args
	SendTelegramMessage(bot, message.Chat.ID, "Project ID обновлен: "+config.GitlabProjectID)
}

func HandleHelpCommand(message *tgbotapi.Message, bot *tgbotapi.BotAPI) {
	helpMessage := "Доступные команды:\n" +
		"/change_project <projectID> - Изменение GitLab project ID\n" +
		"/check - Принудительная проверка новых Merge Request'ов\n" +
		"/help - Вывод описания команд\n" +
		"/start - Запуск бота"
	SendTelegramMessage(bot, message.Chat.ID, helpMessage)
}

func HandleStartCommand(message *tgbotapi.Message, botCtx *BotContext, bot *tgbotapi.BotAPI) {
	botCtx.Lock()
	defer botCtx.Unlock()
	botCtx.ChatID = message.Chat.ID
	SendTelegramMessage(bot, message.Chat.ID, "Привет! Бот успешно запущен!\n"+
		"Со списком команд вы можете ознакомиться, использовав команду /help")
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
		case "change_project":
			HandleChangeProjectCommand(update.Message, config, bot)
		case "check":
			CheckMergeRequests(git, bot, botCtx, config.GitlabProjectID, true)
		case "help":
			HandleHelpCommand(update.Message, bot)
		default:
			SendTelegramMessage(bot, update.Message.Chat.ID, "Неизвестная команда")
		}
	}
}
