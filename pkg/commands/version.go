package commands

import (
	"fmt"

	"github.com/sans-sroc/integrity/pkg/common"
	"github.com/urfave/cli/v2"
)

type versionCommand struct {
}

func (w *versionCommand) Execute(c *cli.Context) error {
	fmt.Printf("%s\n", common.AppVersion.Summary)

	return nil
}

func init() {
	cmd := versionCommand{}

	cliCmd := &cli.Command{
		Name:   "version",
		Usage:  "print version",
		Action: cmd.Execute,
	}

	common.RegisterCommand(cliCmd)
}
