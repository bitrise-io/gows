package cli

import (
	"fmt"

	log "github.com/Sirupsen/logrus"
	"github.com/bitrise-io/go-utils/colorstring"
	"github.com/bitrise-tools/gows/config"
	"github.com/bitrise-tools/gows/gows"
	"github.com/urfave/cli"
)

func clearCmd(c *cli.Context) error {
	projectConfig, err := config.LoadProjectConfigFromFile()
	if err != nil {
		log.Info("Run " + colorstring.Green("gows init") + " to initialize a workspace & gows config for this project")
		return fmt.Errorf("Failed to read Project Config: %s", err)
	}
	if projectConfig.PackageName == "" {
		log.Info("Run " + colorstring.Green("gows init") + " to initialize a workspace & gows config for this project")
		return fmt.Errorf("Package Name is empty")
	}

	if err := gows.Init(projectConfig.PackageName, true); err != nil {
		return fmt.Errorf("Failed to initialize: %s", err)
	}

	log.Println("Done, workspace is clean!")

	return nil
}
