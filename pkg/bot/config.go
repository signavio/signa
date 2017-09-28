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
}

type Component struct {
	Name            string   `json:"name"`
	RepositoryURI   string   `json:"repository-uri"`
	BootstrapConfig string   `json:"bootstrap-config"`
	Namespace       string   `json:"namespace"`
	ExecUsers       []string `json:"exec-users"`
	Alias           string   `json:"alias"`
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

func (c *Component) IsExecUser(username string) bool {
	for _, u := range c.ExecUsers {
		if u == username {
			return true
		}
	}
	return false
}
