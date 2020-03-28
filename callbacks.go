package twirgo

type callbacks struct {
	connectionError     []func(*Twitch, EventConnectionError)
	connected           []func(*Twitch, EventConnected)
	pinged              []func(*Twitch, EventPinged)
	channelJoined       []func(*Twitch, EventChannelJoined)
	channelParted       []func(*Twitch, EventChannelParted)
	userJoined          []func(*Twitch, EventUserJoined)
	userParted          []func(*Twitch, EventUserParted)
	messageReceived     []func(*Twitch, EventMessageReceived)
	whisperReceived     []func(*Twitch, EventWhisperReceived)
	userstate           []func(*Twitch, EventUserstate)
	roomstate           []func(*Twitch, EventRoomstate)
	sub                 []func(*Twitch, EventSub)
	resub               []func(*Twitch, EventResub)
	subgift             []func(*Twitch, EventSubgift)
	anonsubgift         []func(*Twitch, EventAnonsubgift)
	submysterygift      []func(*Twitch, EventSubmysterygift)
	giftpaidupgrade     []func(*Twitch, EventGiftpaidupgrade)
	rewardgift          []func(*Twitch, EventRewardgift)
	anongiftpaidupgrade []func(*Twitch, EventAnongiftpaidupgrade)
	raid                []func(*Twitch, EventRaid)
	unraid              []func(*Twitch, EventUnraid)
	ritual              []func(*Twitch, EventRitual)
	bitsbadgetier       []func(*Twitch, EventBitsbadgetier)
	clearchat           []func(*Twitch, EventClearchat)
	clearmsg            []func(*Twitch, EventClearmsg)
	notice              []func(*Twitch, EventNotice)
	startHosting        []func(*Twitch, EventStartHosting)
	stopHosting         []func(*Twitch, EventStopHosting)
}

func (t *Twitch) callCallbacks(event interface{}) {
	switch ev := event.(type) {
	case EventConnectionError:
		for _, f := range t.callbacks.connectionError {
			go f(t, ev)
		}

	case EventConnected:
		for _, f := range t.callbacks.connected {
			go f(t, ev)
		}

	case EventPinged:
		for _, f := range t.callbacks.pinged {
			go f(t, ev)
		}

	case EventChannelJoined:
		for _, f := range t.callbacks.channelJoined {
			go f(t, ev)
		}

	case EventChannelParted:
		for _, f := range t.callbacks.channelParted {
			go f(t, ev)
		}

	case EventUserJoined:
		for _, f := range t.callbacks.userJoined {
			go f(t, ev)
		}

	case EventUserParted:
		for _, f := range t.callbacks.userParted {
			go f(t, ev)
		}

	case EventMessageReceived:
		for _, f := range t.callbacks.messageReceived {
			go f(t, ev)
			t.runCommand(ev)
		}

	case EventWhisperReceived:
		for _, f := range t.callbacks.whisperReceived {
			go f(t, ev)
		}

	case EventUserstate:
		for _, f := range t.callbacks.userstate {
			go f(t, ev)
		}

	case EventRoomstate:
		for _, f := range t.callbacks.roomstate {
			go f(t, ev)
		}

	case EventSub:
		for _, f := range t.callbacks.sub {
			go f(t, ev)
		}

	case EventResub:
		for _, f := range t.callbacks.resub {
			go f(t, ev)
		}

	case EventSubgift:
		for _, f := range t.callbacks.subgift {
			go f(t, ev)
		}

	case EventAnonsubgift:
		for _, f := range t.callbacks.anonsubgift {
			go f(t, ev)
		}

	case EventSubmysterygift:
		for _, f := range t.callbacks.submysterygift {
			go f(t, ev)
		}

	case EventGiftpaidupgrade:
		for _, f := range t.callbacks.giftpaidupgrade {
			go f(t, ev)
		}

	case EventRewardgift:
		for _, f := range t.callbacks.rewardgift {
			go f(t, ev)
		}

	case EventAnongiftpaidupgrade:
		for _, f := range t.callbacks.anongiftpaidupgrade {
			go f(t, ev)
		}

	case EventRaid:
		for _, f := range t.callbacks.raid {
			go f(t, ev)
		}

	case EventUnraid:
		for _, f := range t.callbacks.unraid {
			go f(t, ev)
		}

	case EventRitual:
		for _, f := range t.callbacks.ritual {
			go f(t, ev)
		}

	case EventBitsbadgetier:
		for _, f := range t.callbacks.bitsbadgetier {
			go f(t, ev)
		}

	case EventClearchat:
		for _, f := range t.callbacks.clearchat {
			go f(t, ev)
		}

	case EventClearmsg:
		for _, f := range t.callbacks.clearmsg {
			go f(t, ev)
		}

	case EventNotice:
		for _, f := range t.callbacks.notice {
			go f(t, ev)
		}

	case EventStartHosting:
		for _, f := range t.callbacks.startHosting {
			go f(t, ev)
		}

	case EventStopHosting:
		for _, f := range t.callbacks.stopHosting {
			go f(t, ev)
		}

	}
}

func (t *Twitch) OnConnectionError(f func(*Twitch, EventConnectionError)) {
	t.callbacks.connectionError = append(t.callbacks.connectionError, f)
}

func (t *Twitch) OnConnected(f func(*Twitch, EventConnected)) {
	t.callbacks.connected = append(t.callbacks.connected, f)
}

func (t *Twitch) OnPinged(f func(*Twitch, EventPinged)) {
	t.callbacks.pinged = append(t.callbacks.pinged, f)
}

func (t *Twitch) OnChannelJoined(f func(*Twitch, EventChannelJoined)) {
	t.callbacks.channelJoined = append(t.callbacks.channelJoined, f)
}

func (t *Twitch) OnChannelParted(f func(*Twitch, EventChannelParted)) {
	t.callbacks.channelParted = append(t.callbacks.channelParted, f)
}

func (t *Twitch) OnUserJoined(f func(*Twitch, EventUserJoined)) {
	t.callbacks.userJoined = append(t.callbacks.userJoined, f)
}

func (t *Twitch) OnUserParted(f func(*Twitch, EventUserParted)) {
	t.callbacks.userParted = append(t.callbacks.userParted, f)
}

func (t *Twitch) OnMessageReceived(f func(*Twitch, EventMessageReceived)) {
	t.callbacks.messageReceived = append(t.callbacks.messageReceived, f)
}

func (t *Twitch) OnWhisperReceived(f func(*Twitch, EventWhisperReceived)) {
	t.callbacks.whisperReceived = append(t.callbacks.whisperReceived, f)
}

func (t *Twitch) OnUserstate(f func(*Twitch, EventUserstate)) {
	t.callbacks.userstate = append(t.callbacks.userstate, f)
}

func (t *Twitch) OnRoomstate(f func(*Twitch, EventRoomstate)) {
	t.callbacks.roomstate = append(t.callbacks.roomstate, f)
}

func (t *Twitch) OnSub(f func(*Twitch, EventSub)) {
	t.callbacks.sub = append(t.callbacks.sub, f)
}

func (t *Twitch) OnResub(f func(*Twitch, EventResub)) {
	t.callbacks.resub = append(t.callbacks.resub, f)
}

func (t *Twitch) OnSubgift(f func(*Twitch, EventSubgift)) {
	t.callbacks.subgift = append(t.callbacks.subgift, f)
}

func (t *Twitch) OnAnonsubgift(f func(*Twitch, EventAnonsubgift)) {
	t.callbacks.anonsubgift = append(t.callbacks.anonsubgift, f)
}

func (t *Twitch) OnSubmysterygift(f func(*Twitch, EventSubmysterygift)) {
	t.callbacks.submysterygift = append(t.callbacks.submysterygift, f)
}

func (t *Twitch) OnGiftpaidupgrade(f func(*Twitch, EventGiftpaidupgrade)) {
	t.callbacks.giftpaidupgrade = append(t.callbacks.giftpaidupgrade, f)
}

func (t *Twitch) OnRewardgift(f func(*Twitch, EventRewardgift)) {
	t.callbacks.rewardgift = append(t.callbacks.rewardgift, f)
}

func (t *Twitch) OnAnongiftpaidupgrade(f func(*Twitch, EventAnongiftpaidupgrade)) {
	t.callbacks.anongiftpaidupgrade = append(t.callbacks.anongiftpaidupgrade, f)
}

func (t *Twitch) OnRaid(f func(*Twitch, EventRaid)) {
	t.callbacks.raid = append(t.callbacks.raid, f)
}

func (t *Twitch) OnUnraid(f func(*Twitch, EventUnraid)) {
	t.callbacks.unraid = append(t.callbacks.unraid, f)
}

func (t *Twitch) OnRitual(f func(*Twitch, EventRitual)) {
	t.callbacks.ritual = append(t.callbacks.ritual, f)
}

func (t *Twitch) OnBitsbadgetier(f func(*Twitch, EventBitsbadgetier)) {
	t.callbacks.bitsbadgetier = append(t.callbacks.bitsbadgetier, f)
}

func (t *Twitch) OnClearchat(f func(*Twitch, EventClearchat)) {
	t.callbacks.clearchat = append(t.callbacks.clearchat, f)
}

func (t *Twitch) OnClearmsg(f func(*Twitch, EventClearmsg)) {
	t.callbacks.clearmsg = append(t.callbacks.clearmsg, f)
}

func (t *Twitch) OnNotice(f func(*Twitch, EventNotice)) {
	t.callbacks.notice = append(t.callbacks.notice, f)
}

func (t *Twitch) OnStartHosting(f func(*Twitch, EventStartHosting)) {
	t.callbacks.startHosting = append(t.callbacks.startHosting, f)
}

func (t *Twitch) OnStopHosting(f func(*Twitch, EventStopHosting)) {
	t.callbacks.stopHosting = append(t.callbacks.stopHosting, f)
}
