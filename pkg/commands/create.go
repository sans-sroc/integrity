package commands

import (
	"os/user"

	"github.com/sans-sroc/integrity/pkg/common"
	"github.com/sans-sroc/integrity/pkg/integrity"
	"github.com/urfave/cli/v2"
)

type createCommand struct {
}

func (w *createCommand) Execute(c *cli.Context) error {
	dir := c.String("directory")
	//ver := c.String("courseware-version")
	//jsonOut := c.Bool("json")
	//pretty := c.Bool("json-pretty")
	//user := c.String("user")

	integrity, err := integrity.New(dir, false)
	if err != nil {
		return err
	}

	if err := integrity.SetName(c.String("name")); err != nil {
		return err
	}

	integrity.SetUser(c.String("user"))
	integrity.SetIgnore(c.StringSlice("ignore"))
	integrity.SetAlgorithm(c.String("algorithm"))

	if err := integrity.Checks(); err != nil {
		return err
	}

	if err := integrity.DiscoverFiles(); err != nil {
		return err
	}

	if err := integrity.HashFiles(); err != nil {
		return err
	}

	if err := integrity.WriteFile(); err != nil {
		return err
	}

	return nil
}

func init() {
	cmd := createCommand{}

	var username string
	u, err := user.Current()
	if err != nil {
		username = "unknown"
	}
	username = u.Username

	flags := []cli.Flag{
		&cli.StringFlag{
			Name:     "name",
			Usage:    "The name that will be given to the ISO volume during USB creation.",
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
			Value:   cli.NewStringSlice(common.IgnoreFiles...),
		},
		&cli.StringFlag{
			Name:    "user",
			Usage:   "allow setting what user created the file",
			Value:   username,
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
