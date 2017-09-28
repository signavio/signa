package gitlab

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
)

type GitlabConn struct {
	APIURL   string
	APIToken string
}

type Projects []*Project

type Project struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

const projectsSearchEndpoint = "/projects?search="

func NewConn(url, token string) (*GitlabConn, error) {
	if url == "" || token == "" {
		return &GitlabConn{}, errors.New(
			"Missing connection params. You must provide a valid attribute.",
		)
	}
	return &GitlabConn{APIURL: url, APIToken: token}, nil
}

// TODO: Implement calls.
func (c *GitlabConn) PostMergeRequests(projectID, userID int, sourceBranch, targetBranch string) (string, error) {
	// Post to the Merge Request API
	return "either response or only error", nil
}

func (c *GitlabConn) FindProjectID(projectName string) (int, error) {
	reqBuf, err := c.newHTTPRequestBuffer("GET", projectsSearchEndpoint+projectName)
	if err != nil {
		return 0, err
	}

	projects := Projects{}
	err = json.Unmarshal(reqBuf.Bytes(), &projects)
	if err != nil {
		return 0, err
	}

	for _, p := range projects {
		if p.Name == projectName {
			return p.ID, nil
		}
	}

	return 0, errors.New("Unable to find the project ID. Name must be valid.")
}

// TODO: Finish it. Implementation will be probably similar to FindProjectID()
func (c *GitlabConn) FindUserID(userName string) (int, error) {
	return 0, nil
}

func (c *GitlabConn) newHTTPRequestBuffer(method, endpoint string) (*bytes.Buffer, error) {
	httpClient := &http.Client{}

	req, err := http.NewRequest(
		method,
		c.APIURL+endpoint,
		nil,
	)
	if err != nil {
		return &bytes.Buffer{}, err
	}

	req.Header.Add("PRIVATE-TOKEN", c.APIToken)
	resp, err := httpClient.Do(req)
	if err != nil {
		return &bytes.Buffer{}, err
	}
	defer resp.Body.Close()

	buf := new(bytes.Buffer)
	buf.ReadFrom(resp.Body)

	return buf, nil
}
