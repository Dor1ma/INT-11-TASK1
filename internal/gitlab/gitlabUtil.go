package gitlab

import (
	"INT-11-TASK1/internal/telegram"
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/xanzy/go-gitlab"
	"log"
	"strconv"
)

func CreateGitlabClient(token string) (*gitlab.Client, error) {
	client, err := gitlab.NewClient(token)
	if err != nil {
		return nil, fmt.Errorf("failed to create Gitlab client: %v", err)
	}
	return client, nil
}

func CheckMergeRequests(git *gitlab.Client, bot *tgbotapi.BotAPI, chatID int64, projectID string) {
	pID, err := strconv.Atoi(projectID)
	if err != nil {
		log.Printf("Invalid project id %v", err)
		return
	}

	mergeRequests, _, err := git.MergeRequests.ListProjectMergeRequests(pID, &gitlab.ListProjectMergeRequestsOptions{
		State: gitlab.String("opened"),
	})
	if err != nil {
		log.Printf("Failed to list project merge requests: %v", err)
		return
	}

	for _, mr := range mergeRequests {
		message := fmt.Sprintf("Новый Merge Request:\nTitle: %s\nURL: %s", mr.Title, mr.WebURL)
		telegram.SendTelegramMessage(bot, chatID, message)
	}
}
