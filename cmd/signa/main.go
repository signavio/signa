package main

import (
	"os"

	yaml "gopkg.in/yaml.v2"

	_ "github.com/signavio/signa/ext/kubernetes/deployment"
	_ "github.com/signavio/signa/ext/kubernetes/get"
	_ "github.com/signavio/signa/ext/kubernetes/info"
	_ "github.com/signavio/signa/ext/kubernetes/jobs"
	"github.com/signavio/signa/pkg/slack"
)

func main() {
	// NOTE: Add the possibility of use a flag to load the conf
	// or this way as default.
	c := loadConfig("/etc/signa.yaml")
	slack.Run(c["slack-token"].(string))
}

func loadConfig(file string) map[string]interface{} {
	f, err := os.Open(file)
	defer f.Close()
	if err != nil {
		panic(err)
	}

	var c map[string]interface{}
	d := yaml.NewDecoder(f)
	d.Decode(&c)

	return c
}
