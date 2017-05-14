package lgtm

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"

	"github.com/StudentRND/lgtm/bot"
	"github.com/StudentRND/lgtm/github"
	sql "github.com/StudentRND/lgtm/sqlite"

	"github.com/pkg/errors"
	. "github.com/y0ssar1an/q"
)

var DefaultConfig = Config{
	SQLiteDB: os.Getenv("SQLITE3"),

	Scheme: os.Getenv("SCHEME"),
	Addr:   os.Getenv("ADDR"),
	Port:   os.Getenv("PORT"),

	SlackToken:   os.Getenv("SLACK_TOKEN"),
	SlackChannel: os.Getenv("SLACK_CHANNEL"),
	SlackbotId:   os.Getenv("SLACKBOT_ID"),

	GithubOauthClientId:     os.Getenv("GITHUB_OAUTH_CLIENT_ID"),
	GithubOauthClientSecret: os.Getenv("GITHUB_OAUTH_CLIENT_SECRET"),
}

type Config struct {
	// Path to SQLite database
	SQLiteDB string

	// e.g. http https
	Scheme string
	// e.g adamryman.com:5060 srnd.org xwl.me:12542
	Addr string

	// Webserver port
	// e.g. 5040 80 12542
	Port string

	// Slack Api Token
	SlackToken   string
	SlackChannel string
	SlackbotId   string

	// Github Auth
	GithubOauthClientId     string
	GithubOauthClientSecret string
}

var debug bool

func init() {
	Q(os.Getenv("DEBUG"))
	if len(os.Getenv("DEBUG")) > 0 {
		debug = true
	}
}

func Start(cfg Config) error {
	// TODO: Validate config
	_, err := strconv.Atoi(cfg.Port)
	if err != nil {
		return errors.Wrap(err, "cfg.Port not a number")
	}

	database, err := sql.Open(cfg.SQLiteDB)

	logger := log.New(os.Stdout, "slack-bot: ", log.Lshortfile|log.LstdFlags)

	lgtm, err := bot.Start(cfg.SlackToken, cfg.SlackbotId, cfg.SlackChannel, bot.SetLogger(logger))
	if err != nil {
		return err
	}
	if debug {
		lgtm.PostMessage("Rebooting....Done " + time.Now().String())
	}

	// handle chat events
	authurl := fmt.Sprintf("%s://%s/lgtm/authorize", cfg.Scheme, cfg.Addr)
	go handleBotEvents(authurl, lgtm, database)

	// start github handlers
	go func() {
		authCallbackURL := fmt.Sprintf("%s://%s/lgtm/authorize/callback", cfg.Scheme, cfg.Addr)
		github.InitAuth(cfg.GithubOauthClientId, cfg.GithubOauthClientSecret, authCallbackURL)

		github.PRWebhook = fmt.Sprintf("%s://%s/lgtm/webhook", cfg.Scheme, cfg.Addr)
		http.HandleFunc("/lgtm/authorize", github.AuthenicateHandler)
		http.HandleFunc("/lgtm/authorize/callback", github.AuthenticateCallbackHandler)
		http.HandleFunc("/lgtm/webhook", github.WebhookHandler)
		http.ListenAndServe(":"+cfg.Port, nil)
	}()

	// handle github events
	go handleGithubEvents(lgtm, database)

	// Wait for termination
	c := make(chan os.Signal, 1)
	signal.Notify(c, syscall.SIGINT, syscall.SIGTERM)
	<-c
	if debug {
		lgtm.PostMessage("TERMINATED")
	}
	time.Sleep(time.Second)

	return nil
}

func handleGithubEvents(lgtm *bot.LGTM, database *sql.Database) {
	for event := range github.IncomingEvents {
		switch ev := event.(type) {
		case github.AuthenticateEvent:
			err := database.CreateUser(ev.Id, ev.Token)

			if err != nil {
				Q(err)
				continue
			}

			lgtm.PostMessage(fmt.Sprintf(
				"<@%s> has Authenticated",
				ev.Id,
			))

		case github.PullRequestEvent:
			Q("PR Event receieved")
			switch ev.Action {
			case "opened":
				Q("OpenPREvent")
				ts, err := lgtm.PostMessage(ev.URL)
				if err != nil {
					Q(err)
					continue
				}
				_ = ts
				// TODO: Store pull request id and timestamp in database
			case "close":
				// TODO: pull out timestamp for pull request id and react to that message
				//lgtm.ReactPullRequest()
			}
		}
	}

}

func handleBotEvents(authurl string, lgtm *bot.LGTM, database *sql.Database) {
	ctx := context.Background()
	for event := range lgtm.IncomingEvents {
		switch ev := event.(type) {
		case bot.WatchRepoEvent:
			lgtm.PostMessage(fmt.Sprintf(
				"someday I will watch `%s` with owner `%s`, <@%s>",
				ev.Repo, ev.Owner, ev.User,
			))

			token, err := database.ReadUserAuth(ev.User)

			if err != nil {
				lgtm.PostMessage(fmt.Sprintf(
					"You need to authorize me <@%s>. %s?id=%s",
					ev.User, authurl, ev.User,
				))
				continue
			}

			hook, err := github.WatchRepo(ctx, token, ev.Owner, ev.Repo)

			if err != nil {
				// TODO: Handle already watching
				Q(err)
				continue
			}
			lgtm.PostMessage(fmt.Sprintf("And that day is today! Look at that! %s",
				*hook.URL,
			))
			Q(hook)
		}
	}
}
