package main

import (
	"bufio"
	"fmt"
	"log"
	"math"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/slack-go/slack"
)

func main() {
	// If you set debugging, it will log all requests to the console
	// Useful when encountering issues
	// slack.New("YOUR_TOKEN_HERE", slack.OptionDebug(true))
	api := slack.New("xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx")
	userConvParams := slack.GetConversationsForUserParameters{
		UserID: "XXXXXXXX",
		Types:  []string{"public_channel"}}

	var err error
	var channels []slack.Channel
	if channels, _, err = api.GetConversationsForUser(&userConvParams); err != nil {
		log.Fatal(err)
	}

	var leaveAnyChannels bool

	for _, channel := range channels {
		var channelInfo *slack.Channel
		if channelInfo, err = api.GetConversationInfo(channel.ID, false); err != nil {
			log.Fatal(err)
		}

		var diffHours float64
		if diffHours, err = calculateTimeDiff(channelInfo); err != nil {
			log.Fatal(err)
		}

		if diffHours > 100 {
			var msgs []slack.Message
			msgs, err = getMsgs(channel, channelInfo, api)
			var simpleMsgsFound int
			if len(msgs) > 0 {
				for _, msg := range msgs {
					if len(msg.SubType) == 0 {
						simpleMsgsFound++
					}
				}
				if simpleMsgsFound > 0 {
					leaveAnyChannels = true
					// log.Infof("%f days", math.Round(diffHours/24))
					fmt.Printf("%s days\n", strconv.FormatFloat(math.Round(diffHours/24), 'f', -1, 64))
					fmt.Printf("%d msgs\n", simpleMsgsFound)
					fmt.Printf("channel: %s\n", channel.Name)
					fmt.Printf("------------------\n")
				}
			}
		}
	}

	if leaveAnyChannels {
		reader := bufio.NewReader(os.Stdin)
		fmt.Print("Leave channels? [Y/n] ")
		var text string
		var err error
		if text, err = reader.ReadString('\n'); err != nil {
			log.Fatal(err)
		}
		text = strings.TrimSuffix(text, "\n")
		if text == "" {
			text = "Y"
		}
		upperText := strings.ToUpper(text)
		if upperText != "Y" && upperText != "N" {
			log.Fatalf("Invalid choice: %s", text)
		}
	}
}

func calculateTimeDiff(channelInfo *slack.Channel) (diffHours float64, err error) {
	var lastRead int64
	if lastRead, err = strconv.ParseInt(strings.Split(channelInfo.LastRead, ".")[0], 10, 64); err != nil {
		return
	}
	t := time.Unix(lastRead, 0)
	now := time.Now()
	diff := now.Sub(t)
	diffHours = diff.Hours()
	return
}

// func getMsgs(channel slack.Channel, channelInfo *slack.Channel) (msgs []slack.Messages, err error) {
func getMsgs(channel slack.Channel, channelInfo *slack.Channel, api *slack.Client) (msgs []slack.Message, err error) {
	params := slack.GetConversationHistoryParameters{
		ChannelID: channel.ID,
		Oldest:    channelInfo.LastRead}
	var histResponse *slack.GetConversationHistoryResponse
	if histResponse, err = api.GetConversationHistory(&params); err != nil {
		return
	}
	msgs = histResponse.Messages
	return
}
