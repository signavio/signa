package jobs

import (
	"fmt"

	"github.com/signavio/signa/pkg/bot"
)

const (
	invalidAmountOfParams = "Invalid amount of parameters"
	jobNotFound           = "Job not found."
	permissionDenied      = "You are not allowed to execute this operation. :sweat_smile:"
	errorMessage          = "Something went wrong."
)

func init() {
	bot.RegisterCommand(
		"run",
		"Run kubernetes jobs.",
		"<job-name>",
		Run,
	)
}

func Run(c *bot.Cmd) (string, error) {
	if len(c.Args) < 1 {
		return invalidAmountOfParams, nil
	}

	job := bot.Cfg().FindJob(c.Args[0])
	if job == nil {
		return jobNotFound, nil
	}

	// TODO: Implement check in a global level.
	username := c.User.Nick
	if bot.Cfg().IsSuperuser(username) || job.IsExecUser(username) {
		j := NewJob(job)
		_, err := j.createJob()
		if err != nil {
			return errorMessage, err
		}

		pods, err := j.getJobPods()
		if err != nil {
			return errorMessage, err
		}

		logs, err := j.getJobLogs(pods)
		if err != nil {
			return errorMessage, err
		}

		return fmt.Sprintf("```%s```", logs), nil
	} else {
		return permissionDenied, nil
	}

	return "", nil
}
