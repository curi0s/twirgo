package twirgo

import (
	"bufio"
	"crypto/tls"
	"errors"
	"net"
	"strings"

	"github.com/sirupsen/logrus"
)

// Options are used to configure the bot.
type Options struct {
	Username       string
	Token          string
	Channels       []string
	DefaultChannel string
	Log            *logrus.Logger
	Unsecure       bool
}

// Twitch is a Twitch IRC bot.
type Twitch struct {
	opts     Options
	channels map[string]*Channel
	users    map[string]*User

	conn    net.Conn
	cSend   chan string
	cEvents chan interface{}

	callbacks callbacks

	commands map[string]*commandDefinition

	log *logrus.Logger
}

var (
	ErrInvalidUsername = errors.New("invalid username provided")
	ErrInvalidChannel  = errors.New("invalid channel provided")
	ErrInvalidToken    = errors.New("invalid token provided")
)

// New returns a new Twitch bot.
func New(options Options) *Twitch {

	// normalize options
	options.Username = strings.ToLower(strings.TrimSpace(options.Username))
	options.DefaultChannel = strings.ToLower(strings.TrimSpace(options.DefaultChannel))

	return &Twitch{
		opts:     options,
		channels: make(map[string]*Channel),
		users:    make(map[string]*User),
		cSend:    make(chan string),
		cEvents:  make(chan interface{}),
		commands: make(map[string]*commandDefinition),
		log:      options.Log,
	}
}

// Options returns the configured options for the this bot
func (t *Twitch) Options() Options {
	return t.opts
}

// Connect establishes a connection to the IRC server,
// returning an event channel.
func (t *Twitch) Connect() (chan interface{}, error) {
	if strings.TrimSpace(t.opts.Token) == "" {
		t.log.Fatal("Invalid token provided")
		return nil, ErrInvalidToken
	}

	var err error
	if t.opts.Unsecure == true {
		t.conn, err = net.Dial("tcp", "irc.chat.twitch.tv:6667")
	} else {
		t.conn, err = tls.Dial("tcp", "irc.chat.twitch.tv:6697", &tls.Config{})
	}
	if err != nil {
		t.log.Fatal("Could not connect ", err)
		return nil, err
	}

	t.log.Info("Connected to Twitch chat server")

	go t.send()
	go t.receive()

	t.SendCommand("PASS " + t.opts.Token)
	t.SendCommand("NICK " + t.opts.Username)

	return t.cEvents, nil
}

// Run loops and handles all callback functions
// Never use this method if you handle the channel (given from Connect()) yourself
func (t *Twitch) Run(chan interface{}) {
	for event := range t.cEvents {
		t.callCallbacks(event)
	}
}

// receive reads the buffer of the connection and parses all events
func (t *Twitch) receive() {
	buf := bufio.NewReader(t.conn)
	for {
		line, err := buf.ReadString('\n')
		if err != nil {
			t.log.Fatal("Connection error", err)
			t.cEvents <- EventConnectionError{Err: err}
			return
		}

		line = strings.TrimSpace(line)

		t.log.Debug("> " + line)

		t.parseLine(line)
	}
}

// send writes the formatted message to the connection
func (t *Twitch) send() {
	for line := range t.cSend {
		t.log.Debug("< " + line)
		t.conn.Write([]byte(line + "\r\n"))
	}
}

// SendMessage sends a message to the given channel
func (t *Twitch) SendMessage(channel, message string) {
	c, err := t.getChannel(channel)
	if err != nil {
		return
	}

	t.SendCommand("PRIVMSG #" + c.Name + " :" + message)
}

// SendWhisper sends a whisper to the given user
func (t *Twitch) SendWhisper(username, message string) {
	t.SendCommand("PRIVMSG #" + username + " :/w " + username + " " + message)
}

// SendCommand sends a command
func (t *Twitch) SendCommand(message string) {
	t.cSend <- message
}

// JoinChannel joins the given channel
func (t *Twitch) JoinChannel(channel string) {
	c, err := t.getChannel(channel)
	if err != nil {
		return
	}

	t.SendCommand("JOIN #" + c.Name)
}

// PartChannel parts the given channel
func (t *Twitch) PartChannel(channel string) {
	c, err := t.getChannel(channel)
	if err != nil {
		return
	}

	t.SendCommand("PART #" + c.Name)
	// has to be in a goroutine otherwise it would block > see examples/main.go#L27
	go func() {
		t.cEvents <- EventChannelParted{}
	}()
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
func (t *Twitch) addUserToChannel(user *User, channel *Channel) {
	channel.Users[user.Username] = user
}

func (t *Twitch) inSlice(s []string, w string) bool {
	for _, e := range s {
		if strings.ToLower(e) == strings.ToLower(w) {
			return true
		}
	}
	return false
}
