package main

import (
	"github.com/Dor1ma/INT-11-TASK1/internal/entity"
	"github.com/Dor1ma/INT-11-TASK1/internal/telegram"
	"github.com/robfig/cron/v3"
	"log"
)

func main() {
	config, err := entity.LoadConfig()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
		return
	}

	git, err := telegram.CreateGitlabClient(config.GitlabToken)
	if err != nil {
		log.Fatalf("Failed to create GitLab client: %v", err)
	}

	bot, err := telegram.CreateTelegramBot(config.TelegramBotToken)
	if err != nil {
		log.Fatalf("Failed to create Telegram bot: %v", err)
	}

	botCtx := &telegram.BotContext{}

	c := cron.New()
	_, err = c.AddFunc(config.ReminderFrequency, func() {
		log.Println("Cron job started.")
		telegram.CheckMergeRequests(git, bot, botCtx, config.GitlabProjectID, false)
		log.Println("Cron job finished.")
	})
	if err != nil {
		log.Fatalf("Failed to add cron job: %v", err)
	}
	c.Start()

	log.Println("Bot started. Listening for commands...")

	telegram.HandleTelegramUpdates(bot, git, config, botCtx)
}
