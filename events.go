package twirgo

import (
	"time"
)

type (
	EventConnectionError struct {
		Err error
	}

	EventConnected     struct{}
	EventPinged        struct{}
	EventJoinedChannel struct{}
	EventPartedChannel struct{}

	EventUserJoined struct {
		Channel *Channel
		User    *User
	}
	EventUserParted struct {
		Channel *Channel
		User    *User
	}
	EventMessageReceived struct {
		Timestamp   time.Time
		Channel     *Channel
		ChannelUser ChannelUser
		Message     Message
	}

	EventUserstate struct {
		Channel *Channel
		User    *User
	}
	EventRoomstate struct {
		Channel *Channel
	}
	EventSub                 struct{}
	EventResub               struct{}
	EventSubgift             struct{}
	EventAnonsubgift         struct{}
	EventSubmysterygift      struct{}
	EventGiftpaidupgrade     struct{}
	EventRewardgift          struct{}
	EventAnongiftpaidupgrade struct{}
	EventRaid                struct{}
	EventUnraid              struct{}
	EventRitual              struct{}
	EventBitsbadgetier       struct{}

	EventClearchat struct {
		Timestamp   time.Time
		BanDuration int64
		Channel     *Channel
		User        *User
	}
	EventClearmsg struct {
		Timestamp time.Time
		User      *User
		Channel   *Channel
		Message   Message
	}

	EventNotice struct {
		MsgId   string
		Channel *Channel
	}

	EventStartHosting struct {
		FromChannel *Channel
		ToChannel   *Channel
		Viewers     int64 // might be 0 if not provided
	}

	EventStopHosting struct {
		FromChannel *Channel
		Viewers     int64 // might be 0 if not provided
	}
)
