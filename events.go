package bot

import "time"

type (
	EventConnected struct{}
	EventPinged    struct{}

	EventUserJoined      struct{}
	EventUserParted      struct{}
	EventMessageReceived struct {
		Timestamp   time.Time
		Channel     *Channel
		ChannelUser ChannelUser
		Message     Message
	}

	EventUserstate  struct{}
	EventUsernotice struct {
		Usernotice Usernotice
	}
	EventRoomstate struct {
		Channel *Channel
	}

	EventClearchat struct {
		BanDuration int64
		Channel     *Channel
		User        *User
	}
	EventClearmsg struct {
		User    *User
		Channel *Channel
		Message Message
	}

	EventConnectionError struct {
		Err error
	}
)

func (t *Twitch) eventTrigger() chan interface{} {
	ch := make(chan interface{})

	go func() {
		for event := range t.cEvents {
			ch <- event

			switch event.(type) {
			case EventConnected:
				t.SendCommand("CAP REQ :twitch.tv/membership")
				t.SendCommand("CAP REQ :twitch.tv/tags")
				t.SendCommand("CAP REQ :twitch.tv/commands")
				for _, channel := range t.opts.Channels {
					t.JoinChannel(channel)
				}

			case EventPinged:
				t.SendCommand("PONG :tmi.twitch.tv")
			}
		}
	}()

	return ch
}
