package telegram

import (
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

func CheckMergeRequests(git *gitlab.Client, bot *tgbotapi.BotAPI, botCtx *BotContext, projectID string) {
	log.Println("Starting check for new merge requests...")
	pID, err := strconv.Atoi(projectID)
	if err != nil {
		log.Printf("Invalid project ID: %v", err)
		return
	}

	mergeRequests, _, err := git.MergeRequests.ListProjectMergeRequests(pID, &gitlab.ListProjectMergeRequestsOptions{
		State: gitlab.String("opened"),
	})
	if err != nil {
		log.Printf("Failed to list merge requests: %v", err)
		return
	}

	botCtx.Lock()
	defer botCtx.Unlock()

	if botCtx.ChatID == 0 {
		log.Printf("Chat ID is not set. Cannot send messages.")
		return
	}

	for _, mr := range mergeRequests {
		message := fmt.Sprintf("Новый Merge Request:\nTitle: %s\nURL: %s", mr.Title, mr.WebURL)
		log.Printf("Sending message: %s", message)
		SendTelegramMessage(bot, botCtx.ChatID, message)
	}
	log.Println("Finished checking for merge requests.")
}
