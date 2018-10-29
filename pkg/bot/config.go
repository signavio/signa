package bot

import (
	"os"

	yaml "gopkg.in/yaml.v2"
)

var cfg = &Config{}

type Config struct {
	BotUsername   string   `yaml:"bot-username"`
	SlackToken    string   `yaml:"slack-token"`
	Log           string   `yaml:"log"`
	RollbackCheck int      `yaml:"rollback-check"`
	Superusers    []string `yaml:"superusers"`
	Components    []Component
	Jobs          []Job
}

type Component struct {
	Name               string `yaml:"name"`
	Clusters           []Cluster
	Containers         []Container
	BootstrapConfig    string             `yaml:"bootstrap-config"`
	Kubeconfig         string             `yaml:"kubeconfig"`
	Namespace          string             `yaml:"namespace"`
	ExecUsers          []string           `yaml:"exec-users"`
	Alias              string             `yaml:"alias"`
	PostProductionStep PostProductionStep `yaml:"post-production-step"`
}

type Job struct {
	Name       string `yaml:"name"`
	Clusters   []Cluster
	Config     string   `yaml:"config"`
	Kubeconfig string   `yaml:"kubeconfig"`
	Namespace  string   `yaml:"namespace"`
	ExecUsers  []string `yaml:"exec-users"`
}

type Container struct {
	Name          string `yaml:"name"`
	RepositoryURI string `yaml:"repository-uri"`
}

type Cluster struct {
	Name       string `yaml:"name"`
	Kubeconfig string `yaml:"kubeconfig"`
}

type PostProductionStep struct {
	Command string `yaml:"command"`
	Cluster string `yaml:"cluster"`
}

func (c *Component) HasPostProductionStep() bool {
	if c.PostProductionStep.Cluster != "" {
		return true
	}
	return false
}

func (c *Config) Load(file string) error {
	f, err := os.Open(file)
	defer f.Close()
	if err != nil {
		return err
	}

	parser := yaml.NewDecoder(f)
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

func (c *Component) FindCluster(name string) *Cluster {
	for _, cluster := range c.Clusters {
		if cluster.Name == name {
			return &cluster
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

func (j *Job) FindCluster(name string) *Cluster {
	for _, cluster := range j.Clusters {
		if cluster.Name == name {
			return &cluster
		}
	}
	return nil
}
