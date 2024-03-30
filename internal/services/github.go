package services

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"time"

	"github-cli/internal/config/env"
	"github-cli/internal/file"
	"github-cli/internal/logging"

	"github.com/pkg/browser"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/github"
)

type Github struct {
	config *oauth2.Config
}

type IssuesSearchResult struct {
	TotalCount int `json:"total_count"`
	Items      []*Issue
}

type Issue struct {
	Number    int
	HTMLURL   string `json:"html_url"`
	Title     string
	State     string
	User      *User
	CreatedAt time.Time `json:"created_at"`
	Body      string    // em formato Markdown
}

type User struct {
	Login   string
	HTMLURL string `json:"html_url"`
}

type CreateIssueResult struct {
	Title         string
	Body          string
	Number        int
	URL           string
	RepositoryURL string `json:"repository_url"`
	User          *User
	CreatedAt     time.Time `json:"created_at"`
}

func NewGithubService() *Github {
	service := &Github{}
	service.config = &oauth2.Config{
		ClientID:     env.Get("GITHUB_CLIENT_ID"),
		ClientSecret: env.Get("GITHUB_CLIENT_SECRET"),
		Scopes:       []string{"read:org", "read:user", "read:project", "public_repo", "gist"},
		Endpoint:     github.Endpoint,
		RedirectURL:  "http://localhost:9999/oauth/callback",
	}
	return service
}

func (g *Github) GetToken() (*oauth2.Token, error) {
	// create a server
	server := &http.Server{Addr: ":9999"}

	// create a channel to receive authorization code
	codeChan := make(chan string)

	// github callback route
	http.HandleFunc("/oauth/callback", func(w http.ResponseWriter, r *http.Request) {
		code := r.URL.Query().Get("code")
		codeChan <- code
	})

	go func() {
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("failed to start the server: %v", err)
		}
	}()

	// get the OAuth authorization URL
	url := g.config.AuthCodeURL("state", oauth2.AccessTypeOffline)

	logging.Info("Your browser has been opened to visit:\n\n%s\n\n", url)

	// Redirect user to consent page to ask for permission
	// for the scopes specified above
	if err := browser.OpenURL(url); err != nil {
		return nil, fmt.Errorf("failed to open browser for authentication: %s", err.Error())
	}

	// wait for the authorization code to be received
	code := <-codeChan

	// exchange the authorization code for an access token
	token, err := g.config.Exchange(context.Background(), code)
	if err != nil {
		return nil, fmt.Errorf("failed to exchange authorization code for token: %v", err)
	}

	if !token.Valid() {
		return nil, fmt.Errorf("can't get source information without accessToken: %v", err)
	}

	if err := server.Shutdown(context.Background()); err != nil {
		return nil, fmt.Errorf("failed to shut down server: %v", err)
	}

	logging.Success("Authenticated successfuly!\n")

	return token, nil
}

func (g *Github) GetIssues(repo, is, author string) (*IssuesSearchResult, error) {
	query := buildQuery(repo, is, author)
	res, err := http.Get("https://api.github.com/search/issues" + "?q=" + query)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("search query failed with status (%s)", res.Status)
	}

	var result IssuesSearchResult
	if err := json.NewDecoder(res.Body).Decode(&result); err != nil {
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

func (g *Github) CreateIssue(repo, title, body string) (*CreateIssueResult, error) {
	url := fmt.Sprintf("https://api.github.com/repos/%s/issues", repo)

	reqBody, err := buildRequestBody(title, body)
	if err != nil {
		return &CreateIssueResult{}, err
	}

	req, err := http.NewRequest(http.MethodPost, url, reqBody)
	if err != nil {
		return &CreateIssueResult{}, err
	}

	credentials, err := file.Read("credentials.json")
	if err != nil {
		return &CreateIssueResult{}, err
	}

	token := fmt.Sprintf("Bearer %s", credentials["access_token"])

	req.Header.Set("Accept", "application/vnd.github+json")
	req.Header.Set("Authorization", token)
	req.Header.Set("X-GitHub-Api-Version", "2022-11-28")

	client := http.Client{Timeout: 3 * time.Second}
	res, err := client.Do(req)
	if err != nil {
		return &CreateIssueResult{}, err
	}
	defer res.Body.Close()

	var result *CreateIssueResult
	if err := json.NewDecoder(res.Body).Decode(&result); err != nil {
		return &CreateIssueResult{}, err
	}

	return result, nil
}

func buildRequestBody(title, body string) (*bytes.Reader, error) {
	jsonBody, err := json.Marshal(map[string]string{"title": title, "body": body})
	if err != nil {
		return nil, fmt.Errorf("invalid params format")
	}
	return bytes.NewReader(jsonBody), nil
}
