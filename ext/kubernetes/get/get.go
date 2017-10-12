package get

import (
	"fmt"
	"strings"

	"github.com/signavio/signa/pkg/bot"
	"github.com/signavio/signa/pkg/kubectl"
	"github.com/signavio/signa/pkg/logger"
)

const (
	invalidAmountOfParams = "Invalid amount of parameters"
	invalidParams         = "Invalid parameters"
)

func init() {
	bot.RegisterCommand(
		"get",
		"Query resources running in the Kubernetes cluster.",
		"-n <namespace> pods|services|deployments|namespaces|etc",
		Get)
}

func Get(c *bot.Cmd) (string, error) {
	if len(c.Args) < 1 {
		return invalidAmountOfParams, nil
	}

	args := parseArgs(c.Args)
	for _, a := range args {
		if strings.Contains(a, "secret") {
			return "not so fast :wink:", nil
		}
	}

	k, err := kubectl.NewKubectl("default", args)
	checkErrors(err)
	o, err := k.Exec()
	checkErrors(err)

	if err != nil {
		return invalidParams, nil
	}

	r := fmt.Sprintf("```%v```", o)
	// NOTE: move to the bot package
	//logAction(c.User.Nick, c.Channel, c.Command, c.RawArgs)

	return r, err
}

func parseArgs(args []string) (sl []string) {
	sl = []string{"get"}
	for _, a := range args {
		sl = append(sl, a)
	}

	return
}

func checkErrors(e error) {
	if e != nil {
		logger.Error.Println(e)
	}
}
