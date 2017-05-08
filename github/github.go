package main

import (
	"context"
	"fmt"
	//"math/rand"
	"net/http"
	"os"
	//"strconv"
	//"time"

	"github.com/google/go-github/github"
	github_goth "github.com/markbates/goth/providers/github"
	"golang.org/x/oauth2"
	ogithub "golang.org/x/oauth2/github"

	. "github.com/y0ssar1an/q"
)

import _ "github.com/joho/godotenv/autoload"

var SlackWebhook string

func init() {
	SlackWebhook = os.Getenv("SLACK_WEBHOOK")

}

func main() {
	clientId := os.Getenv("GITHUB_OAUTH_CLIENT_ID")
	clientSecret := os.Getenv("GITHUB_OAUTH_CLIENT_SECRET")
	redirectURL := os.Getenv("GITHUB_OAUTH_REDIRECT")

	ctx := context.Background()
	conf := &oauth2.Config{
		ClientID:     clientId,
		ClientSecret: clientSecret,
		Scopes:       []string{"admin:repo_hook"},
		Endpoint:     ogithub.Endpoint,
		RedirectURL:  redirectURL,
	}

	// Redirect user to consent page to ask for permission
	// for the scopes specified above.
	//rand.Seed(time.Now().UnixNano())
	//state := strconv.Atoi(rand.Int())
	provider := github_goth.New(clientId, clientSecret, redirectURL, "admin:repo_hook")
	session, err := provider.BeginAuth("state")

	session.

	url := conf.AuthCodeURL("state", oauth2.AccessTypeOffline)
	fmt.Printf("Visit the URL for the auth dialog: %v", url)

	// Use the authorization code that is pushed to the redirect
	// URL. Exchange will do the handshake to retrieve the
	// initial access token. The HTTP Client returned by
	// conf.Client will refresh the token as necessary.
	var code string
	if _, err := fmt.Scan(&code); err != nil {
		//log.Fatal(err)
	}
	tok, err := conf.Exchange(ctx, code)
	if err != nil {
		//log.Fatal(err)
	}

	//oauth2
	tc := conf.Client(ctx, tok)
	//client.Get("...")

	//token := os.Getenv("GITHUB_TOKEN")
	//Q(token)

	//fmt.Println("You can do anything!")
	//Q("Lets debug some shit")

	//ctx := context.Background()
	//ts := oauth2.StaticTokenSource(
	//&oauth2.Token{AccessToken: token},
	//)
	//tc := oauth2.NewClient(ctx, ts)

	client := github.NewClient(tc)

	// list all repositories for the authenticated user
	//hooks, _, err := client.Repositories.List(ctx, "adamryman", nil)
	//hooks, _, err := client.ListServiceHooks(ctx)
	//_ = err
	//Q(err)

	//hooks, _, err := client.Repositories.ListHooks(ctx, "adamryman", "kit", nil)
	hook, err := WatchRepo(ctx, client, "adamryman", "kit")
	if err != nil {
		Q(err)
		Q("ERROR")
		os.Exit(1)
	}
	Q(hook)
	//Q(hooks[0])

}

type auth struct{}

func (s *auth) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	r.FormValue("code")

}

func ServeHTTP(w http.ResponseWriter, r *http.Request) {
	//s.webhookSecretKey
	payload, err := github.ValidatePayload(r, nil)
	if err != nil {
		return
	}
	event, err := github.ParseWebHook(github.WebHookType(r), payload)
	if err != nil {
		return
	}
	switch event := event.(type) {
	case *github.PullRequestEvent:
		switch action := event.GetAction(); action {
		case "open":
			// To put into slack message
			event.PullRequest.GetURL()
			// To store so that when it is closed @lgtm can react to it
			event.PullRequest.GetID()

		case "closed":
			//// to find the PR in chat to react to it
			event.PullRequest.GetID()

		}
	}
}

func WatchRepo(ctx context.Context, client *github.Client, owner, repo string) (*github.Hook, error) {
	req := new(github.Hook)
	name := "web"
	active := true
	req = &github.Hook{
		Name:   &name,
		Active: &active,
		Events: []string{"pull_request"},
		Config: map[string]interface{}{
			"url":          SlackWebhook,
			"content_type": "json",
		},
	}
	hook, _, err := client.Repositories.CreateHook(ctx, owner, repo, req)
	return hook, err
}
