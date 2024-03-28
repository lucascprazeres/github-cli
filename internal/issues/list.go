package issues

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
)

func List(repo, is, author string) (*IssuesSearchResult, error) {
	query := buildQuery(repo, is, author)
	response, err := http.Get(IssuesURL + "?q=" + query)
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("search query failed with status (%s)", response.Status)
	}

	var result IssuesSearchResult
	if err := json.NewDecoder(response.Body).Decode(&result); err != nil {
		return nil, err
	}

	return &result, nil
}

func buildQuery(repo, is, author string) string {
	query := fmt.Sprintf("repo:%s is:issue", repo)

	if is != "" {
		query += fmt.Sprintf(" is:%s", is)
	}

	if author != "" {
		query += fmt.Sprintf(" author:%s", author)
	}

	return url.QueryEscape(query)
}
