package deployment

import (
	"encoding/json"
	"io"

	"github.com/nlopes/slack"
	"github.com/signavio/signa/pkg/bot"
	"github.com/signavio/signa/pkg/kubectl"
)

func executeKubectlCmd(namespace string, args ...string) (string, error) {
	baseArgs := []string{"-n", namespace}
	kubectlArgs := append(baseArgs, args...)

	k, err := kubectl.NewKubectl("default", kubectlArgs)
	if err != nil {
		return "", err
	}

	return k.Exec()
}

// NOTE: move this function to the bot package. It should have an API
// for sending arbitrary messages to Slack.
func postMessageToSlackChannel(channel, message string) error {
	slackApi := slack.New(bot.Cfg().SlackToken)
	messageParams := slack.PostMessageParameters{
		Username: bot.Cfg().BotUsername,
		AsUser:   true,
	}

	_, _, err := slackApi.PostMessage(channel, message, messageParams)
	if err != nil {
		return err
	}

	return nil
}

func decodeJson(reader io.Reader, receiver interface{}) error {
	parser := json.NewDecoder(reader)
	return parser.Decode(&receiver)
}
