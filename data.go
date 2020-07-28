package twirgo

type (
	User struct {
		ID          int64
		Username    string
		DisplayName string
		Color       string
		IsPartner   bool
	}

	ChannelUser struct {
		User             *User
		SubscriberMonths int64
		Badges           Badges

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
		ID      string
		Emotes  map[string][]struct {
			From int
			To   int
		}
	}

	SubTier interface{}

	SubTierPrime struct {
		Name string
	}
	SubTierOne struct {
		Name string
	}
	SubTierTwo struct {
		Name string
	}
	SubTierThree struct {
		Name string
	}

	Sub struct {
		Months       int64    // msg-param-cumulative-months
		ShareStreak  bool     // msg-param-should-share-streak
		StreakMonths int64    // msg-param-streak-months
		SubTier      *SubTier // msg-param-sub-plan & msg-param-sub-plan-name - can hold SubTierPrime, SubTierOne, SubTierTwo, SubTierThree
	}
	Resub struct {
		Months       int64    // msg-param-cumulative-months
		ShareStreak  bool     // msg-param-should-share-streak
		StreakMonths int64    // msg-param-streak-months
		SubTier      *SubTier // msg-param-sub-plan & msg-param-sub-plan-name - can hold SubTierPrime, SubTierOne, SubTierTwo, SubTierThree
	}
	Subgift struct {
		Months  int64    // msg-param-months
		User    *User    // msg-param-recipient-display-name & msg-param-recipient-id & msg-param-recipient-user-name
		SubTier *SubTier // msg-param-sub-plan & msg-param-sub-plan-name - can hold SubTierPrime, SubTierOne, SubTierTwo, SubTierThree
	}
	Anonsubgift struct {
		Months  int64    // msg-param-months
		User    *User    // msg-param-recipient-display-name & msg-param-recipient-id & msg-param-recipient-user-name
		SubTier *SubTier // msg-param-sub-plan & msg-param-sub-plan-name - can hold SubTierPrime, SubTierOne, SubTierTwo, SubTierThree
	}
	Submysterygift  struct{}
	Giftpaidupgrade struct {
		Gifts int64  // msg-param-months - number of gifts the user gifted
		Name  string // msg-param-promo-name
		User  *User  // msg-param-sender-login & msg-param-sender-name
	}
	Rewardgift          struct{}
	Anongiftpaidupgrade struct {
		Gifts int64  // msg-param-months - number of gifts the user gifted
		Name  string // msg-param-promo-name
	}
	Raid struct {
		User        *User // msg-param-displayName & msg-param-login ignored
		ViewerCount int64 // msg-param-viewerCount
	}
	Unraid struct{}
	Ritual struct {
		Name string // msg-param-ritual-name
	}
	Cheer struct {
		BadgeTier int64 // msg-param-threshold
	}

	Badges map[string]int64

	Tags map[string]string
)
