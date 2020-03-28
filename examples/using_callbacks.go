package main

import (
	"fmt"
	"log"
	"os"

	"github.com/curi0s/twirgo"
)

func handleMessage(t *twirgo.Twitch, event twirgo.EventMessageReceived) {
	fmt.Println(event.Message.Content)
}

func handleUserJoin(t *twirgo.Twitch, event twirgo.EventUserJoined) {
	t.SendMessage(event.Channel.Name, "Welcome "+event.User.Username)
}

func main() {
	options := twirgo.Options{
		Username:       "curi0sde_bot",                       // the name of your bot account
		Token:          os.Getenv("TOKEN"),                   // provide your token in any way you like
		Channels:       []string{"curi0sde", "curi0sde_bot"}, // all channels will be joined at connect
		DefaultChannel: "curi0sde",                           // have a look an #L16
	}

	t := twirgo.NewTwirgo(options)

	ch, err := t.Connect()
	if err == twirgo.ErrInvalidToken {
		log.Fatal(err)
	}

	t.OnMessageReceived(handleMessage)
	t.OnUserJoined(handleUserJoin)

	t.Callbacks(ch)
}
