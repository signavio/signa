package main

import (
	"flag"
	"os"

	yaml "gopkg.in/yaml.v2"

	_ "github.com/signavio/signa/ext/kubernetes/deployment"
	_ "github.com/signavio/signa/ext/kubernetes/get"
	_ "github.com/signavio/signa/ext/kubernetes/info"
	_ "github.com/signavio/signa/ext/kubernetes/jobs"
	"github.com/signavio/signa/pkg/slack"
)

func main() {
	configFile := flag.String(
		"config", "/etc/signa.yaml", "Path to the configuration file.")
	flag.Parse()

	c := loadConfig(*configFile)
	slack.Run(*configFile, c["slack-token"].(string))
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
