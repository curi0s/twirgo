package bot

import (
	"bufio"
	"errors"
	"log"
	"net"
	"os"
	"strings"
)

// Options are used to configure the bot.
type Options struct {
	Username       string
	Token          string
	Channels       []string
	DefaultChannel string
}

// Twitch is a twitch IRC bot.
type Twitch struct {
	opts     Options
	channels map[string]*Channel
	users    map[string]*User

	conn    net.Conn
	cSend   chan string
	cEvents chan interface{}
}

var (
	ErrInvalidUsername = errors.New("invalid username")
	ErrInvalidChannel  = errors.New("invalid channel")
)

// NewTwitch returns a new Twitch bot.
func NewTwitch(options Options) *Twitch {
	return &Twitch{
		opts:     options,
		channels: make(map[string]*Channel),
		users:    make(map[string]*User),
		cSend:    make(chan string),
		cEvents:  make(chan interface{}),
	}
}

func (t *Twitch) Options() Options {
	return t.opts
}

// Connect establishes a connection to the IRC server,
// returning an event channel.
func (t *Twitch) Connect() chan interface{} {
	var err error
	t.conn, err = net.Dial("tcp", "irc.chat.twitch.tv:6667")
	if err != nil {
		log.Fatal(err)
		os.Exit(1)
	}

	go t.send()
	go t.receive()

	t.SendCommand("PASS " + t.opts.Token)
	t.SendCommand("NICK " + t.opts.Username)

	return t.eventTrigger()
}

// send writes the formatted message to the connection
func (t *Twitch) send() {
	for line := range t.cSend {
		t.conn.Write([]byte(line + "\r\n"))
	}
}

// SendMessage sends a message
func (t *Twitch) SendMessage(channel, message string) {
	t.SendCommand("PRIVMSG #" + strings.TrimLeft(channel, "#") + " : " + message)
}

// SendCommand sends a command
func (t *Twitch) SendCommand(message string) {
	t.cSend <- message
}

// JoinChannel joins the given channel
func (t *Twitch) JoinChannel(channel string) {
	t.getChannel(channel)
	t.SendCommand("JOIN #" + strings.TrimLeft(channel, "#"))
}

// PartChannel parts the given channel
func (t *Twitch) PartChannel(channel string) {
	t.SendCommand("PART #" + strings.TrimLeft(channel, "#"))
}

// getUser returns or creates the given user in the internal global list
func (t *Twitch) getUser(username string) (*User, error) {
	username = strings.TrimSpace(strings.ToLower(username))

	if username == "" {
		return nil, ErrInvalidUsername
	}

	u, ok := t.users[username]
	if !ok {
		u = &User{Username: username}
		t.users[username] = u
	}
	return u, nil
}

// getChannel returns or creates the given channel in the internal global list
func (t *Twitch) getChannel(channel string) (*Channel, error) {
	channel = strings.TrimSpace(channel)

	if channel == "" {
		return nil, ErrInvalidChannel
	}

	c, ok := t.channels[channel]
	if !ok {
		c = &Channel{
			Name:  channel,
			Users: make(map[string]*User),
		}
		t.channels[channel] = c
	}
	return c, nil
}

// addUserToChannel links an internal global user to an internal global channel
func (t *Twitch) addUserToChannel(username string, channel string) {
	c, err := t.getChannel(channel)
	if err == ErrInvalidChannel {
		return
	}

	u, err := t.getUser(username)
	if err == ErrInvalidUsername {
		return
	}

	c.Users[username] = u
}

// receive reads the buffer of the connection and parses all events
func (t *Twitch) receive() {
	buf := bufio.NewReader(t.conn)
	for {
		line, err := buf.ReadString('\n')
		if err != nil {
			t.cEvents <- EventConnectionError{Err: err}
			return
		}

		line = strings.TrimSpace(line)

		t.parseLine(line)
	}
}
