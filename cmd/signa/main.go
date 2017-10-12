package main

import (
	"encoding/json"
	"os"

	_ "github.com/signavio/signa/ext/kubernetes/deployment"
	_ "github.com/signavio/signa/ext/kubernetes/get"
	_ "github.com/signavio/signa/ext/kubernetes/info"
	"github.com/signavio/signa/pkg/slack"
)

func loadConfig(file string) map[string]interface{} {
	f, err := os.Open(file)
	defer f.Close()
	if err != nil {
		panic(err)
	}

	var c map[string]interface{}
	d := json.NewDecoder(f)
	d.Decode(&c)

	return c
}

func main() {
	// NOTE: Add the possibility of use a flag to load the conf
	// or this way as default.
	c := loadConfig("/etc/signa.conf")
	slack.Run(c["slack-token"].(string))
}
