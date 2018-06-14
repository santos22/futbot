package main

import (
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"github.com/guardian/gocapiclient"
	"github.com/guardian/gocapiclient/queries"
	"github.com/nlopes/slack"
)

var (
	slackToken    = getenv("SLACK_TOKEN")
	guardianToken = getenv("GUARDIAN_TOKEN")
	events        = make(map[string]bool)
)

func main() {
	api := slack.New(slackToken)
	api.SetDebug(true)

	rtm := api.NewRTM()
	go rtm.ManageConnection()

Loop:
	for {
		select {
		case msg := <-rtm.IncomingEvents:
			fmt.Print("Event Received: ")
			switch ev := msg.Data.(type) {
			case *slack.ConnectedEvent:
				fmt.Println("Connection counter:", ev.ConnectionCount)

			case *slack.MessageEvent:
				fmt.Printf("Message: %v\n", ev)
				info := rtm.GetInfo()
				prefix := fmt.Sprintf("<@%s> ", info.User.ID)

				if ev.User != info.User.ID && strings.HasPrefix(ev.Text, prefix) {
					sendEvent(rtm, ev, prefix)
				}

			case *slack.RTMError:
				fmt.Printf("Error: %s\n", ev.Error())

			case *slack.InvalidAuthEvent:
				fmt.Printf("Invalid credentials")
				break Loop

			default:
				//Take no action
			}
		}
	}
}

// Retrieve access tokens stored as environment variable
func getenv(name string) string {
	v := os.Getenv(name)
	if v == "" {
		panic("missing required environment variable " + name)
	}
	return v
}

// Send football events to Slack channel
func sendEvent(rtm *slack.RTM, msg *slack.MessageEvent, prefix string) {
	text := msg.Text
	text = strings.TrimPrefix(text, prefix)
	text = strings.TrimSpace(text)
	text = strings.ToLower(text)

	client := gocapiclient.NewGuardianContentClient("https://content.guardianapis.com/", guardianToken)

	for {
		getEvents(text, client, rtm, msg)
		time.Sleep(2 * time.Minute)
	}
}

func getEvents(gameURL string, client *gocapiclient.GuardianContentClient, rtm *slack.RTM, msg *slack.MessageEvent) {
	itemQuery := queries.NewItemQuery(gameURL)

	showParam := queries.StringParam{"show-blocks", "body:latest"}
	params := []queries.Param{&showParam}

	itemQuery.Params = params

	err := client.GetResponse(itemQuery)

	if err != nil {
		log.Fatal(err)
	}

	blocks := itemQuery.Response.Content.Blocks.RequestedBodyBlocks["body:latest"]
	fmt.Println(blocks)

	for _, event := range blocks {
		if events[event.ID] {
			continue
		} else {
			// add event ID to map and send event to channel
			events[event.ID] = true
			rtm.SendMessage(rtm.NewOutgoingMessage(event.BodyTextSummary, msg.Channel))
		}
	}
}
