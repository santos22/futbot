package guardian

import (
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/guardian/gocapiclient"
	"github.com/guardian/gocapiclient/queries"
	"github.com/nlopes/slack"
)

// Guardian struct to hold api client
type Guardian struct {
	GuardianClient *gocapiclient.GuardianContentClient
}

// map to keep track of event IDs already sent
var (
	events = make(map[string]bool)
)

// SendEvent - send football events to Slack channel
func (c *Guardian) SendEvent(rtm *slack.RTM, msg *slack.MessageEvent, prefix string) {
	gameURL := msg.Text
	gameURL = strings.TrimPrefix(gameURL, prefix)
	gameURL = strings.TrimSpace(gameURL)
	gameURL = strings.ToLower(gameURL)

	for {
		c.getEvents(gameURL, rtm, msg, c.GuardianClient)
		time.Sleep(30 * time.Second) // request updates every 30 seconds
	}
}

// get live updates from www.theguardian.com/football/live/
func (c *Guardian) getEvents(gameURL string, rtm *slack.RTM, msg *slack.MessageEvent, client *gocapiclient.GuardianContentClient) {
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
			// send event to channel
			events[event.ID] = true
			rtm.SendMessage(rtm.NewOutgoingMessage(event.BodyTextSummary, msg.Channel))
		}
	}
}
