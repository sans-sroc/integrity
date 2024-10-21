package commands

import (
	"io"
	"os"
	"runtime"

	"github.com/shiena/ansicolor"
	"github.com/sirupsen/logrus"
	"github.com/urfave/cli/v2"

	"github.com/sans-sroc/integrity/pkg/common"
)

func globalFlags() []cli.Flag {
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
			Usage:   "The directory that will be the current working directory for the tool when it runs",
			Aliases: []string{"d"},
			EnvVars: []string{"DIRECTORY"},
			Value:   ".",
		},
		&cli.StringFlag{
			Name:   "filename",
			Usage:  "The integrity file",
			Hidden: true,
			Value:  common.Filename,
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
	case "none":
		logrus.SetOutput(io.Discard)
	}

	if runtime.GOOS == "windows" {
		logrus.SetFormatter(&logrus.TextFormatter{ForceColors: true})
		logrus.SetOutput(ansicolor.NewAnsiColorWriter(os.Stdout))
	}

	if c.Bool("json") {
		logrus.SetOutput(os.Stderr)

		if runtime.GOOS == "windows" {
			logrus.SetOutput(ansicolor.NewAnsiColorWriter(os.Stderr))
		}
	}

	return nil
}
