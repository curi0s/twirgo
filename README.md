# TWIRGO

**TW**itch**IR**c**GO** is a Twitch IRC library written in GO

You have the choice of using this library with callbacks, handling the whole channel by yourself or define commands directly.
For either way there is a example in the [examples directory](examples/).

## Usage

You have to pass a `twirgo.Options` slice to the `New` method of TWIRGO.

| Option         | Description                                                                                            | Mandatory |
| -------------- | ------------------------------------------------------------------------------------------------------ | --------- |
| Username       | The username of your bot account                                                                       | X         |
| Token          | An oauth token (format: oauth:xxxx)                                                                    | X         |
| Channels       | A slice of channelnames to connect to at start                                                         | X         |
| Log            | A [logrus](https://github.com/sirupsen/logrus) instance                                                | X         |
| DefaultChannel | Not necessary for the library. Have a look at [examples/channel/main.go](examples/channel/main.go#L16) |           |
| Unsecure       | Should the connection unsecure to the Twitch IRC server?                                               |           |

```golang
options := twirgo.Options{
    Username:       "curi_bot_",                   // the name of your bot account
    Token:          os.Getenv("TOKEN"),            // provide your token in any way you like
    Channels:       []string{"curi", "curi_bot_"}, // all channels will be joined at connect
    DefaultChannel: "curi",
    Log:            logrus.New(),
    Unsecure:       false,
}

// You can set the loglevel (default is Info) or even change the formatter
// => https://github.com/sirupsen/logrus
// Be aware that the debug level will reveal your token in the log!
// options.Log.SetLevel(logrus.DebugLevel)

t := twirgo.New(options)
```

## Examples

Check out the [examples](examples/) directory.

-   [Channel](examples/channel/main.go)
-   [Callbacks](examples/callbacks/main.go)
-   [Commands](examples/commands/main.go)
