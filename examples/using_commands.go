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

func cmdBar(t *twirgo.Twitch, c twirgo.Command) {
	fmt.Println("command " + c.Command + " invoked")
	fmt.Printf("%+v\n", c.Mentions)
}

func cmdFoo(t *twirgo.Twitch, c twirgo.Command) {
	fmt.Println("command " + c.Command + " invoked")
	fmt.Printf("%+v\n", c.Mentions)
}

func cmdAsd(t *twirgo.Twitch, c twirgo.Command) {
	fmt.Println("command " + c.Command + " invoked")
	fmt.Printf("%+v\n", c.Mentions)
}

func main() {
	options := twirgo.Options{
		Username:       "curi0sde_bot",                       // the name of your bot account
		Token:          os.Getenv("TOKEN"),                   // provide your token in any way you like
		Channels:       []string{"curi0sde", "curi0sde_bot"}, // all channels will be joined at connect
		DefaultChannel: "curi0sde",
	}

	t := twirgo.NewTwirgo(options)

	ch, err := t.Connect()
	if err == twirgo.ErrInvalidToken {
		log.Fatal(err)
	}

	t.OnMessageReceived(handleMessage)

	t.Command([]string{"!bar"}, []string{"curi0sde"}, twirgo.PermEveryone, cmdBar)
	t.Command([]string{"$foo"}, []string{"curi0sde"}, twirgo.PermMods, cmdFoo)
	t.Command([]string{"$asd"}, []string{"curi0sde"}, twirgo.PermBroadcaster, cmdAsd)

	t.Run(ch)
}
