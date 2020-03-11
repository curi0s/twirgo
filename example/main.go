package main

import (
	"log"
	"os"

	bot "github.com/curi0s/twirgo"
)

func handleEvents(t *bot.Twitch, ch chan interface{}) error {
	for event := range ch {
		switch ev := event.(type) {
		case bot.EventConnected:
			log.Println("Connected!")
			t.SendMessage(t.Options().DefaultChannel, "HeyGuys")

		// case bot.EventMessageReceived:
		// 	fmt.Printf("%+v\n", ev)
		// 	fmt.Printf("%+v\n", ev.ChannelUser)
		// 	fmt.Printf("%+v\n", ev.ChannelUser.User)
		// 	fmt.Println(ev.Channel.Name, ev.Message)

		case bot.ConnectionError:
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

	t := bot.NewTwitch(bot.Options{
		Username:       "curi0sde_bot",
		Token:          token,
		Channels:       []string{"curi0sde"},
		DefaultChannel: "curi0sde",
	})

	ch := t.Connect()
	err := handleEvents(t, ch)
	if err != nil {
		log.Fatal(err)
	}
}
