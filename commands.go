package twirgo

import (
	"strings"
)

type Command struct {
	Command    string
	Parameters []string
	Text       string
	Mentions   []*User
	Message    *EventMessageReceived
}

type commandPermission int8

const (
	PermEveryone commandPermission = iota
	PermMods
	PermBroadcaster
)

type commandDefinition struct {
	channels   []string
	permission commandPermission
	callback   func(*Twitch, Command)
}

// Command saved the command definition for an user defined command
func (t *Twitch) Command(cs []string, channels []string, permission commandPermission, cb func(*Twitch, Command)) {
	for _, c := range cs {
		c = strings.ToLower(strings.TrimSpace(c))
		t.commands[c] = &commandDefinition{
			channels:   channels,
			callback:   cb,
			permission: permission,
		}
	}
}

func (t *Twitch) runCommand(message EventMessageReceived) {
	if len(t.commands) < 1 {
		return
	}

	for cmd, def := range t.commands {
		if strings.HasPrefix(strings.ToLower(message.Message.Content), cmd) {
			if !t.inSlice(def.channels, message.Channel.Name) ||
				(def.permission == PermBroadcaster && !message.ChannelUser.IsBroadcaster) ||
				(def.permission == PermMods && !message.ChannelUser.IsMod) {
				return
			}

			parameters := strings.Split(message.Message.Content, " ")[1:]

			command := Command{
				Command:    cmd,
				Parameters: parameters,
				Text:       strings.Join(parameters, " "),
				Message:    &message,
			}

			for _, p := range parameters {
				if strings.HasPrefix(p, "@") {
					username := strings.ToLower(strings.TrimLeft(p, "@"))
					user, _ := t.getUser(username)
					command.Mentions = append(command.Mentions, user)
				}
			}

			go def.callback(t, command)

			return
		}
	}
}
