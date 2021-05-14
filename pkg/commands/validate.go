package commands

import (
	"fmt"
	"io/ioutil"
	"os"

	"github.com/sans-sroc/integrity/pkg/common"
	"github.com/sans-sroc/integrity/pkg/integrity"
	"github.com/sirupsen/logrus"
	"github.com/urfave/cli/v2"
)

type validateCommand struct {
}

func (w *validateCommand) Execute(c *cli.Context) error {
	if c.Args().Len() > 0 {
		return fmt.Errorf("Positional arguments are not supported with this command.\n\nDid you mean to use `-d` to change the directory that the command runs against?\n\n")
	}

	integrity, err := integrity.New(c.String("directory"), true)
	if err != nil {
		return err
	}

	integrity.SetFilename(c.String("filename"))
	integrity.SetIgnore(common.IgnoreAlways)

	if err := integrity.Checks(); err != nil {
		return err
	}

	if err := integrity.DiscoverFiles(); err != nil {
		return err
	}

	if err := integrity.HashFiles(); err != nil {
		return err
	}

	identical, err := integrity.CompareFiles()
	if err != nil {
		return err
	}

	if identical {
		logrus.Info("Success! All files successfully validated")
	}

	if c.String("output-format") == "json" {
		b, err := integrity.GetValidationOutput("json")
		if err != nil {
			return err
		}

		if c.String("output") == "-" {
			os.Stdout.Write(b)
			os.Stdout.Write([]byte("\n"))
		} else {
			ioutil.WriteFile(c.String("output"), b, 0644)
		}
	}

	if !identical {
		return fmt.Errorf("Validation Failed")
	}

	return nil
}

func init() {
	cmd := validateCommand{}

	flags := []cli.Flag{
		&cli.StringFlag{
			Name:    "output-format",
			Usage:   "Chose which format to output the validation results (default is none) (valid options: none, json)",
			Aliases: []string{"format"},
			EnvVars: []string{"OUTPUT_FORMAT"},
			Value:   "none",
		},
		&cli.StringFlag{
			Name:    "output",
			Usage:   "When output-format is specified, this controls where it goes, (defaults to stdout)",
			Aliases: []string{"o"},
			EnvVars: []string{"OUTPUT"},
			Value:   "-",
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
