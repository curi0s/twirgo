package main

import (
	"fmt"
	"os"

	"github.com/curiTTV/twirgo"
	"github.com/sirupsen/logrus"
)

func handleMessage(t *twirgo.Twitch, event twirgo.EventMessageReceived) {
	fmt.Println(event.Message.Content)
	fmt.Println(event.ChannelUser.Badges)
}

// func handleUserJoin(t *twirgo.Twitch, event twirgo.EventUserJoined) {
// 	t.SendMessage(event.Channel.Name, "Welcome "+event.User.Username)
// }

func main() {
	options := twirgo.Options{
		Username:       "curi_bot_",                   // the name of your bot account
		Token:          os.Getenv("TOKEN"),            // provide your token in any way you like
		Channels:       []string{"curi", "curi_bot_"}, // all channels will be joined at connect
		DefaultChannel: "curi",
		Log:            logrus.New(),
	}

	// options.Log.SetLevel(logrus.DebugLevel)

	t := twirgo.New(options)

	ch, err := t.Connect()
	if err == twirgo.ErrInvalidToken {
		options.Log.Fatal(err)
	}

	t.OnMessageReceived(handleMessage)
	// t.OnUserJoined(handleUserJoin)

	t.Run(ch)
}
