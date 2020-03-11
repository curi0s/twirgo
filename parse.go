package bot

import (
	"log"
	"strconv"
	"strings"
	"time"
)

type (
	User struct {
		Username    string
		DisplayName string
		Color       string
		IsPartner   bool
	}

	ChannelUser struct {
		User             *User
		SubscriberMonths int64

		IsMod         bool
		IsBroadcaster bool
		IsSubscriber  bool
		IsVIP         bool
	}

	Channel struct {
		Name          string
		Users         map[string]*User
		EmoteOnly     bool
		FollowersOnly bool
		// only messages with more than 9 chars allowed & must be unique
		R9k      bool
		Slow     bool
		SubsOnly bool
	}

	Message struct {
		Content string
		Id      string
	}

	// Usernotice holds all information of an USERNOTICE irc event
	// https://dev.twitch.tv/docs/irc/tags#usernotice-twitch-tags
	Usernotice struct {
		Type               string
		Message            string
		Channel            *Channel
		SubscriberMonths   int64
		Color              string
		DisplayName        string
		User               *User
		SystemMsg          string
		Timestamp          time.Time
		Raider             *User
		Promo              *Promo
		Recipient          *User
		Sender             *User
		ShouldShareStreak  bool
		StreakMonths       int64
		SubTier            SubTier
		SubPlanName        string
		RaidingViewerCount int64
		RitualName         string
		EarnedBadgeTier    int64
	}

	Promo struct {
		GiftTotal int64
		Name      string
	}

	SubTier interface{}

	SubTierPrime struct{}
	SubTierOne   struct{}
	SubTierTwo   struct{}
	SubTierThree struct{}
)

func (t *Twitch) convertTmiTs(tmiTs string) (time.Time, error) {
	unixTimestamp, err := strconv.ParseInt(tmiTs, 10, 64)
	if err != nil {
		return time.Now(), err
	}
	timestamp := time.Unix(unixTimestamp/1000, 0)
	return timestamp, nil
}

func (t *Twitch) toInt(s string) int64 {
	i, _ := strconv.ParseInt(s, 10, 64)
	return i
}

// parsePRIVMSG parses the PRIVMSG event of the IRC protocol
func (t *Twitch) parsePRIVMSG(line string) (time.Time, *Channel, ChannelUser, Message) {
	parts := strings.Split(strings.TrimLeft(line, "@"), " :")
	infos := strings.Split(parts[0], ";")
	message := Message{
		Content: strings.Join(parts[2:], " :"),
	}

	timestamp := time.Now()
	channel, _ := t.getChannel(strings.Split(parts[1], "#")[1])

	username := strings.Split(parts[1], "!")[0]
	user, _ := t.getUser(username)
	channelUser := ChannelUser{
		User: user,
	}

	for _, info := range infos {
		infoSplit := strings.Split(info, "=")
		k, v := infoSplit[0], infoSplit[1]
		switch k {
		case "tmi-sent-ts":
			timestamp, _ = t.convertTmiTs(v)

		case "badge-info":
			if strings.Contains(v, "subscriber") {
				channelUser.IsSubscriber = true
				// TODO: check if there is also the streak information available
				channelUser.SubscriberMonths = t.toInt(strings.Split(v, "/")[1])
			}

		case "badges":
			channelUser.IsBroadcaster = strings.Contains(v, "broadcaster")
			channelUser.User.IsPartner = strings.Contains(v, "partner")
			channelUser.IsMod = strings.Contains(v, "moderator") || channelUser.IsBroadcaster
			channelUser.IsVIP = strings.Contains(v, "vip")

		case "color":
			channelUser.User.Color = v

		case "display-name":
			channelUser.User.DisplayName = v

		case "mod":
			channelUser.IsMod, _ = strconv.ParseBool(v)

		case "id":
			message.Id = v
		}
	}

	t.addUserToChannel(username, channel.Name)

	return timestamp, channel, channelUser, message
}

// parseJOINPART parses the JOIN and PART events of the IRC protocol
func (t *Twitch) parseJOINPART(line string) {
	username := strings.TrimLeft(strings.Split(line, "!")[0], ":")
	channel := strings.Split(line, "#")[1]
	c := t.channels[channel]
	if strings.Contains(line, "JOIN") {
		t.addUserToChannel(username, channel)
	} else {
		delete(c.Users, username)
	}
}

// parseROOMSTATE parses the ROOMSTATE event of the irc protocol
func (t *Twitch) parseROOMSTATE(line string) (*Channel, error) {
	channel, err := t.getChannel(strings.Split(line, "#")[1])
	if err == ErrInvalidChannel {
		return nil, err
	}

	channel.EmoteOnly = strings.Contains(line, "emote-only=1")
	channel.FollowersOnly = !strings.Contains(line, "followers-only=-1")
	channel.R9k = strings.Contains(line, "r9k=1")
	channel.Slow = strings.Contains(line, "slow=1")
	channel.SubsOnly = strings.Contains(line, "subs-only=1")

	return channel, nil
}

func (t *Twitch) parseUSERNOTICE(line string) Usernotice {
	line = strings.TrimLeft(line, "@")
	// > @badge-info=<badge-info>;badges=<badges>;color=<color>;
	// display-name=<display-name>;emotes=<emotes>;id=<id-of-msg>;
	// login=<user>;mod=<mod>;msg-id=<msg-id>;room-id=<room-id>;
	// subscriber=<subscriber>;system-msg=<system-msg>;tmi-sent-ts=<timestamp>;
	// turbo=<turbo>;user-id=<user-id>;user-type=<user-type>
	//  :tmi.twitch.tv USERNOTICE #<channel> :<message>
	parts := strings.Split(line, " :")

	channel, _ := t.getChannel(strings.Split(parts[1], "#")[1])
	usernotice := Usernotice{
		Channel: channel,
		Promo:   &Promo{},
	}

	// Usernotice is with message
	switch {
	case len(parts) == 3:
		usernotice.Message = parts[2]
	case len(parts) > 3:
		usernotice.Message = strings.Join(parts[2:], " :")
	}

	isBroadcaster := false

	infos := strings.Split(parts[0], ";")
	for _, info := range infos {
		infoSplit := strings.Split(info, "=")
		k, v := infoSplit[0], infoSplit[1]
		switch k {
		case "badge-info":
			if strings.Contains(v, "subscriber") {
				usernotice.SubscriberMonths = t.toInt(strings.Split(v, "/")[1])
			}

		case "badges":
			// TODO: parse badges and place them in own structs
			isBroadcaster = strings.Contains(v, "broadcaster")

		case "color":
			usernotice.Color = v

		case "display-name":
			usernotice.DisplayName = v

		case "emotes":
			// TODO: parse emotes and place them in own structs

		case "login":
			// check if usernotice was triggered anonymously
			if channel.Name != v || isBroadcaster {
				usernotice.User, _ = t.getUser(v)
			}
		case "msg-id":
			usernotice.Type = v

		case "system-msg":
			usernotice.SystemMsg = v

		case "tmi-sent-ts":
			usernotice.Timestamp, _ = t.convertTmiTs(v)

		case "msg-param-cumulative-months":
			fallthrough
		case "msg-param-months":
			usernotice.SubscriberMonths = t.toInt(strings.Split(v, "/")[1])

		case "msg-param-displayName":
			usernotice.Raider, _ = t.getUser(v)
			usernotice.Raider.DisplayName = v

		case "msg-param-login":
			// we do not parse this field because its already handled in msg-param-displayName

		case "msg-param-promo-gift-total":
			usernotice.Promo.GiftTotal = t.toInt(v)

		case "msg-param-promo-name":
			usernotice.Promo.Name = v

		case "msg-param-recipient-display-name":
			usernotice.Recipient, _ = t.getUser(v)

		case "msg-param-recipient-id":
			// we actually do not need this field

		case "msg-param-recipient-user-name":
			// we do not parse this field because its already handled in msg-param-recipient-display-name

		case "msg-param-sender-login":
			// we do not parse this field because its already handled in msg-param-sender-name

		case "msg-param-sender-name":
			usernotice.Sender, _ = t.getUser(v)

		case "msg-param-should-share-streak":
			usernotice.ShouldShareStreak, _ = strconv.ParseBool(v)

		case "msg-param-streak-months":
			usernotice.StreakMonths = t.toInt(v)

		case "msg-param-sub-plan":
			switch v {
			case "Prime":
				usernotice.SubTier = SubTierPrime{}
			case "1000":
				usernotice.SubTier = SubTierOne{}
			case "2000":
				usernotice.SubTier = SubTierTwo{}
			case "3000":
				usernotice.SubTier = SubTierThree{}
			}

		case "msg-param-sub-plan-name":
			usernotice.SubPlanName = v

		case "msg-param-viewerCount":
			usernotice.RaidingViewerCount = t.toInt(v)

		case "msg-param-ritual-name":
			usernotice.RitualName = v

		case "msg-param-threshold":
			usernotice.EarnedBadgeTier = t.toInt(v)
		}
	}

	return usernotice
}

func (t *Twitch) parseCLEARCHAT(line string) (int64, *Channel, *User, error) {
	line = strings.TrimSpace(line)

	var parts []string
	var banDuration int64

	// check if ban-duration is delivered
	if strings.HasPrefix(line, "@") {
		parts = strings.Split(strings.TrimLeft(line, "@"), " :")

		banDuration = t.toInt(strings.Split(parts[0], "=")[1])
	} else {
		parts = strings.Split(line, " :")
	}

	user, err := t.getUser(parts[2])
	if err != nil {
		return 0, nil, nil, err
	}

	channel, err := t.getChannel(strings.Split(parts[1], "#")[1])
	if err != nil {
		return 0, nil, nil, err
	}

	return banDuration, channel, user, nil
}

func (t *Twitch) parseCLEARMSG(line string) (*User, *Channel, Message, error) {
	parts := strings.Split(strings.TrimLeft(line, "@"), " :")

	var message Message
	var user *User

	channel, err := t.getChannel(strings.Split(parts[1], "#")[1])
	if err != nil {
		return nil, nil, message, err
	}

	if len(parts) > 3 {
		message.Content = strings.Join(parts[2:], " :")
	} else {
		message.Content = parts[2]
	}

	items := strings.Split(parts[0], ";")
	for _, item := range items {
		itemSplit := strings.Split(item, "=")
		k, v := itemSplit[0], itemSplit[1]
		switch k {
		case "login":
			user, err = t.getUser(v)
			if err != nil {
				return nil, nil, message, err
			}
		case "target-msg-id":
			message.Id = v
		}
	}

	return user, channel, message, nil
}

// parseLine
func (t *Twitch) parseLine(line string) {
	switch {
	case strings.HasPrefix(line, "PING"):
		t.cEvents <- EventPinged{}

	case strings.HasPrefix(line, ":tmi.twitch.tv 001"):
		t.cEvents <- EventConnected{}

	case strings.HasPrefix(line, ":") && strings.Contains(line, "JOIN"):
		t.parseJOINPART(line)
		t.cEvents <- EventUserJoined{}

	case strings.HasPrefix(line, ":") && strings.Contains(line, "PART"):
		t.parseJOINPART(line)
		t.cEvents <- EventUserParted{}

	case strings.Contains(line, "PRIVMSG"):
		timestamp, channel, user, message := t.parsePRIVMSG(line)
		t.cEvents <- EventMessageReceived{Timestamp: timestamp, Channel: channel, Message: message, ChannelUser: user}

	case strings.Contains(line, "USERSTATE"):
		// we don't need to parse this event cause every message a user writes
		// all these informations are also provided
		t.cEvents <- EventUserstate{}

	case strings.Contains(line, "ROOMSTATE"):
		channel, err := t.parseROOMSTATE(line)
		if err == ErrInvalidChannel {
			// TODO: log
			return
		}

		t.cEvents <- EventRoomstate{Channel: channel}

	case strings.Contains(line, "USERNOTICE"):
		// TODO: split this up in different events like sub, resub, anonsub, subgift etc
		usernotice := t.parseUSERNOTICE(line)
		t.cEvents <- EventUsernotice{Usernotice: usernotice}

	case strings.Contains(line, "CLEARCHAT"):
		banDuration, channel, user, err := t.parseCLEARCHAT(line)
		if err != nil {
			// TODO: log
			return
		}
		t.cEvents <- EventClearchat{BanDuration: banDuration, Channel: channel, User: user}

	case strings.Contains(line, "CLEARMSG"):
		user, channel, message, err := t.parseCLEARMSG(line)
		if err != nil {
			// TODO: log
			return
		}
		t.cEvents <- EventClearmsg{User: user, Channel: channel, Message: message}

	default:
		log.Println("unhandled event", line)
	}
}
