package deployment

import (
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/signavio/signa/pkg/bot"
)

var k8sErrorStatuses = []string{
	"ImagePullBackOff",
	"ImageInspectError",
	"ErrImagePull",
	"ErrImageNeverPull",
	"RegistryUnavailable",
	"InvalidImageName",
	"CrashLoopBackOff",
	"RunContainerError",
}

type Deployment struct {
	Name          string
	Namespace     string
	ContainerName string
	RepositoryURI string
	Version       string
	Config        string
}

func NewDeployment(component *bot.Component, container *bot.Container, version string) *Deployment {
	return &Deployment{
		Name:          component.Name,
		Namespace:     component.Namespace,
		ContainerName: container.Name,
		RepositoryURI: container.RepositoryURI,
		Version:       version,
		Config:        component.BootstrapConfig,
	}
}

func (d *Deployment) Apply() (bool, error) {
	deploymentName := fmt.Sprintf("deployment/%v", d.Name)
	currentDeployment, _ := executeKubectlCmd(d.Namespace, "get", deploymentName)

	if strings.Contains(currentDeployment, "NotFound") {
		_, err := executeKubectlCmd(d.Namespace, "create", "-f", d.Config)
		if err != nil {
			return false, err
		}

		return false, nil
	} else {
		image := fmt.Sprintf("%v=%v:%v", d.ContainerName, d.RepositoryURI, d.Version)

		_, err := executeKubectlCmd(d.Namespace, "set", "image", deploymentName, image)
		if err != nil {
			bot.LogError(err)
			return false, err
		}

		time.Sleep(time.Duration(bot.Cfg().RollbackCheck) * time.Second)
		return d.rollbackInCaseOfError()
	}

	return false, errors.New("Apply(): Some error happened.")
}

func (d *Deployment) GetPods() []string {
	pods := NewPods(d.Namespace)
	deployedPods := []string{}
	for _, p := range pods.Items {
		if strings.Contains(p.Metadata.Name, d.Name) {
			deployedPods = append(deployedPods, p.Metadata.Name)
		}
	}

	return deployedPods
}

func (d *Deployment) rollbackInCaseOfError() (bool, error) {
	if !d.isDeploySuccessful() {
		_, err := executeKubectlCmd(d.Namespace, "rollout", "undo", "deployment", d.Name)
		if err != nil {
			bot.LogError(err)
			return false, err
		}

		return true, nil
	}

	return false, nil
}

func (d *Deployment) isDeploySuccessful() bool {
	deploymentStatus := d.getDeploymentStatus()
	for _, e := range k8sErrorStatuses {
		if deploymentStatus == e {
			return false
		}
	}
	return true
}

// NOTE: Refactor function to maybe make use of a type struct.
func (d *Deployment) getDeploymentStatus() string {
	pods := d.GetPods()
	execOutput, err := executeKubectlCmd(
		d.Namespace,
		"get",
		"pod",
		pods[len(pods)-1],
		"-o",
		"json",
	)
	if err != nil {
		bot.LogError(err)
		return ""
	}

	kubectlOutput := NewKubectlOutput(execOutput)
	state := kubectlOutput.Status.ContainerStatuses[0].State
	if state["waiting"] != nil {
		return state["waiting"]["reason"]
	} else if state["running"] != nil {
		return "Running"
	}

	return ""
}
