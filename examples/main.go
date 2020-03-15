package main

import (
	"log"
	"os"

	"github.com/curi0s/twirgo"
)

func handleEvents(t *twirgo.Twitch, ch chan interface{}) error {
	for event := range ch {
		switch ev := event.(type) {
		// case twirgo.EventConnected:
		// 	log.Println("Connected!")
		// 	t.SendMessage(t.Options().DefaultChannel, "HeyGuys")

		// case twirgo.EventMessageReceived:
		// 	fmt.Printf("%+v\n", ev)
		// 	fmt.Printf("%+v\n", ev.ChannelUser)
		// 	fmt.Printf("%+v\n", ev.ChannelUser.User)
		// 	fmt.Println(ev.Channel.Name, ev.Message)

		case twirgo.EventConnectionError:
			return ev.Err
		}
	}

	return nil
}

func main() {
	token := os.Getenv("TOKEN")
	if token == "" {
		log.Fatal("Empty TOKEN")
	}

	t := twirgo.NewTwirgo(twirgo.Options{
		Username:       "curi0sde_bot",
		Token:          token,
		Channels:       []string{"curi0sde", "rocketleague", "AdmiralBahroo", "summit1g"},
		DefaultChannel: "curi0sde",
	})

	ch := t.Connect()
	err := handleEvents(t, ch)
	if err != nil {
		log.Fatal(err)
	}
}
