// Package bot provides a simple framework to create Slack bots
package bot

import (
	"fmt"
	"log"
	"math/rand"
	"time"
)

const (
	// CmdPrefix is the identifier of a command
	// The current definition implies that the prefix for calling
	// a command is "!", therefore !hello would be read as a command
	// by the bot in the chat stream.
	CmdPrefix = "!"
)

// A Bot implements a Handler.
type Bot struct {
	Handler *Handler
}

// A ResponseHandler handles the bot responses
type ResponseHandler func(target, message string, sender *User)

// A Handler receives callbacks from the bot
type Handler struct {
	Response ResponseHandler
}

// New creates a new bot instance
func New(h *Handler) *Bot {
	return &Bot{Handler: h}
}

// MessageReceived is called by the protocol upon receiving a message
func (b *Bot) MessageReceived(channel *ChannelData, message *Message, sender *User) {
	command, err := parse(message.Text, channel, sender)
	if err != nil {
		b.Handler.Response(channel.Channel, err.Error(), sender)
		return
	}

	if command != nil {
		switch command.Command {
		case helpCommand:
			b.help(command)
		default:
			b.handleCmd(command)
		}
	}
}

func (b *Bot) handleCmd(c *Cmd) {
	cmd := commands[c.Command]

	if cmd == nil {
		log.Printf("Command not found %v", c.Command)
		return
	}

	message, err := cmd.CmdFunc(c)
	b.checkCmdError(err, c)
	if message != "" {
		b.Handler.Response(c.Channel, message, c.User)
	}
}

func (b *Bot) checkCmdError(err error, c *Cmd) {
	if err != nil {
		errorMsg := fmt.Sprintf(errorExecutingCommand, c.Command, err.Error())
		log.Printf(errorMsg)
		b.Handler.Response(c.Channel, errorMsg, c.User)
	}
}

func Cfg() *Config {
	return cfg
}

func init() {
	rand.Seed(time.Now().UnixNano())

	if err := cfg.Load("/etc/signa.conf"); err != nil {
		panic(err)
	}
	if err := setupLogger(); err != nil {
		panic(err)
	}
}
