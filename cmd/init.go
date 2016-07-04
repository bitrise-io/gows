package cmd

import (
	"errors"
	"fmt"

	log "github.com/Sirupsen/logrus"
	"github.com/bitrise-io/go-utils/colorstring"
	"github.com/bitrise-tools/gows/gows"
	"gopkg.in/viktorbenei/cobra.v0"
)

var (
	isAllowReset = false
)

// initCmd represents the init command
var initCmd = &cobra.Command{
	Use:           "init",
	Short:         "Initialize gows for your Go project",
	Long:          ``,
	SilenceUsage:  true,
	SilenceErrors: true,
	RunE: func(cmd *cobra.Command, args []string) error {
		if len(args) > 1 {
			return errors.New("More than one package argument specified")
		}
		packageName := ""
		if len(args) < 1 {
			log.Info("No package name specified, scanning it automatically ...")
			scanRes, err := gows.AutoScanPackageName()
			if err != nil {
				return fmt.Errorf("Failed to auto-scan the package name: %s", err)
			}
			if scanRes == "" {
				return errors.New("Empty package name scanned")
			}
			packageName = scanRes
			log.Infof(" Scanned package name: %s", packageName)
		} else {
			packageName = args[0]
		}

		if isAllowReset {
			log.Warning(colorstring.Red("Will reset the related workspace"))
		}

		if err := gows.Init(packageName, isAllowReset); err != nil {
			return fmt.Errorf("Failed to initialize: %s", err)
		}

		log.Info("Successful init - " + colorstring.Green("gows is ready for use!"))

		return nil
	},
}

func init() {
	RootCmd.AddCommand(initCmd)
	initCmd.Flags().BoolVarP(&isAllowReset,
		"reset", "",
		false,
		"Delete previous workspace (if any) and initialize a new one")
}
