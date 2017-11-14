package jobs

import (
	"fmt"
	"log"

	"github.com/signavio/signa/pkg/bot"
)

const (
	invalidAmountOfParams = "Invalid amount of parameters"
	jobNotFound           = "Job not found."
	permissionDenied      = "You are not allowed to execute this operation. :sweat_smile:"
	errorMessage          = "Something went wrong."
	jobOutputNotFound     = "Job executed but output not found."
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
		status := make(chan string)
		go execJob(j, status)

		for {
			if current := <-status; current != "" {
				return fmt.Sprintf("```%s```", current), nil
			}
		}
	} else {
		return permissionDenied, nil
	}

	return jobOutputNotFound, nil
}

func execJob(j *Job, status chan string) {
	if _, err := j.createJob(); err != nil {
		log.Print(err)
		status <- errorMessage
		return
	}

	pods, err := j.getJobPods()
	if err != nil {
		log.Print(err)
		status <- errorMessage
		return
	}

	for {
		state, err := j.getJobState()
		if err != nil {
			log.Print(err)
			status <- errorMessage
			return
		}

		if state == "Completed" {
			logs, err := j.getJobLogs(pods)
			if err != nil {
				log.Print(err)
				status <- errorMessage
				return
			}
			status <- logs
			break
		}
	}

	if _, err := j.deleteJob(); err != nil {
		log.Print(err)
		status <- errorMessage
		return
	}

	close(status)
}
