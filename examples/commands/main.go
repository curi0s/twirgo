package main

import (
	"fmt"
	"os"

	"github.com/curi0s/twirgo"
	"github.com/sirupsen/logrus"
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

	t.Command([]string{"!bar"}, []string{"curi"}, twirgo.PermEveryone, cmdBar)
	t.Command([]string{"$foo"}, []string{"curi"}, twirgo.PermMods, cmdFoo)
	t.Command([]string{"$asd"}, []string{"curi"}, twirgo.PermBroadcaster, cmdAsd)

	t.Run(ch)
}
