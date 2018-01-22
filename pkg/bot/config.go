package bot

import (
	"encoding/json"
	"os"
)

var cfg = &Config{}

type Config struct {
	BotUsername   string   `json:"bot-username"`
	SlackToken    string   `json:"slack-token"`
	Log           string   `json:"log"`
	RollbackCheck int      `json:"rollback-check"`
	Superusers    []string `json:"superusers"`
	Components    []Component
	Jobs          []Job
}

type Component struct {
	Name            string `json:"name"`
	Containers      []Container
	BootstrapConfig string   `json:"bootstrap-config"`
	Kubeconfig      string   `json:"kubeconfig"`
	Namespace       string   `json:"namespace"`
	ExecUsers       []string `json:"exec-users"`
	Alias           string   `json:"alias"`
}

type Job struct {
	Name       string   `json:"name"`
	Config     string   `json:"config"`
	Kubeconfig string   `json:"kubeconfig"`
	Namespace  string   `json:"namespace"`
	ExecUsers  []string `json:"exec-users"`
}

type Container struct {
	Name          string `json:"name"`
	RepositoryURI string `json:"repository-uri"`
}

func (c *Config) Load(file string) error {
	f, err := os.Open(file)
	defer f.Close()
	if err != nil {
		return err
	}

	parser := json.NewDecoder(f)
	if err = parser.Decode(&c); err != nil {
		return err
	}

	return nil
}

func (c *Config) AvailableComponents() (components []string) {
	for _, comp := range c.Components {
		components = append(components, comp.Name)
		if comp.Alias != "" {
			components = append(components, comp.Alias)
		}
	}
	return
}

func (c *Config) IsSuperuser(username string) bool {
	for _, su := range c.Superusers {
		if su == username {
			return true
		}
	}
	return false
}

func (c *Config) FindComponent(name string) *Component {
	for _, comp := range c.Components {
		if comp.Name == name || comp.Alias == name {
			return &comp
		}
	}
	return nil
}

func (c *Config) FindJob(name string) *Job {
	for _, j := range c.Jobs {
		if j.Name == name {
			return &j
		}
	}
	return nil
}

func (c *Component) FindContainer(name string) *Container {
	for _, cont := range c.Containers {
		if cont.Name == name {
			return &cont
		}
	}
	return nil
}

func (c *Component) IsExecUser(username string) bool {
	for _, u := range c.ExecUsers {
		if u == username {
			return true
		}
	}
	return false
}

func (j *Job) IsExecUser(username string) bool {
	for _, u := range j.ExecUsers {
		if u == username {
			return true
		}
	}
	return false
}
