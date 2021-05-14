package commands

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/sans-sroc/integrity/pkg/common"
	"github.com/sans-sroc/integrity/pkg/integrity"
	"github.com/sirupsen/logrus"
	"github.com/urfave/cli/v2"
)

type validateCommand struct {
}

func (w *validateCommand) Execute(c *cli.Context) error {
	// Validate existing VERSION file(s)
	dir := filepath.ToSlash(c.String("directory"))
	//ver := c.String("courseware-version")
	//parts := c.Bool("parts")
	//first := c.Bool("first")
	//json := c.Bool("json")
	//pretty := c.Bool("json-pretty")

	integrity, err := integrity.New(dir, true)
	if err != nil {
		return err
	}

	integrity.SetIgnore(common.IgnoreAlways)

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
		logrus.Info("Success - all files validate")
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
