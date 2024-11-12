package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/sirupsen/logrus"
	"github.com/slawo/go-dediarize/dediarize"
	"github.com/urfave/cli/v2"
)

func main() {
	if err := runApp(); err != nil {
		logrus.Fatalf("failed with error: %s", err.Error())
	}
}

func runApp() error {

	exe, err := os.Executable()
	if err != nil {
		return err
	}

	return (&cli.App{
		Name:                 exe,
		EnableBashCompletion: true,

		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:    "file",
				Usage:   "the file to convert",
				EnvVars: []string{"FILE"},
			},
			&cli.StringFlag{
				Name:    "output",
				Aliases: []string{"out", "o"},
				Usage:   "the file to convert",
				EnvVars: []string{"OUTPUT"},
			},
			&cli.StringFlag{
				Name:    "loglevel",
				Aliases: []string{"ll", "level"},
				Value:   logrus.WarnLevel.String(),
				Usage: fmt.Sprintf("log level (%s)",
					strings.Join(getAllLogLevels(), "|"),
				),
				EnvVars: []string{"LOGLEVEL"},
			},
		},
		Action: RunDediarize,
	}).Run(os.Args)
}

func getAllLogLevels() []string {
	levels := make([]string, len(logrus.AllLevels))
	for i, ll := range logrus.AllLevels {
		levels[i] = ll.String()
	}
	return levels
}

func RunDediarize(cCtx *cli.Context) error {
	ll, err := logrus.ParseLevel(cCtx.String("loglevel"))
	if err != nil {
		return err
	}
	logrus.SetLevel(ll)

	if cCtx.String("file") == "" {
		return fmt.Errorf("file is required")
	}

	if cCtx.String("output") == "" {
		return fmt.Errorf("output is required")
	}

	err = dediarize.TranscribeJsonFile(cCtx.String("file"), cCtx.String("output"))
	if err != nil {
		return fmt.Errorf("failed to load json: %w", err)
	}
	return nil
}
