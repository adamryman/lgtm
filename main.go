package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/adamryman/gophersay/gopherart"
	"github.com/nlopes/slack"

	_ "github.com/joho/godotenv/autoload"
)

const prPartyID = "C3YJF4GP5"
const playgroundId = "C03LPQF0Y"

func main() {
	token := os.Getenv("SLACK_API_TOKEN")
	gopher, _ := gopherart.Asset("gopherart/gopher.ascii")

	api := slack.New(token)
	logger := log.New(os.Stdout, "slack-bot: ", log.Lshortfile|log.LstdFlags)
	slack.SetLogger(logger)
	api.SetDebug(true)

	rtm := api.NewRTM()
	go rtm.ManageConnection()
	go func() {
		c := make(chan os.Signal, 1)
		signal.Notify(c, syscall.SIGINT, syscall.SIGTERM)
		<-c
		rtm.SendMessage(rtm.NewOutgoingMessage("TERMINATED", playgroundId))
		time.Sleep(time.Second)
		os.Exit(0)
	}()

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
			// Replace #general with your Channel ID
			rtm.SendMessage(rtm.NewOutgoingMessage("Rebooting....Done "+time.Now().String(), playgroundId))

		case *slack.MessageEvent:
			fmt.Printf("Message: %v\n", ev)
			txt := ev.Text

			if strings.HasPrefix(txt, "gophersay") {
				txt = strings.TrimPrefix(txt, "gophersay ")
				message := fmt.Sprintf("```\n%s\n%s\n%s\n```\n",
					" ------------------------",
					txt,
					gopher,
				)
				rtm.SendMessage(rtm.NewOutgoingMessage(message, playgroundId))
			}

			if strings.Contains(txt, "<@U456ZSLSJ>") {
				newtxt := strings.Replace(txt, "<@U456ZSLSJ>", "", -1)
				rtm.SendMessage(rtm.NewOutgoingMessage(Reverse(newtxt), playgroundId))
			}

			//if ev.User == "U03L3L1NC" {
			//count := strings.Count(txt, " ")
			//rtm.SendMessage(rtm.NewOutgoingMessage(fmt.Sprintf("John said %d words", count), playgroundId))
			//}

		case *slack.PresenceChangeEvent:
			fmt.Printf("Presence Change: %v\n", ev)

		case *slack.LatencyReport:
			fmt.Printf("Current latency: %v\n", ev.Value)

		case *slack.RTMError:
			fmt.Printf("Error: %s\n", ev.Error())

		case *slack.InvalidAuthEvent:
			fmt.Printf("Invalid credentials")
			return

		default:

			// Ignore other events..
			// fmt.Printf("Unexpected: %v\n", msg.Data)
		}
	}
}

func Reverse(s string) string {
	runes := []rune(s)
	for i, j := 0, len(runes)-1; i < j; i, j = i+1, j-1 {
		runes[i], runes[j] = runes[j], runes[i]
	}
	return string(runes)
}
