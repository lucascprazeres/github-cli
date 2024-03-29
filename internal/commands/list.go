package commands

import (
	"fmt"

	"github.com/lucascprazeres/github-cli/internal/logging"
	"github.com/lucascprazeres/github-cli/internal/services"
)

func List(repo, is, author string) error {
	github := services.NewGithubService()
	issues, err := github.GetIssues(repo, is, author)
	if err != nil {
		return err
	}

	if is != "" && is != "open" && is != "closed" {
		return fmt.Errorf("the only valid options for 'is' are 'open' and 'closed'")
	}

	logging.Prompt("%v issues found:\n", issues.TotalCount)

	if issues.TotalCount == 0 {
		return nil
	}

	logging.Info("[#%-5s]\t %15s\t %.55s\n", "ISSUE", "AUTHOR", "TITLE")
	for _, item := range issues.Items {
		logging.Success("[#%-5d]\t %15.15s\t %.55s\n", item.Number, item.User.Login, item.Title)
	}

	return nil
}
