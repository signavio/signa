package exec

import (
	"fmt"

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

func extractArgs(namespace string, a []string) (sl []string) {
	pod := a[0]
	a = append(a[:0], a[0+1:]...)
	sl = []string{"exec", "-n", namespace, pod, "--"}
	for _, arg := range a {
		sl = append(sl, arg)
	}

	return
}

func Exec(c *bot.Cmd) (r string, err error) {
	if len(c.Args) < 2 {
		return invalidAmountOfParams, nil
	}
	// NOTE: Find a way of set the namespace where to execute.
	namespace := "foobar"
	a := extractArgs(namespace, c.Args)

	bin, err := kubectl.WhereIs()
	checkErrors(err)

	k, _ := kubectl.NewKubectl(bin, a)
	o, err := k.Exec()
	checkErrors(err)

	if err != nil {
		return invalidParams, nil
	}

	r = fmt.Sprintf("```%v```", o)

	return
}

func init() {
	bot.RegisterCommand(
		"exec",
		"Run commands inside a pod in the Kubernetes cluster.",
		"<pod> <command>",
		Exec)
}
