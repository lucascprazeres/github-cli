package commands

import (
	"encoding/json"
	"fmt"
	"strings"

	"github-cli/internal/logging"
	"github-cli/internal/services"
)

func Create(repo, title, body string) error {
	if err := validateParams(repo, title, body); err != nil {
		return err
	}

	github := services.NewGithubService()
	result, err := github.CreateIssue(repo, title, body)
	if err != nil {
		return err
	}

	formattedResult, err := format(result)
	if err != nil {
		return err
	}

	logging.Success("%s\n", formattedResult)
	return nil
}

func validateParams(repo, title, body string) error {
	if !strings.Contains(repo, "/") {
		return fmt.Errorf("repository be passed as [owner/repo]")
	}

	if title == "" {
		return fmt.Errorf("title is required")
	}

	if body == "" {
		return fmt.Errorf("body is required")
	}

	return nil
}

func format(data *services.CreateIssueResult) (string, error) {
	result, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		return "", fmt.Errorf("error formatting json response")
	}
	return string(result), nil
}
