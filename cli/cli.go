package cli

import (
	"fmt"
	"os"
	"path"

	log "github.com/Sirupsen/logrus"
	"github.com/bitrise-tools/gows/gows"
	"github.com/bitrise-tools/gows/version"
	"github.com/urfave/cli"
)

func before(c *cli.Context) error {
	// Log level
	if logLevel, err := log.ParseLevel(c.String(LogLevelKey)); err != nil {
		log.Fatal("Failed to parse log level:", err)
	} else {
		log.SetLevel(logLevel)
	}

	return nil
}

func printVersion(c *cli.Context) {
	fmt.Fprintf(c.App.Writer, "%v\n", c.App.Version)
}

// Run CLI
func Run() {
	cli.VersionPrinter = printVersion

	app := cli.NewApp()
	app.Name = path.Base(os.Args[0])
	app.Usage = "gows"
	app.Version = version.VERSION

	app.Author = ""
	app.Email = ""

	app.Before = before

	app.Flags = appFlags
	app.Commands = commands

	app.Action = func(c *cli.Context) error {
		cmdName := c.Args()[0]
		cmdArgs := []string{}
		if len(c.Args()) > 1 {
			cmdArgs = c.Args()[1:]
		}
		exitCode, err := gows.PrepareEnvironmentAndRunCommand(cmdName, cmdArgs...)
		if exitCode != 0 {
			return cli.NewExitError("", exitCode)
		}
		if err != nil {
			return fmt.Errorf("Exit Code was 0, but an error happened: %s", err)
		}
		return nil
	}

	if err := app.Run(os.Args); err != nil {
		log.Fatalf("Error: %s", err)
	}
}
