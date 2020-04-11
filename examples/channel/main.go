package main

import (
	"fmt"
	"os"

	"github.com/curi0s/twirgo"
	"github.com/sirupsen/logrus"
)

func handleEvents(t *twirgo.Twitch, ch chan interface{}) error {
	for event := range ch {
		switch ev := event.(type) {
		// case twirgo.EventConnected:
		// 	log.Println("Connected!")
		// 	t.SendMessage(t.Options().DefaultChannel, "HeyGuys")

		case twirgo.EventChannelJoined:
			fmt.Println("channel joined...")

		case twirgo.EventChannelParted:
			fmt.Println("channel parted...")

		case twirgo.EventWhisperReceived:
			fmt.Println("whisper received...")
			t.SendWhisper(ev.User.Username, "hi")
			t.PartChannel("curi0sde")

		case twirgo.EventMessageReceived:
			fmt.Println("message received...")
			fmt.Println(ev.Message.Content)
			t.SendMessage(ev.Channel.Name, "Thank you for your message :)")

		case twirgo.EventConnectionError:
			return ev.Err
		}
	}

	return nil
}

func main() {
	options := twirgo.Options{
		Username:       "curi",                        // the name of your bot account
		Token:          os.Getenv("TOKEN"),            // provide your token in any way you like
		Channels:       []string{"curi", "curi_bot_"}, // all channels will be joined at connect
		DefaultChannel: "curi",                        // have a look an #L16
		Log:            logrus.New(),
	}

	// options.Log.SetLevel(logrus.DebugLevel)

	t := twirgo.New(options)

	ch, err := t.Connect()
	if err == twirgo.ErrInvalidToken {
		options.Log.Fatal(err)
	}

	err = handleEvents(t, ch)
	if err != nil {
		options.Log.Fatal(err)
	}
}
