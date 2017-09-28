package deployment

import "bytes"

type KubectlOutput struct {
	Status struct {
		ContainerStatuses []ContainerStatus
	} `json:"status"`
}

type ContainerStatus struct {
	Name  string                       `json:"name"`
	State map[string]map[string]string `json:"state"`
}

func NewKubectlOutput(output string) (ko *KubectlOutput) {
	ko = &KubectlOutput{}
	decodeJson(bytes.NewBufferString(output), ko)
	return
}
