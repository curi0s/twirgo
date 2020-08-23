package twirgo

import (
	"regexp"
	"strconv"
	"strings"
	"time"
	"unicode"
)

type parsedLine struct {
	tags    Tags
	user    *User
	t       string
	channel *Channel
	message Message
}

// convertTmiTs converts a TMI-TS which is a unix timestamp in milliseconds to a time.Time
func (t *Twitch) convertTmiTs(tmiTs string) (time.Time, error) {
	unixTimestamp, err := strconv.ParseInt(tmiTs, 10, 64)
	if err != nil {
		return time.Now(), err
	}
	timestamp := time.Unix(unixTimestamp/1000, 0)
	return timestamp, nil
}

// toInt converts a string into an int64
func (t *Twitch) toInt(s string) int64 {
	i, _ := strconv.ParseInt(s, 10, 64)
	return i
}

// parseJOINPART parses the JOIN and PART events of the IRC protocol
func (t *Twitch) parseJOINPART(line string) {
	user, _ := t.getUser(strings.TrimLeft(strings.Split(line, "!")[0], ":"))
	channel, _ := t.getChannel(strings.Split(line, "#")[1])
	if strings.Contains(line, "JOIN") {
		t.addUserToChannel(user, channel)
	} else {
		delete(channel.Users, user.Username)
	}
}

func (t *Twitch) parseTags(tags string) Tags {
	tagMap := make(Tags)

	for _, tag := range strings.Split(tags, ";") {
		s := strings.Split(tag, "=")
		tagMap[s[0]] = s[1]
	}

	return tagMap
}

func (t *Twitch) parseBadges(badgesString string) Badges {
	b := make(Badges)

	if badgesString != "" {
		if strings.Contains(badgesString, ",") {
			for _, badge := range strings.Split(badgesString, ",") {
				badgeInfo := strings.Split(badge, "/")
				b[badgeInfo[0]] = badgeInfo[1]
			}
		} else {
			badgeInfo := strings.Split(badgesString, "/")
			b[badgeInfo[0]] = badgeInfo[1]
		}
	}

	return b
}

func (t *Twitch) buildChannelUser(parsedLine *parsedLine) ChannelUser {
	channelUser := ChannelUser{
		User:   parsedLine.user,
		Badges: t.parseBadges(parsedLine.tags["badges"]),
	}

	badgeInfo := parsedLine.tags["badge-info"]
	if badgeInfo != "" {
		// there is currently only the subscriber months in it
		parts := strings.Split(badgeInfo, "/")
		channelUser.IsSubscriber = true
		channelUser.SubscriberMonths = t.toInt(parts[1])
	}

	channelUser.IsMod, _ = strconv.ParseBool(parsedLine.tags["mod"])
	channelUser.IsBroadcaster = channelUser.Badges["broadcaster"] == "1"
	// a broadcaster is also a mod
	if channelUser.IsBroadcaster {
		channelUser.IsMod = true
	}
	channelUser.IsVIP = channelUser.Badges["vip"] == "1"
	channelUser.User.IsPartner = channelUser.Badges["parner"] == "1"
	channelUser.User.DisplayName = parsedLine.tags["display-name"]
	channelUser.User.Color = parsedLine.tags["color"]

	return channelUser
}

func (t *Twitch) parseUSERNOTICE(parsedLine *parsedLine) {
	var subTier SubTier

	switch parsedLine.tags["msg-param-sub-plan"] {
	case "Prime":
		subTier = SubTierPrime{Name: parsedLine.tags["msg-param-sub-plan-name"]}

	case "1000":
		subTier = SubTierOne{Name: parsedLine.tags["msg-param-sub-plan-name"]}

	case "2000":
		subTier = SubTierTwo{Name: parsedLine.tags["msg-param-sub-plan-name"]}

	case "3000":
		subTier = SubTierThree{Name: parsedLine.tags["msg-param-sub-plan-name"]}

	}

	switch parsedLine.tags["msg-id"] {
	case "sub":
		fallthrough

	case "resub":
		shareStreak, _ := strconv.ParseBool(parsedLine.tags["msg-param-should-share-streak"])

		t.cEvents <- Sub{
			Months:       t.toInt(parsedLine.tags["msg-param-cumulative-months"]),
			ShareStreak:  shareStreak,
			StreakMonths: t.toInt(parsedLine.tags["msg-param-streak-months"]),
			SubTier:      &subTier,
		}

	case "subgift":
		fallthrough

	case "anonsubgift":
		user, _ := t.getUser(parsedLine.tags["msg-param-recipient-user-name"])
		user.DisplayName = parsedLine.tags["msg-param-recipient-display-name"]
		user.ID = t.toInt(parsedLine.tags["msg-param-recipient-id"])

		t.cEvents <- Subgift{
			Months:  t.toInt(parsedLine.tags["msg-param-cumulative-months"]),
			User:    user,
			SubTier: &subTier,
		}

	case "submysterygift":
		t.cEvents <- Submysterygift{}

	case "giftpaidupgrade":
		user, _ := t.getUser(parsedLine.tags["msg-param-sender-login"])
		user.DisplayName = parsedLine.tags["msg-param-sender-name"]

		t.cEvents <- Giftpaidupgrade{
			Gifts: t.toInt(parsedLine.tags["msg-param-months"]),
			Name:  parsedLine.tags["msg-param-promo-name"],
			User:  user,
		}

	case "rewardgift":
		t.cEvents <- Rewardgift{}

	case "anongiftpaidupgrade":
		t.cEvents <- Anongiftpaidupgrade{
			Gifts: t.toInt(parsedLine.tags["msg-param-months"]),
			Name:  parsedLine.tags["msg-param-promo-name"],
		}

	case "raid":
		user, _ := t.getUser(parsedLine.tags["msg-param-login"])
		user.DisplayName = parsedLine.tags["msg-param-displayName"]

		t.cEvents <- Raid{
			User:        user,
			ViewerCount: t.toInt(parsedLine.tags["msg-param-viewerCount"]),
		}

	case "unraid":
		t.cEvents <- Unraid{}

	case "ritual":
		t.cEvents <- Ritual{
			Name: parsedLine.tags["msg-param-ritual-name"],
		}

	case "bitsbadgetier":
		t.cEvents <- Cheer{
			BadgeTier: t.toInt(parsedLine.tags["msg-param-threshold"]),
		}

	}
}

// parseLine parses every line that was received from the IRC server
func (t *Twitch) parseLine(line string) {
	switch {
	case strings.HasPrefix(line, ":tmi.twitch.tv 001"):
		t.SendCommand("CAP REQ :twitch.tv/membership")
		t.SendCommand("CAP REQ :twitch.tv/tags")
		t.SendCommand("CAP REQ :twitch.tv/commands")
		for _, channel := range t.opts.Channels {
			t.JoinChannel(channel)
		}

		t.cEvents <- EventConnected{}
		return

	case strings.HasPrefix(line, "PING"):
		t.SendCommand("PONG :tmi.twitch.tv")

		t.cEvents <- EventPinged{}
		return

	case strings.HasPrefix(line, ":"+t.opts.Username+".tmi.twitch.tv 353"):
		channel, _ := t.getChannel(strings.Split(strings.Split(line, " :")[0], "#")[1])
		t.cEvents <- EventChannelJoined{Channel: channel}

	case strings.HasPrefix(line, ":tmi.twitch.tv NOTICE * :"):
		t.log.Error(strings.Split(line, " :")[1])
	}

	// parse every other event of the irc protocol

	var parsedLine parsedLine
	var timestamp time.Time
	var err error

	matches := regexp.MustCompile(`^(@(.+)\s+)?:(([^!]+).+)?tmi\.twitch\.tv\s+([A-Z]+)\s+#?(\w+)(\s+:(.+))?$`).FindAllStringSubmatch(line, -1)

	if matches != nil {
		parsedLine.channel, _ = t.getChannel(matches[0][6])
		parsedLine.user, _ = t.getUser(matches[0][4])
		parsedLine.t = matches[0][5]

		if strings.HasPrefix(line, "@") {
			parsedLine.tags = t.parseTags(matches[0][2])

			timestamp, _ = t.convertTmiTs(parsedLine.tags["tmi-sent-ts"])
			if parsedLine.tags["tmi-sent-ts"] != "" {
				timestamp, _ = t.convertTmiTs(parsedLine.tags["tmi-sent-ts"])
			}

			parsedLine.message = Message{
				Content: matches[0][8],
			}

			if tag, ok := parsedLine.tags["msg-id"]; ok && tag == "highlighted-message" {
				parsedLine.message.Highlighted = true
			} else if len(parsedLine.message.Content) > 0 && []byte(parsedLine.message.Content)[0] == byte(1) {
				parsedLine.message.Content = strings.TrimFunc(parsedLine.message.Content, func(r rune) bool {
					return !unicode.IsGraphic(r)
				})

				if strings.HasPrefix(parsedLine.message.Content, "ACTION ") {
					parsedLine.message.Me = true
					parsedLine.message.Content = parsedLine.message.Content[7:]
				}
			}

			if emotes, ok := parsedLine.tags["emotes"]; ok && emotes != "" {
				parsedLine.message.Emotes = make(map[string][]struct {
					From int
					To   int
				})
				for _, emote := range strings.Split(emotes, "/") {
					if strings.Contains(emote, ":") {
						emoteDetails := strings.Split(emote, ":")
						for _, r := range strings.Split(emoteDetails[1], ",") {
							fromTo := strings.Split(r, "-")

							f, err := strconv.Atoi(fromTo[0])
							if err != nil {
								continue
							}

							t, err := strconv.Atoi(fromTo[1])
							if err != nil {
								continue
							}

							parsedLine.message.Emotes[emoteDetails[0]] = append(parsedLine.message.Emotes[emoteDetails[0]], struct {
								From int
								To   int
							}{
								From: f,
								To:   t + 1,
							})
						}
					}
				}
			}
		}
	}

	if parsedLine.user != nil && parsedLine.channel != nil {
		t.addUserToChannel(parsedLine.user, parsedLine.channel)
	}

	switch parsedLine.t {
	case "PRIVMSG":
		parsedLine.message.ID = parsedLine.tags["id"]
		t.cEvents <- EventMessageReceived{Timestamp: timestamp, Channel: parsedLine.channel, Message: parsedLine.message, ChannelUser: t.buildChannelUser(&parsedLine)}

	case "JOIN":
		t.cEvents <- EventUserJoined{Channel: parsedLine.channel, User: parsedLine.user}

	case "PART":
		t.cEvents <- EventUserParted{Channel: parsedLine.channel, User: parsedLine.user}

	case "USERSTATE":
		event := EventUserstate{Channel: parsedLine.channel, User: parsedLine.user}
		if emoteSets, ok := parsedLine.tags["emote-sets"]; ok && emoteSets != "" {
			event.EmoteSets = strings.Split(emoteSets, ",")
		}
		t.cEvents <- event

	case "ROOMSTATE":
		parsedLine.channel.EmoteOnly, _ = strconv.ParseBool(parsedLine.tags["emote-only"])
		parsedLine.channel.FollowersOnly = parsedLine.tags["followers-only"] != "-1"
		parsedLine.channel.R9k, _ = strconv.ParseBool(parsedLine.tags["r9k"])
		parsedLine.channel.Slow, _ = strconv.ParseBool(parsedLine.tags["slow"])
		parsedLine.channel.SubsOnly, _ = strconv.ParseBool(parsedLine.tags["subs-only"])

		t.cEvents <- EventRoomstate{Channel: parsedLine.channel}

	case "USERNOTICE":
		// USERNOTICE holds many different events - is handled in own method
		t.parseUSERNOTICE(&parsedLine)

	case "CLEARCHAT":
		parsedLine.user, err = t.getUser(matches[0][8])
		if err != nil {
			// stop further processing, username is mandatory for this event
			t.log.Error("Invalid username on CLEARCHAT event", err)
			return
		}
		parsedLine.user.ID = t.toInt(parsedLine.tags["target-user-id"])

		t.cEvents <- EventClearchat{Timestamp: timestamp, BanDuration: t.toInt(parsedLine.tags["bad-duration"]), Channel: parsedLine.channel, User: parsedLine.user}

	case "CLEARMSG":
		parsedLine.user, err = t.getUser(parsedLine.tags["login"])
		if err != nil {
			return
		}
		parsedLine.message.ID = parsedLine.tags["target-msg-id"]

		t.cEvents <- EventClearmsg{Timestamp: timestamp, User: parsedLine.user, Channel: parsedLine.channel, Message: parsedLine.message}

	case "NOTICE":
		t.cEvents <- EventNotice{MsgId: parsedLine.tags["msg-id"], Channel: parsedLine.channel}

	case "HOSTTARGET":
		var viewers int64
		var toChannel string

		if strings.Contains(parsedLine.message.Content, " ") {
			parts := strings.Split(parsedLine.message.Content, " ")
			toChannel = parts[0]
			viewers = t.toInt(parts[1])
		} else {
			toChannel = parsedLine.message.Content
		}

		if toChannel != "-" {
			tC, _ := t.getChannel(toChannel)
			t.cEvents <- EventStartHosting{FromChannel: parsedLine.channel, ToChannel: tC, Viewers: viewers}
		} else {
			t.cEvents <- EventStopHosting{FromChannel: parsedLine.channel, Viewers: viewers}
		}

	// whisper has PRIVMSG like syntax - has to be parsed separately
	case "WHISPER":
		badges := t.parseBadges(parsedLine.tags["badges"])
		parsedLine.user.IsPartner = badges["partner"] == "1"
		parsedLine.user.DisplayName = parsedLine.tags["display-name"]
		parsedLine.user.Color = parsedLine.tags["color"]
		parsedLine.message.ID = parsedLine.tags["message-id"]
		t.cEvents <- EventWhisperReceived{User: parsedLine.user, Message: parsedLine.message}

	default:
		t.log.Info("unhandled event", line)
	}
}
