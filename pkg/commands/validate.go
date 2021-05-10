package commands

import (
	"fmt"
	"path/filepath"

	"github.com/sans-sroc/integrity/pkg/common"
	"github.com/sans-sroc/integrity/pkg/utils"
	"github.com/urfave/cli/v2"
)

type validateCommand struct {
}

func (w *validateCommand) Execute(c *cli.Context) error {
	// Validate existing VERSION file(s)
	dir := filepath.ToSlash(c.String("directory"))
	ver := c.String("courseware-version")
	parts := c.Bool("parts")
	first := c.Bool("first")
	json := c.Bool("json")
	pretty := c.Bool("json-pretty")

	failed := utils.ValidateFiles(dir, ver, parts, first, json, pretty)
	if !json {
		if failed {
			fmt.Println("[!] Result: FAIL!")
		} else {
			fmt.Println("[+] Result: SUCCESS!")
		}
	}

	return nil
}

func init() {
	cmd := validateCommand{}

	flags := []cli.Flag{
		&cli.BoolFlag{
			Name:    "parts",
			Usage:   "Validate the VERSION-part.txt file",
			Aliases: []string{"p"},
		},
		&cli.BoolFlag{
			Name:    "first",
			Usage:   "Validate the VERSION-first.txt file",
			Aliases: []string{"f"},
		},
	}

	cliCmd := &cli.Command{
		Name:   "validate",
		Usage:  "validate integrity files",
		Action: cmd.Execute,
		Flags:  append(flags, globalFlags()...),
		Before: globalBefore,
	}

	common.RegisterCommand(cliCmd)
}
