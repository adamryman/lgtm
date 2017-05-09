package lgtm

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/StudentRND/lgtm/bot"
	"github.com/StudentRND/lgtm/github"

	_ "github.com/joho/godotenv/autoload"
	. "github.com/y0ssar1an/q"
)

func Start() error {
	slackToken := os.Getenv("SLACK_API_TOKEN")
	githubToken := os.Getenv("GITHUB_TOKEN")
	Q(githubToken)

	Q(slackToken)
	logger := log.New(os.Stdout, "slack-bot: ", log.Lshortfile|log.LstdFlags)
	ctx := context.Background()

	lgtm, err := bot.Start(ctx, slackToken, bot.SetLogger(logger))
	if err != nil {
		return err
	}
	lgtm.PostMessage("Rebooting....Done " + time.Now().String())

	// handle chat events
	go func() {
		for event := range lgtm.IncomingEvents {
			switch ev := event.(type) {
			case bot.WatchRepoEvent:
				lgtm.PostMessage(fmt.Sprintf("someday I will watch `%s` with owner `%s`, <@%s>", ev.Repo, ev.Owner, ev.User))
				hook, err := github.WatchRepo(ctx, githubToken, ev.Owner, ev.Repo)
				if err != nil {
					fmt.Println(err)
					Q(err)
					continue
				}
				lgtm.PostMessage(fmt.Sprintf("And that day is today! Look at that! %s", *hook.URL))
				Q(hook)
			}
		}
	}()

	// handle github events
	go func() {

	}()

	// Wait for termination
	c := make(chan os.Signal, 1)
	signal.Notify(c, syscall.SIGINT, syscall.SIGTERM)
	<-c
	lgtm.PostMessage("TERMINATED")
	time.Sleep(time.Second)

	return nil
}
