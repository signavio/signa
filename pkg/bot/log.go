package bot

import (
	"fmt"
	"io/ioutil"
	"os"

	"github.com/signavio/signa/pkg/logger"
)

const slackActionLogMessage = "User %v on Channel %v executed: %v %v"

func setupLogger() {
	// NOTE: Using stdout and stderr for now.
	//f, err := os.OpenFile(cfg.Log, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	logger.NewLogger(ioutil.Discard, os.Stdout, os.Stdout, os.Stderr)
}

func LogError(err error) {
	logger.Error.Println(err)
}

func LogSlackAction(username, channel, command, arg string) {
	logger.Info.Println(fmt.Sprintf(
		slackActionLogMessage,
		username,
		channel,
		command,
		arg,
	))
}
