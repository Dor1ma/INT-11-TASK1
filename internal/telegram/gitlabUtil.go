package telegram

import (
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/xanzy/go-gitlab"
	"log"
	"strconv"
	"strings"
)

func CreateGitlabClient(token string) (*gitlab.Client, error) {
	client, err := gitlab.NewClient(token)
	if err != nil {
		return nil, fmt.Errorf("failed to create Gitlab client: %v", err)
	}
	return client, nil
}

var notifiedMRs = make(map[int]bool)

func CheckMergeRequests(git *gitlab.Client, bot *tgbotapi.BotAPI, botCtx *BotContext, projectID string, forceCheck bool) {
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
		log.Printf("Chat ID is not set. Cannot send messages")
		return
	}

	newMRs := []string{}
	unresolvedMRs := []string{}

	for _, mr := range mergeRequests {
		message := fmt.Sprintf("Title: %s\nURL: %s", mr.Title, mr.WebURL)
		if _, notified := notifiedMRs[mr.IID]; !notified {
			newMRs = append(newMRs, message)
			notifiedMRs[mr.IID] = true
		} else {
			unresolvedMRs = append(unresolvedMRs, message)
		}
	}

	if len(newMRs) > 0 {
		fullMessage := "You have new Merge Requests:\n" + strings.Join(newMRs, "\n\n")
		log.Printf("Sending message: %s", fullMessage)
		SendTelegramMessage(bot, botCtx.ChatID, fullMessage)
	}

	if len(unresolvedMRs) > 0 {
		fullMessage := "You have unresolved Merge Requests:\n" + strings.Join(unresolvedMRs, "\n\n")
		log.Printf("Sending message: %s", fullMessage)
		SendTelegramMessage(bot, botCtx.ChatID, fullMessage)
	} else if forceCheck && len(newMRs) == 0 {
		SendTelegramMessage(bot, botCtx.ChatID, "There are no new Merge Requests for this project")
	}

	for iid := range notifiedMRs {
		mr, _, err := git.MergeRequests.GetMergeRequest(pID, iid, nil)
		if err != nil {
			log.Printf("Failed to get merge request: %v", err)
			continue
		}

		if mr.State == "merged" || mr.State == "closed" {
			delete(notifiedMRs, iid)
		}
	}

	log.Println("Finished checking for merge requests")
}
