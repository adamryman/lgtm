package bot

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/nlopes/slack"
	"github.com/pkg/errors"

	_ "github.com/joho/godotenv/autoload"
	. "github.com/y0ssar1an/q"
)

var debug bool

func init() {
	Q(os.Getenv("DEBUG"))
	if len(os.Getenv("DEBUG")) > 0 {
		debug = true
	}
}

type LGTM struct {
	id             string
	channel        string
	api            *slack.Client
	emoji          string
	logger         *log.Logger
	IncomingEvents chan interface{}
}

type WatchRepoEvent struct {
	User  string
	Owner string
	Repo  string
}

// LGTMOption is a function that modifies the client config
type LGTMOption func(*LGTM) error

func SetLogger(logger *log.Logger) LGTMOption {
	return func(o *LGTM) error {
		o.logger = logger
		return nil
	}
}

func Start(token, id, channel string, options ...LGTMOption) (*LGTM, error) {
	lgtm := LGTM{
		IncomingEvents: make(chan interface{}),
		id:             id,
		channel:        channel,
		emoji:          "white_check_mark",
	}

	for _, f := range options {
		err := f(&lgtm)
		if err != nil {
			return nil, errors.Wrap(err, "cannot apply option")
		}
	}

	slack.SetLogger(lgtm.logger)
	lgtm.api = slack.New(token)
	go lgtm.start(token)

	return &lgtm, nil
}

func (lgtm LGTM) PostMessage(msg string) (timestamp string, err error) {
	pmp := slack.NewPostMessageParameters()
	pmp.AsUser = true
	pmp.LinkNames = 1
	pmp.EscapeText = false
	pmp.UnfurlLinks = true
	pmp.UnfurlMedia = true
	_, timestamp, err = lgtm.api.PostMessage(lgtm.channel, msg, pmp)
	return
}

func (lgtm LGTM) ReactPullRequest(timestamp string) error {
	pmp := slack.NewPostMessageParameters()
	pmp.AsUser = true
	itemref := slack.NewRefToMessage(lgtm.channel, timestamp)
	return lgtm.api.AddReaction(lgtm.emoji, itemref)
}

func (lgtm *LGTM) start(token string) {

	lgtm.api.SetDebug(debug)

	rtm := lgtm.api.NewRTM()
	go rtm.ManageConnection()

	for msg := range rtm.IncomingEvents {
		fmt.Print("Event Received: ")
		switch ev := msg.Data.(type) {
		case *slack.HelloEvent:
			// Ignore hello
		case *slack.ChannelJoinedEvent:
			fmt.Println(ev.Channel.ID)

		case *slack.ConnectedEvent:
			fmt.Println("Infos:", ev.Info)
			fmt.Println("Connection counter:", ev.ConnectionCount)

		case *slack.MessageEvent:
			text := ev.Text
			// Not right channel
			if ev.Channel != lgtm.channel {
				continue
			}
			// Does not mention @lgtm
			if !strings.Contains(text, lgtm.id) {
				continue
			}

			// TODO: Make this more robust
			// " watch repo adamryman/lgtm "
			watchRequest := strings.Split(text, " watch repo ")
			switch {
			case len(watchRequest) > 1:

				Q(ev.User)
				Q(watchRequest)

				// should be "adamryman/lgtm blah blah"
				repoPart := watchRequest[len(watchRequest)-1]
				go lgtm.handleWatchRequest(ev, repoPart)
			}

			fmt.Printf("Message: %v\n", ev)
		}
	}
}

func (lgtm *LGTM) handleWatchRequest(msg *slack.MessageEvent, repoPart string) {

	Q(repoPart)

	// Move to the first "word" space seperated
	scanner := bufio.NewScanner(strings.NewReader(repoPart))
	scanner.Split(bufio.ScanWords)
	if !scanner.Scan() {
		return
	}

	repoText := scanner.Text()
	// Break it by slash to get owner/repo
	ownerRepo := strings.Split(repoText, "/")

	Q(repoText, ownerRepo)

	if len(ownerRepo) < 2 {
		return
	}
	// TODO: Handle punctuation. "watch repo adamryman/lgtm."
	owner, repo := ownerRepo[0], ownerRepo[1]

	Q(owner, repo)

	lgtm.IncomingEvents <- WatchRepoEvent{User: msg.User, Owner: owner, Repo: repo}
}
