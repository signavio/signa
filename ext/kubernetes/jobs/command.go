// TODO: Move it to the top package, offering the possibility of reuse.
package jobs

import "github.com/signavio/signa/pkg/kubectl"

type Command struct {
	Arguments []string
}

func NewCommand(args []string) *Command {
	return &Command{args}
}

func (c *Command) Exec() (string, error) {
	k, err := kubectl.NewKubectl("default", c.Arguments)
	if err != nil {
		return "", err
	}

	output, err := k.Exec()
	if err != nil {
		return "", err
	}

	return output, nil
}
