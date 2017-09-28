package deployment

import (
	"bytes"

	"github.com/signavio/signa/pkg/bot"
)

type Pods struct {
	Items []Pod
}

type Pod struct {
	Metadata struct {
		Name              string `json:"name"`
		CreationTimestamp string `json:"creationTimestamp"`
	} `json:"metadata"`
}

func NewPods(namespace string) (p *Pods) {
	p = &Pods{}
	p.ParseFromKubectlOutput(namespace)
	return
}

func (p *Pods) ParseFromKubectlOutput(namespace string) error {
	execOutput, err := executeKubectlCmd(
		namespace,
		"get",
		"pods",
		// NOTE: kubectl has a bug on this parameter. It's broken on the 1.7.x
		// and it should be fixed in the next minor version upgrade as the
		// issue got merged recently. 1.6.8 is fine though.
		// https://github.com/kubernetes/kubernetes/pull/48659
		"--sort-by=.metadata.creationTimestamp",
		"-o",
		"json",
	)
	if err != nil {
		bot.LogError(err)
		return err
	}

	if err = decodeJson(bytes.NewBufferString(execOutput), p); err != nil {
		bot.LogError(err)
		return err
	}

	return nil
}
