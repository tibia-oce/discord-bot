package githubapi

import (
	"fmt"

	"github.com/tibia-oce/discord-bot/src/logger"
)

type GitHubClient struct{}

func NewGitHubClient() *GitHubClient {
	return &GitHubClient{}
}

func (g *GitHubClient) CreateIssue(title, body string) {
	logger.Info(fmt.Sprintf("GitHub Issue Created - Title: %s\nBody: %s", title, body))
}
