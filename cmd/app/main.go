package main

import (
	"INT-11-TASK1/internal/entity"
	"INT-11-TASK1/internal/gitlab"
	"INT-11-TASK1/internal/telegram"
	"github.com/robfig/cron"
	"log"
)

func main() {
	config, err := entity.LoadConfig()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
		return
	}

	git, err := gitlab.CreateGitlabClient(config.GitlabToken)
	if err != nil {
		log.Fatalf("Failed to create GitLab client: %v", err)
	}

	bot, err := telegram.CreateTelegramBot(config.TelegramBotToken)
	if err != nil {
		log.Fatalf("Failed to create Telegram bot: %v", err)
	}

	botCtx := &telegram.BotContext{}

	c := cron.New()
	err = c.AddFunc(config.ReminderFrequency, func() { gitlab.CheckMergeRequests(git, bot, botCtx, config.GitlabProjectID) })
	if err != nil {
		return
	}
	c.Start()

	telegram.HandleTelegramUpdates(bot, git, config, botCtx)
}
