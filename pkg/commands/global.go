package commands

import (
	"os/user"

	"github.com/sirupsen/logrus"
	"github.com/urfave/cli/v2"
)

func globalFlags() []cli.Flag {
	var username string
	u, err := user.Current()
	if err != nil {
		username = "unknown"
	}
	username = u.Username

	globalFlags := []cli.Flag{
		&cli.StringFlag{
			Name:    "log-level",
			Usage:   "Log Level",
			Aliases: []string{"l"},
			EnvVars: []string{"LOG_LEVEL"},
			Value:   "info",
		},
		&cli.StringFlag{
			Name:    "directory",
			Usage:   "Target Directory",
			Aliases: []string{"d"},
			EnvVars: []string{"DIRECTORY"},
			Value:   ".",
		},
		&cli.StringFlag{
			Name:     "courseware-version",
			Usage:    "Courseware Version Identifier",
			Aliases:  []string{"c"},
			EnvVars:  []string{"COURSEWARE_VERSION"},
			Required: true,
		},
		&cli.BoolFlag{
			Name:    "json",
			Usage:   "Output in JSON",
			Aliases: []string{"j"},
		},
		&cli.BoolFlag{
			Name:  "json-pretty",
			Usage: "Output JSON in Pretty Print Format",
			Value: true,
		},
		&cli.StringFlag{
			Name:  "user",
			Usage: "allow setting what user created the file",
			Value: username,
		},
	}

	return globalFlags
}

func globalBefore(c *cli.Context) error {
	switch c.String("log-level") {
	case "debug":
		logrus.SetLevel(logrus.DebugLevel)
	case "info":
		logrus.SetLevel(logrus.InfoLevel)
	case "warn":
		logrus.SetLevel(logrus.WarnLevel)
	case "error":
		logrus.SetLevel(logrus.ErrorLevel)
	}

	return nil
}
