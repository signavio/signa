package deployment

import (
	"bytes"
	"fmt"
	"os/exec"
	"strings"
	"text/template"

	"github.com/signavio/signa/pkg/bot"
)

const (
	commandName        = "deploy"
	commandDescription = "Deploy a service to production."
	commandUsage       = "%v <image-tag>"

	invalidAmountOfParams = "Invalid amount of parameters"
	invalidParams         = "Invalid parameters"
	componentNotFound     = "Component not found."
	containerNotFound     = "Container not found."
	clusterNotFound       = "Cluster not found."
	deployInfo            = "Deploying the image tag `%v` from `%v`. It may take several seconds."
	deployErrors          = "Some errors occurred :thinking_face:. Please, check the logs and try again in a few moments."
	rollbackExecuted      = "Problems identified during the deployment, the rollback was executed successfully."
	deploySuccess         = "The deployment was successful! Pods: `%v`."
	permissionDenied      = "You don't have enough permissions to execute this operation. :sweat_smile:"
	postDeploymentFailed  = "Something went wrong during the post deployment step. :morty:"
)

var messageChannel = make(chan string)

func init() {
	availableComponents := strings.Join(bot.Cfg().AvailableComponents(), "|")
	bot.RegisterCommand(
		commandName,
		commandDescription,
		fmt.Sprintf(commandUsage, availableComponents),
		Deploy,
	)
}

func Deploy(botCommand *bot.Cmd) (string, error) {
	if len(botCommand.Args) < 2 {
		return invalidAmountOfParams, nil
	}

	component := bot.Cfg().FindComponent(botCommand.Args[0])
	if component == nil {
		return componentNotFound, nil
	}
	container := component.FindContainer(botCommand.Args[1])
	if container == nil {
		return containerNotFound, nil
	}
	cluster := component.FindCluster(botCommand.Args[2])
	if cluster == nil {
		return clusterNotFound, nil
	}

	username := botCommand.User.Nick
	if bot.Cfg().IsSuperuser(username) || component.IsExecUser(username) {
		deployment := NewDeployment(component, container, cluster, botCommand.Args[3])
		err := postMessageToSlackChannel(
			botCommand.Channel,
			fmt.Sprintf(
				deployInfo,
				deployment.Version,
				deployment.Name,
			),
		)
		if err != nil {
			bot.LogError(err)
		}
		initiateDeploymentProcedure(deployment, component.Name, cluster.Name, username)
	} else {
		return permissionDenied, nil
	}

	// NOTE: This should be moved to the package bot.
	//logUserAction(
	//  botCommand.User.Nick,
	//  botCommand.Channel,
	//  botCommand.Command,
	//  botCommand.RawArgs,
	//)

	return <-messageChannel, nil
}

func initiateDeploymentProcedure(d *Deployment, componentName string, clusterName string, username string) {
	go func() {
		rollback, err := d.Apply()
		if err != nil {
			messageChannel <- deployErrors
		} else if rollback {
			messageChannel <- rollbackExecuted
		} else {
			messageChannel <- fmt.Sprintf(
				deploySuccess,
				strings.Join(d.GetPods(), " "),
			)
			initiatePostDeploymentStep(componentName, clusterName, username)
		}
	}()
}

var initiatePostDeploymentStep = func(componentName string, clusterName string, username string) {
	component := bot.Cfg().FindComponent(componentName)
	if component.HasPostDeploymentStep() && (clusterName == component.PostDeploymentStep.Cluster) {
		var err error
		command, err := addDeploymentInfoToCommand(component.PostDeploymentStep.Command, componentName, clusterName, username)
		if err != nil {
			messageChannel <- postDeploymentFailed
		}
		_, err = triggerRequest(command)
		if err != nil {
			messageChannel <- postDeploymentFailed
		}
	}
}

var addDeploymentInfoToCommand = func(command string, componentName string, clusterName string, username string) (string, error) {
	deploymentInfo := struct {
		Username      string
		ComponentName string
		ClusterName   string
	}{
		Username:      username,
		ComponentName: componentName,
		ClusterName:   clusterName,
	}

	t, _ := template.New("Command").Parse(command)
	var tpl bytes.Buffer
	if err := t.Execute(&tpl, deploymentInfo); err != nil {
		return "", err
	}
	result := tpl.String()
	return result, nil
}

var triggerRequest = func(request string) ([]byte, error) {
	cmd := exec.Command("/bin/sh", "-c", request)
	out, error := cmd.Output()
	return out, error
}
