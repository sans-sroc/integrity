package commands

import (
	"fmt"
	"github.com/sans-sroc/integrity/pkg/common"
	"github.com/sans-sroc/integrity/pkg/integrity"
	"github.com/sirupsen/logrus"
	"github.com/urfave/cli/v2"
)

type createCommand struct {
}

func (w *createCommand) Execute(c *cli.Context) error {
	if c.Args().Len() > 0 {
		return fmt.Errorf("positional arguments are not supported with this command.\n\n" + //nolint:stylecheck
			"Did you mean to use `-d` to change the directory that the command runs against?\n\n")
	}

	run, err := integrity.New(c.String("directory"), false)
	if err != nil {
		return err
	}

	if err := run.SetName(c.String("name")); err != nil {
		return err
	}

	run.SetFilename(c.String("filename"))
	run.SetIgnore(c.StringSlice("ignore"))

	if err := run.SetAlgorithm(c.String("algorithm")); err != nil {
		return err
	}

	if err := run.Checks(); err != nil {
		return err
	}

	if err := run.DiscoverFiles(); err != nil {
		return err
	}

	if err := run.HashFiles(); err != nil {
		return err
	}

	if err := run.WriteFile(); err != nil {
		return err
	}

	logrus.Info("Integrity file created successfully!")

	return nil
}

func init() {
	cmd := createCommand{}

	flags := []cli.Flag{
		&cli.StringFlag{
			Name:     "name",
			Usage:    fmt.Sprintf("The name that will be given to the ISO volume during USB creation. Format: %s", common.NameFormat),
			Aliases:  []string{"n"},
			EnvVars:  []string{"NAME"},
			Required: true,
		},
		&cli.StringFlag{
			Name:    "algorithm",
			Usage:   "Algorithm to use for hashing the files",
			Value:   "sha256",
			Aliases: []string{"a"},
			Hidden:  true,
		},
		&cli.StringSliceFlag{
			Name:    "ignore",
			Usage:   "Ignore files or directories as a direct match, prefix, or as a regular expressions",
			Aliases: []string{"i"},
			Hidden:  true,
			Value:   cli.NewStringSlice(common.IgnoreOnCreate...),
		},
		&cli.StringFlag{
			Name:    "user",
			Usage:   "DEPRECATED: no longer used -- allow setting what user created the file",
			Aliases: []string{"u"},
			EnvVars: []string{"USER"},
		},
	}

	cliCmd := &cli.Command{
		Name:   "create",
		Usage:  "create integrity files",
		Action: cmd.Execute,
		Flags:  append(flags, globalFlags()...),
		Before: globalBefore,
	}

	common.RegisterCommand(cliCmd)
}
