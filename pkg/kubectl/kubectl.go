package kubectl

import "os/exec"

type Kubectl struct {
	bin  string
	args []string
}

func NewKubectl(bin string, args []string) (kubectl *Kubectl, err error) {
	if bin == "default" {
		bin, err = WhereIs()
	}
	kubectl = &Kubectl{bin: bin, args: args}
	return
}

func WhereIs() (kubectlBinPath string, err error) {
	kubectlBinPath, err = exec.LookPath("kubectl")
	return
}

func (k *Kubectl) Exec() (string, error) {
	command := exec.Command(k.bin, k.args...)
	commandOutput, err := command.CombinedOutput()
	return string(commandOutput), err
}
