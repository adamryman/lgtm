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

	_ "github.com/joho/godotenv/autoload"
)

func Start() error {
	token := os.Getenv("SLACK_API_TOKEN")

	logger := log.New(os.Stdout, "slack-bot: ", log.Lshortfile|log.LstdFlags)
	ctx := context.Background()

	lgtm, err := bot.Start(ctx, token, bot.SetLogger(logger))
	if err != nil {
		return err
	}
	lgtm.PostMessage("Rebooting....Done " + time.Now().String())

	go func() {
		c := make(chan os.Signal, 1)
		signal.Notify(c, syscall.SIGINT, syscall.SIGTERM)
		<-c
		lgtm.PostMessage("TERMINATED")
		time.Sleep(time.Second)
	}()

	for event := range lgtm.IncomingEvents {
		switch ev := event.(type) {
		case bot.WatchRepoEvent:
			lgtm.PostMessage(fmt.Sprintf("someday I will watch `%s` with owner `%s`", ev.Repo, ev.Owner))
		}

	}

	return nil
}
