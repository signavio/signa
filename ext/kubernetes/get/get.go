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

func checkErrors(e error) {
	if e != nil {
		logger.Error.Println(e)
	}
}

func extractArgs(args []string) (sl []string) {
	sl = []string{"get"}
	for _, a := range args {
		sl = append(sl, a)
	}

	return
}

func Get(c *bot.Cmd) (string, error) {
	if len(c.Args) < 1 {
		return invalidAmountOfParams, nil
	}
	a := extractArgs(c.Args)
	for _, arg := range a {
		if strings.Contains(arg, "secret") {
			return "not so fast :wink:", nil
		}
	}

	k, err := kubectl.NewKubectl("default", a)
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

func init() {
	bot.RegisterCommand(
		"get",
		"Query resources running in the Kubernetes cluster.",
		"-n <namespace> pods|services|deployments|namespaces|etc",
		Get)
}
