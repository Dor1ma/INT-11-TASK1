package entity

import (
	"fmt"
	"github.com/joho/godotenv"
	"os"
)

type Config struct {
	GitlabToken       string
	GitlabProjectID   string
	TelegramBotToken  string
	TelegramChatID    string
	ReminderFrequency string
}

func LoadConfig() (*Config, error) {
	err := godotenv.Load()
	if err != nil {
		return nil, err
	}

	config := &Config{
		GitlabToken:       os.Getenv("GITLAB_TOKEN"),
		GitlabProjectID:   os.Getenv("GITLAB_PROJECT_ID"),
		TelegramBotToken:  os.Getenv("TELEGRAM_BOT_TOKEN"),
		TelegramChatID:    os.Getenv("TELEGRAM_CHAT_ID"),
		ReminderFrequency: os.Getenv("REMINDER_FREQUENCY"),
	}

	if config.GitlabToken == "" || config.GitlabProjectID == "" ||
		config.TelegramBotToken == "" || config.TelegramChatID == "" ||
		config.ReminderFrequency == "" {
		return nil, fmt.Errorf("missing one or more required environment variables")
	}

	return config, nil
}
