package info

import (
	"errors"
	"fmt"
	"strings"

	"github.com/signavio/signa/pkg/bot"
	"github.com/signavio/signa/pkg/kubectl"
)

const (
	invalidAmountOfParams = "Invalid amount of parameters"
	invalidParams         = "Invalid parameters"

	noNamespaceNorDeployment = "No Namespace nor Deployment found in argument"
	currentImageVersion      = "Current deployed image and version for `%s`: ```%s```"
)

func init() {
	bot.RegisterCommand(
		"info",
		"Retrieve information from Deployments in the k8s Cluster.",
		"version <namespace>/<deployment-name>",
		Info,
	)
}

func Info(c *bot.Cmd) (string, error) {
	if len(c.Args) < 1 {
		return invalidAmountOfParams, nil
	}
	if c.Args[0] != "version" {
		return invalidParams, nil
	}

	ns, depl, err := parseDeploymentAndNamespace(c.Args[1])
	if err != nil {
		// NOTE: Implement general logging later.
		return "", err
	}
	args := []string{
		"get",
		"deployment",
		depl,
		"-o=jsonpath='{$.spec.template.spec.containers[:1].image}'",
		"-n",
		ns,
	}
	k, err := kubectl.NewKubectl("default", args)
	if err != nil {
		// NOTE: Implement general logging later.
		return "", err
	}

	output, err := k.Exec()
	if err != nil {
		// NOTE: Implement general logging later.
		return "", err
	}

	return fmt.Sprintf(currentImageVersion, depl, output), nil
}

func parseDeploymentAndNamespace(s string) (string, string, error) {
	sl := strings.Split(s, "/")
	if len(sl) != 2 {
		return "", "", errors.New(noNamespaceNorDeployment)
	}
	for _, str := range sl {
		if str == "" {
			return "", "", errors.New(noNamespaceNorDeployment)
		}
	}
	return sl[0], sl[1], nil
}
