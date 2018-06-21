package main

import (
	"fmt"
	"strings"

	"github.com/guardian/gocapiclient"
	"github.com/nlopes/slack"
	"github.com/santos22/slack-wc/guardian"
	"github.com/santos22/slack-wc/utils"
)

var (
	slackToken    = utils.GetEnv("SLACK_TOKEN")
	guardianToken = utils.GetEnv("GUARDIAN_TOKEN")
)

func main() {
	client := &guardian.Guardian{
		GuardianClient: gocapiclient.NewGuardianContentClient("https://content.guardianapis.com/", guardianToken),
	}

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
					client.SendEvent(rtm, ev, prefix)
				}

			case *slack.RTMError:
				fmt.Printf("Error: %s\n", ev.Error())

			case *slack.InvalidAuthEvent:
				fmt.Printf("Invalid credentials")
				break Loop

			default:
				// take no action
			}
		}
	}
}
