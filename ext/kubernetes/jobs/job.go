package jobs

import (
	"crypto/rand"
	"fmt"
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

func (j *Job) exec(kubeconfig string, status chan string) {
	if _, err := j.createJob(kubeconfig); err != nil {
		log.Print(err)
		status <- errorMessage
		return
	}

	for {
		state, err := j.getJobState(kubeconfig)
		if err != nil {
			log.Print(err)
			status <- errorMessage
			break
		}

		if state == "Completed" || state == "Error" {
			logs, err := j.getJobLogs(kubeconfig)
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
	}

	if _, err := j.deleteJob(kubeconfig); err != nil {
		log.Print(err)
		status <- errorMessage
		return
	}

	close(status)
}

func (j *Job) createJob(kubeconfig string) (string, error) {
	cmd := NewCommand([]string{
		"--kubeconfig=" + kubeconfig,
		"create",
		"-f",
		j.Config,
		"-n",
		j.Namespace,
	})
	output, err := cmd.Exec()
	if err != nil {
		return "", err
	}
	return output, nil
}

func (j *Job) getJobState(kubeconfig string) (string, error) {
	cmd := NewCommand([]string{
		"--kubeconfig=" + kubeconfig,
		"get",
		"pods",
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

func (j *Job) getJobLogs(kubeconfig string) (string, error) {
	pods, err := j.getJobPods(kubeconfig)
	if err != nil {
		return "", err
	}

	cmd := NewCommand([]string{
		"--kubeconfig=" + kubeconfig,
		"logs",
		pods,
		"--tail=20",
		"-n",
		j.Namespace,
	})
	output, err := cmd.Exec()
	if err != nil {
		return "", err
	}

	return output, nil
}

func (j *Job) getJobPods(kubeconfig string) (string, error) {
	cmd := NewCommand([]string{
		"--kubeconfig=" + kubeconfig,
		"get",
		"pods",
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

func (j *Job) deleteJob(kubeconfig string) (string, error) {
	cmd := NewCommand([]string{
		"--kubeconfig=" + kubeconfig,
		"delete",
		"-f",
		j.Config,
		"-n",
		j.Namespace,
	})
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

	parsedCfgFile := "/tmp/parsed-" + randToken() + ".yaml"
	f, err := os.Create(parsedCfgFile)
	if err != nil {
		return err
	}

	err = t.Execute(f, j)
	if err != nil {
		return err
	}

	j.Config = parsedCfgFile

	return nil
}

func randToken() string {
	b := make([]byte, 8)
	rand.Read(b)
	return fmt.Sprintf("%x", b)
}
