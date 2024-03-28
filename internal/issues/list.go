package issues

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
)

func List(repo, is, author string) (*IssuesSearchResult, error) {
	raw := fmt.Sprintf("repo:%s is:issue", repo)

	if len(is) > 0 {
		raw += " is:" + is
	}

	if len(author) > 0 {
		raw += " author:" + author
	}

	q := url.QueryEscape(raw)
	r, err := http.Get(IssuesURL + "?q=" + q)
	if err != nil {
		return nil, err
	}
	defer r.Body.Close()

	if r.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("search query failed with status (%s)", r.Status)
	}

	var result IssuesSearchResult
	if err := json.NewDecoder(r.Body).Decode(&result); err != nil {
		return nil, err
	}

	return &result, nil
}
