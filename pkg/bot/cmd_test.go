package bot

import (
	"errors"
	"fmt"
	"reflect"
	"strings"
	"testing"
)

var (
	channel string
	replies []string
	user    *User
)

const (
	expectedMsg    = "msg"
	cmd            = "cmd"
	cmdDescription = "Command description"
	cmdExampleArgs = "arg1 arg2"
)

func responseHandler(target string, message string, sender *User) {
	channel = target
	user = sender
	replies = append(replies, message)
}

func resetResponses() {
	channel = ""
	user = &User{Nick: ""}
	replies = []string{}
	commands = make(map[string]*customCommand)
}

func newBot() *Bot {
	return New(&Handler{
		Response: responseHandler,
	})
}

func registerValidCommand() {
	RegisterCommand(cmd, cmdDescription, cmdExampleArgs,
		func(c *Cmd) (string, error) {
			return expectedMsg, nil
		})
}

func TestCommandNotRegistered(t *testing.T) {
	resetResponses()

	newBot().MessageReceived(&ChannelData{Channel: "#go-bot"}, &Message{Text: "!not_a_cmd"}, &User{})

	if len(replies) != 0 {
		t.Fatal("Should not reply if a command is not found")
	}
}

func TestInvalidCmdArgs(t *testing.T) {
	resetResponses()
	registerValidCommand()

	newBot().MessageReceived(&ChannelData{Channel: "#go-bot"}, &Message{Text: "!cmd \"invalid arg"}, &User{Nick: "user"})

	if channel != "#go-bot" {
		t.Error("Should reply to #go-bot channel")
	}
	if len(replies) != 1 {
		t.Fatal("Invalid reply")
	}
	if !strings.HasPrefix(replies[0], "Error parsing") {
		t.Fatal("Should reply with an error message")
	}
}

func TestErroredCmd(t *testing.T) {
	resetResponses()
	cmdError := errors.New("error")
	RegisterCommand("cmd", "", "",
		func(c *Cmd) (string, error) {
			return "", cmdError
		})

	newBot().MessageReceived(&ChannelData{Channel: "#go-bot"}, &Message{Text: "!cmd"}, &User{Nick: "user"})

	if channel != "#go-bot" {
		t.Fatal("Invalid channel")
	}
	if len(replies) != 1 {
		t.Fatal("Invalid reply")
	}
	if replies[0] != fmt.Sprintf(errorExecutingCommand, "cmd", cmdError.Error()) {
		t.Fatal("Reply should contain the error message")
	}
}

func TestValidCmdOnChannel(t *testing.T) {
	resetResponses()
	registerValidCommand()

	newBot().MessageReceived(&ChannelData{Channel: "#go-bot"}, &Message{Text: "!cmd"}, &User{Nick: "user"})

	if channel != "#go-bot" {
		t.Fatal("Command called on channel should reply to channel")
	}
	if len(replies) != 1 {
		t.Fatal("Should have one reply on channel")
	}
	if replies[0] != expectedMsg {
		t.Fatal("Invalid command reply")
	}
}

func TestChannelData(t *testing.T) {
	cd := ChannelData{
		Protocol: "irc",
		Server:   "myserver",
		Channel:  "#mychan",
	}
	if cd.URI() != "irc://myserver/#mychan" {
		t.Fatal("URI should return a valid IRC URI")
	}
}

func TestHelpWithNoArgs(t *testing.T) {
	resetResponses()
	registerValidCommand()
	newBot().MessageReceived(&ChannelData{Channel: "#go-bot"}, &Message{Text: "!help"}, &User{Nick: "user"})

	expectedReply := []string{
		fmt.Sprintf(helpAboutCommand, CmdPrefix),
		fmt.Sprintf(availableCommands, "cmd"),
	}

	if !reflect.DeepEqual(replies, expectedReply) {
		t.Fatalf("Invalid reply. Expected %v got %v", expectedReply, replies)
	}
}

func TestHelpForACommand(t *testing.T) {
	resetResponses()
	registerValidCommand()
	newBot().MessageReceived(&ChannelData{Channel: "#go-bot"}, &Message{Text: "!help cmd"}, &User{Nick: "user"})

	expectedReply := []string{
		fmt.Sprintf(helpDescripton, cmdDescription),
		fmt.Sprintf(helpUsage, CmdPrefix, cmd, cmdExampleArgs),
	}

	if !reflect.DeepEqual(replies, expectedReply) {
		t.Fatalf("Invalid reply. Expected %v got %v", expectedReply, replies)
	}
}

func TestHelpWithNonExistingCommand(t *testing.T) {
	resetResponses()
	registerValidCommand()
	newBot().MessageReceived(&ChannelData{Channel: "#go-bot"}, &Message{Text: "!help not_a_cmd"}, &User{Nick: "user"})

	expectedReply := []string{
		fmt.Sprintf(helpAboutCommand, CmdPrefix),
		fmt.Sprintf(availableCommands, "cmd"),
	}

	if !reflect.DeepEqual(replies, expectedReply) {
		t.Fatalf("Invalid reply. Expected %v got %v", expectedReply, replies)
	}
}

func TestHelpWithInvalidArgs(t *testing.T) {
	resetResponses()
	registerValidCommand()
	newBot().MessageReceived(&ChannelData{Channel: "#go-bot"}, &Message{Text: "!help cmd \"invalid arg"}, &User{Nick: "user"})

	if len(replies) != 1 {
		t.Fatal("Invalid reply")
	}
	if !strings.HasPrefix(replies[0], "Error parsing") {
		t.Fatal("Should reply with an error message")
	}
}
