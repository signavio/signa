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

func NewPods(kubeconfig, namespace string) (p *Pods) {
	p = &Pods{}
	p.ParseFromKubectlOutput(kubeconfig, namespace)
	return
}

func (p *Pods) ParseFromKubectlOutput(kubeconfig, namespace string) error {
	execOutput, err := executeKubectlCmd(
		kubeconfig,
		namespace,
		"get",
		"pods",
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
