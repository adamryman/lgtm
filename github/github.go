package github

import (
	"context"
	//"fmt"
	"math"
	//"math/rand"
	"net/http"
	"os"
	//"strconv"
	//"time"

	"github.com/google/go-github/github"
	//"github.com/markbates/goth"
	gothic "github.com/markbates/goth/gothic"
	//github_goth "github.com/markbates/goth/providers/github"

	"github.com/gorilla/sessions"

	"golang.org/x/oauth2"
	//ogithub "golang.org/x/oauth2/github"
	. "github.com/y0ssar1an/q"
)

import _ "github.com/joho/godotenv/autoload"

var SlackWebhook string

func init() {
	SlackWebhook = os.Getenv("SLACK_WEBHOOK")
	Q(SlackWebhook)

}

func init() {
	store := sessions.NewFilesystemStore(os.TempDir(), []byte("goth-example"))

	// set the maxLength of the cookies stored on the disk to a larger number to prevent issues with:
	// securecookie: the value is too long
	// when using OpenID Connect , since this can contain a large amount of extra information in the id_token

	// Note, when using the FilesystemStore only the session.ID is written to a browser cookie, so this is explicit for the storage on disk
	store.MaxLength(math.MaxInt64)

	gothic.Store = store
}

type auth struct{}

type AuthServer struct {
}

//func AuthCallback() http.Handler {

//return
//}

func ServeHTTPCallback(w http.ResponseWriter, r *http.Request) {
	user, err := gothic.CompleteUserAuth(w, r)
	if err != nil {
		return
	}
	_ = user.AccessToken
	// TODO: Store user.AccessToken, assoitate with user
	// TODO: Maybe hook in a way to response to user?

}

func ServeHTTPAuthenicate(w http.ResponseWriter, r *http.Request) {
	gothic.BeginAuthHandler(w, r)
}

// Wehook endpoint
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

type authError struct {
	err error
}

func WatchRepo(ctx context.Context, token, owner, repo string) (*github.Hook, error) {
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

	ts := oauth2.StaticTokenSource(
		&oauth2.Token{
			AccessToken: token,
		},
	)

	tc := oauth2.NewClient(ctx, ts)

	client := github.NewClient(tc)
	hook, _, err := client.Repositories.CreateHook(ctx, owner, repo, req)
	return hook, err
}

/*
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
	gothic.GetProviderName = func(req *http.Request) (string, error) { return "github", nil }
	provider := github_goth.New(clientId, clientSecret, redirectURL, "admin:repo_hook")
	goth.UseProviders(provider)

	// Must get auth token out of database to make a tc

	//client := github.NewClient(tc)

	//// list all repositories for the authenticated user
	////hooks, _, err := client.Repositories.List(ctx, "adamryman", nil)
	////hooks, _, err := client.ListServiceHooks(ctx)
	////_ = err
	////Q(err)

	////hooks, _, err := client.Repositories.ListHooks(ctx, "adamryman", "kit", nil)
	//hook, err := WatchRepo(ctx, client, "adamryman", "kit")
	//if err != nil {
	//Q(err)
	//Q("ERROR")
	//os.Exit(1)
	//}
	//Q(hook)
	//Q(hooks[0])

}
*/
