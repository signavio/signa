package deployment

import (
	"strings"
	"testing"
)

func TestGetPostDeploymentStepCommandSuccessful(t *testing.T) {
	var commandTemplate = "curl -u componentname: {{.ComponentName}} clustername: {{.ClusterName}} username: {{.Username}} "
	var expected = "curl -u componentname: Gryffindor clustername: Hogwarts username: Harry Potter"
	command, error := addDeploymentInfoToCommand(commandTemplate, "Gryffindor", "Hogwarts", "Harry Potter")
	if strings.Compare(command, expected) == 0 {
		t.Fatal("Deployment info not correctly accessed. Expected " + expected + ", got: " + command)
	}
	if error != nil {
		t.Fatal("Error during post deployment info retrieval")
	}
}

func TestGetPostDeploymentStepCommandUnsuccessful(t *testing.T) {
	var commandTemplate = "curl"
	command, _ := addDeploymentInfoToCommand(commandTemplate, "Gryffindor", "Hogwarts", "Harry Potter")
	if command != commandTemplate {
		t.Fatal("Curl command not correctly formmatted. Expected " + commandTemplate + ", got: " + command)
	}
}
