package github

import (
	"fmt"

	"github.com/tibia-oce/discord-bot/src/logger"
)

type GitHubClient struct{}

func NewGitHubClient() *GitHubClient {
	return &GitHubClient{}
}

func (g *GitHubClient) CreateIssue(repository, issueType, title, description, imageLink string) {
	issueBody := fmt.Sprintf(
		"**Issue Type:** %s\n**Description:** %s",
		issueType, description,
	)
	if imageLink != "" {
		issueBody += fmt.Sprintf("\n**Image Link:** %s", imageLink)
	}

	logger.Info(fmt.Sprintf("GitHub Issue Created in Repository: %s\nTitle: %s\nBody:\n%s", repository, title, issueBody))
}
