# slack-wc
I sometimes find myself at work during an afternoon Champions League game (and in 2018, a World Cup game). Instead of switching back and forth between theguardian.com/football/live and Slack, I thought it would be neat to bring all of the live updates to Slack. @wc-bot is a Slack bot that posts live game updates to your Slack channel of choice. Want to share your thoughts during/after the game? Having a channel with live updates is also the perfect place for those watching the game to discuss what is going on.

## Tokens
You can request a developer key from the Guardian at the following link:

* https://open-platform.theguardian.com/access/

After receiving both tokens, set them as environment variables:

```
export SLACK_TOKEN="SLACK_TOKEN"

export GUARDIAN_TOKEN="GUARDIAN_TOKEN"
```

## Running the program
```
go run main.go
```

## Using @wc-bot
Let's say you can't watch the highly anticpated World Cup game between Portugal and Spain. You can keep switching between Slack and https://www.theguardian.com/football/live/2018/jun/15/portugal-v-spain-world-cup-2018-live...

**OR** you can get live updates in Slack by calling upon @wc-bot!

In your Slack channel, enter the following: @wc-bot football/live/2018/jun/15/portugal-v-spain-world-cup-2018-live

## Improvements
- [ ] Include any media (e.g. photos) as attachments in Slack messages