package commands

import (
	"errors"
	"fmt"
	"os"
	"strings"

	"github.com/sans-sroc/integrity/pkg/common"
	"github.com/sans-sroc/integrity/pkg/integrity"
	"github.com/sirupsen/logrus"
	"github.com/urfave/cli/v2"
)

type validateCommand struct {
}

func (w *validateCommand) Execute(c *cli.Context) error {
	if c.Args().Len() > 0 {
		return fmt.Errorf("positional arguments are not supported with this command.\n\n" + //nolint:stylecheck
			"did you mean to use `-d` to change the directory that the command runs against?\n\n")
	}

	run, err := integrity.New(c.String("directory"), true)
	if err != nil {
		return err
	}

	if _, err := os.Stat(c.String("filename")); err != nil && strings.Contains(err.Error(), "no such file") {
		return errors.New("the sans-integrity.yml checksum file does not exist")
	}

	run.SetFilename(c.String("filename"))
	run.SetIgnore(common.IgnoreAlways)

	if err := run.Checks(); err != nil {
		return err
	}

	if err := run.DiscoverFiles(); err != nil {
		return err
	}

	if err := run.HashFiles(); err != nil {
		return err
	}

	identical, err := run.CompareFiles()
	if err != nil {
		return err
	}

	if identical {
		logrus.Info("Success! All files successfully validated")
	}

	if c.String("output-format") == "json" {
		b, err := run.GetValidationOutput("json")
		if err != nil {
			return err
		}

		if c.String("output") == "-" {
			_, _ = os.Stdout.Write(b)
			_, _ = os.Stdout.WriteString("\n")
		} else {
			if err := os.WriteFile(c.String("output"), b, 0600); err != nil {
				return err
			}
		}
	}

	if !identical {
		return fmt.Errorf("validation Failed")
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
