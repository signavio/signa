package bot

import "fmt"

// Cmd holds the parsed user's input for easier handling of commands
type Cmd struct {
	Raw         string       // Raw is full string passed to the command
	Channel     string       // Channel where the command was called
	ChannelData *ChannelData // More info about the channel, including network
	User        *User        // User who sent the message
	Message     string       // Full string without the prefix
	MessageData *Message     // Message with extra flags
	Command     string       // Command is the first argument passed to the bot
	RawArgs     string       // Raw arguments after the command
	Args        []string     // Arguments as array
}

// ChannelData holds the improved channel info, which includes protocol and server
type ChannelData struct {
	Protocol  string // What protocol the message was sent on (irc, slack, telegram)
	Server    string // The server hostname the message was sent on
	Channel   string // The channel name the message appeared in
	IsPrivate bool   // Whether the channel is a group or private chat
}

// Message holds the message info - for IRC and Slack networks, this can include whether the message was an action.
type Message struct {
	Text     string // The actual content of this Message
	IsAction bool   // True if this was a '/me does something' message
}

// User holds user id, nick and real name
type User struct {
	ID       string
	Nick     string
	RealName string
	IsBot    bool
}

type customCommand struct {
	Cmd         string
	CmdFunc     activeCmdFunc
	Description string
	ExampleArgs string
}

const (
	commandNotAvailable   = "Command %v not available."
	noCommandsAvailable   = "No commands available."
	errorExecutingCommand = "Error executing %s: %s"
)

type activeCmdFunc func(cmd *Cmd) (string, error)

var commands = make(map[string]*customCommand)

// URI gives back an URI-fied string containing protocol, server and channel.
func (c *ChannelData) URI() string {
	return fmt.Sprintf("%s://%s/%s", c.Protocol, c.Server, c.Channel)
}

// RegisterCommand adds a new command to the bot.
// The command(s) should be registered in the Init() func of your package
// command: String which the user will use to execute the command, example: reverse
// decription: Description of the command to use in !help, example: Reverses a string
// exampleArgs: Example args to be displayed in !help <command>, example: string to be reversed
// cmdFunc: Function which will be executed. It will received a parsed command as a Cmd value
func RegisterCommand(command, description, exampleArgs string, cmdFunc activeCmdFunc) {
	commands[command] = &customCommand{
		Cmd:         command,
		CmdFunc:     cmdFunc,
		Description: description,
		ExampleArgs: exampleArgs,
	}
}
