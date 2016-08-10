package cmd

import (
	"errors"
	"fmt"
	"os"

	"github.com/bitrise-tools/gows/gows"
	"gopkg.in/viktorbenei/cobra.v0"
)

var cfgFile string
var isSyncBack bool

// RootCmd represents the base command when called without any subcommands
var RootCmd = &cobra.Command{
	Use:   "gows",
	Short: "Go Workspace / Environment Manager, to easily manage the Go Workspace during development.",
	Long: `Go Workspace / Environment Manager, to easily manage the Go Workspace during development.

Work in isolated (development) environment when you're working on your Go projects.
No cross-project dependency version missmatch, no more packages left out from vendor/.

No need for initializing a go workspace either, your project can be located anywhere,
not just in a predefined $GOPATH workspace. gows will take care about crearing
the (per-project isolated) workspace directory structure, no matter where your project is located.

gows works perfectly with other Go tools, all it does is it ensures that every project
gets it's own, isolated Go workspace and sets $GOPATH accordingly.`,

	DisableFlagParsing: true,
}

// Execute adds all child commands to the root command sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if err := RootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(-1)
	}
}

func init() {
	RootCmd.PersistentFlags().BoolVarP(&isSyncBack, "sync-back", "", false, "Sync back when command finishes")
	RootCmd.PreRunE = func(cmd *cobra.Command, args []string) error {
		if len(args) < 1 {
			return errors.New("No command specified!")
		}
		RootCmd.SilenceErrors = true
		RootCmd.SilenceUsage = true
		return nil
	}
	RootCmd.RunE = func(cmd *cobra.Command, args []string) error {
		cmdName := args[0]
		if cmdName == "-h" || cmdName == "--help" {
			if err := RootCmd.Help(); err != nil {
				return err
			}
			return nil
		}

		cmdArgs := []string{}
		if len(args) > 1 {
			cmdArgs = args[1:]
		}
		exitCode, err := gows.PrepareEnvironmentAndRunCommand(isSyncBack, cmdName, cmdArgs...)
		if exitCode != 0 {
			os.Exit(exitCode)
		}
		if err != nil {
			return fmt.Errorf("Exit Code was 0, but an error happened: %s", err)
		}
		return nil
	}
}
