package exec

import (
	"strings"

	"github.com/signavio/signa/pkg/bot"
	"github.com/signavio/signa/pkg/gitlab"
)

type MergeRequest struct {
	Project  string
	Assignee string
	Source   string
	Target   string
}

func CreateGitMergeRequest(c *bot.Cmd) (string, error) {
	parsedArgs := parseArguments(c.Args)
	mergeRequest := NewMergeRequest(parsedArgs)
	// TODO: Implement config to setup api URI and token.
	mergeRequest.createOnGitlab("API URI FROM CONFIG", "TOKEN FROM CONFIG")
	return "", nil
}

func parseArguments(args []string) map[string]string {
	newArgs := make(map[string]string)
	for idx, a := range args {
		if isEven(idx) {
			key := strings.Trim(a, ":")
			value := args[idx+1]
			newArgs[key] = value
		}
	}
	return newArgs
}

func isEven(number int) bool {
	return number%2 == 0
}

func NewMergeRequest(args map[string]string) *MergeRequest {
	return &MergeRequest{
		Project:  args["project"],
		Assignee: args["assignee"],
		Source:   args["source"],
		Target:   args["target"],
	}
}

func (m *MergeRequest) createOnGitlab(apiURL, apiToken string) {
	gitlabConn, err := gitlab.NewConn(apiURL, apiToken)
	if err != nil {
		return
	}
	// TODO: Error handling.
	projectID, _ := gitlabConn.FindProjectID(m.Project)
	userID, _ := gitlabConn.FindUserID(m.Assignee)
	gitlabConn.PostMergeRequests(projectID, userID, m.Source, m.Target)
}

func init() {
	bot.RegisterCommand(
		"gmr",
		"Creates merge requests",
		"",
		CreateGitMergeRequest)
}
