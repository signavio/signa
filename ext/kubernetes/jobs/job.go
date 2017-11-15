package jobs

import (
	"html/template"
	"log"
	"os"

	"github.com/signavio/signa/pkg/bot"
)

type Job struct {
	*bot.Job
	ImageTag string
}

func NewJob(j *bot.Job, imageTag string) *Job {
	return &Job{j, imageTag}
}

func (j *Job) exec(status chan string) {
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
			break
		}

		if state == "Completed" {
			logs, err := j.getJobLogs(pods)
			if err != nil {
				log.Print(err)
				status <- errorMessage
				break
			}
			status <- logs
			break
		} else if state == "ErrImagePull" {
			log.Print(err)
			status <- errorMessage + ": " + state
			break
		}
		// TODO: Implement error condition.
	}

	if _, err := j.deleteJob(); err != nil {
		log.Print(err)
		status <- errorMessage
		return
	}

	close(status)
}

func (j *Job) createJob() (string, error) {
	cmd := NewCommand([]string{"create", "-f", j.Config, "-n", j.Namespace})
	output, err := cmd.Exec()
	if err != nil {
		return "", err
	}
	return output, nil
}

func (j *Job) getJobPods() (string, error) {
	cmd := NewCommand([]string{
		"get",
		"pods",
		"--show-all",
		"--selector=job-name=" + j.Name,
		"--output=jsonpath={.items..metadata.name}",
		"-n",
		j.Namespace,
	})
	output, err := cmd.Exec()
	if err != nil {
		return "", err
	}
	return output, nil
}

func (j *Job) getJobState() (string, error) {
	cmd := NewCommand([]string{
		"get",
		"pods",
		"--show-all",
		"--selector=job-name=" + j.Name,
		"--output=jsonpath={.items..status.containerStatuses..reason}",
		"-n",
		j.Namespace,
	})
	output, err := cmd.Exec()
	if err != nil {
		return "", err
	}
	return output, nil
}

func (j *Job) getJobLogs(pods string) (string, error) {
	cmd := NewCommand([]string{"logs", pods, "--tail=20", "-n", j.Namespace})
	output, err := cmd.Exec()
	if err != nil {
		return "", err
	}
	return output, nil
}

func (j *Job) deleteJob() (string, error) {
	cmd := NewCommand([]string{"delete", "-f", j.Config, "-n", j.Namespace})
	output, err := cmd.Exec()
	if err != nil {
		return "", err
	}
	return output, nil
}

func (j *Job) parseImageTag() error {
	t, err := template.ParseFiles(j.Config)
	if err != nil {
		return err
	}

	parsedCfg := j.Config + ".parsed"
	f, err := os.Create(parsedCfg)
	if err != nil {
		return err
	}

	err = t.Execute(f, j)
	if err != nil {
		return err
	}

	j.Config = parsedCfg

	return nil
}
