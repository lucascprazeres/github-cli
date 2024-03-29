package utils

import (
	"fmt"
	"net/url"
)

func BuildQuery(repo, is, author string) string {
	query := fmt.Sprintf("repo:%s is:issue", repo)

	if is != "" {
		query += fmt.Sprintf(" is:%s", is)
	}

	if author != "" {
		query += fmt.Sprintf(" author:%s", author)
	}

	return url.QueryEscape(query)
}
