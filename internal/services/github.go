package services

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/lucascprazeres/github-cli/internal/config/env"
	"github.com/lucascprazeres/github-cli/internal/logging"
	"github.com/lucascprazeres/github-cli/internal/utils"
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

const IssuesURL = "https://api.github.com/search/issues"

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
	query := utils.BuildQuery(repo, is, author)
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
