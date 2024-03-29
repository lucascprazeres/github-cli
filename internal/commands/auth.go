package commands

import (
	"github.com/lucascprazeres/github-cli/internal/file"
	"github.com/lucascprazeres/github-cli/internal/services"
)

func Auth() error {
	github := services.NewGithubService()

	token, err := github.GetToken()
	if err != nil {
		return err
	}

	file.Write("credentials.json", token)

	return nil
}
