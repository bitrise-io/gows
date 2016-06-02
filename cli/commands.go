package cli

import "github.com/urfave/cli"

const (
	// --- Only available through EnvVar flags

	// --- App flags

	// LogLevelEnvKey ...
	LogLevelEnvKey = "LOGLEVEL"
	// LogLevelKey ...
	LogLevelKey      = "loglevel"
	logLevelKeyShort = "l"

	// HelpKey ...
	HelpKey      = "help"
	helpKeyShort = "h"

	// VersionKey ...
	VersionKey      = "version"
	versionKeyShort = "v"

	// -- For "run"

	// SyncBackKey ...
	SyncBackKey = "sync-back"

	// --- Command flags

	// InitResetKey ...
	InitResetKey = "reset"
)

var (
	commands = []cli.Command{
		{
			Name:   "version",
			Usage:  "Version",
			Action: versionCmd,
		},
		{
			Name:   "init",
			Usage:  "Initialize gows for your Go project",
			Action: initCmd,
			Flags: []cli.Flag{
				cli.BoolFlag{
					Name:  InitResetKey,
					Usage: "Delete previous workspace (if any) and initialize a new one",
				},
			},
		},
	}

	appFlags = []cli.Flag{
		cli.StringFlag{
			Name:   LogLevelKey + ", " + logLevelKeyShort,
			Value:  "info",
			Usage:  "Log level (options: debug, info, warn, error, fatal, panic).",
			EnvVar: LogLevelEnvKey,
		},
		cli.BoolFlag{
			Name:  SyncBackKey,
			Usage: "Sync back when command finishes",
		},
	}
)

func init() {
	// Override default help and version flags
	cli.HelpFlag = cli.BoolFlag{
		Name:  HelpKey + ", " + helpKeyShort,
		Usage: "Show help.",
	}

	cli.VersionFlag = cli.BoolFlag{
		Name:  VersionKey + ", " + versionKeyShort,
		Usage: "Print the version.",
	}
}
