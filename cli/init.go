package cli

import (
	"errors"
	"fmt"

	log "github.com/Sirupsen/logrus"
	"github.com/bitrise-tools/gows/config"
	"github.com/bitrise-tools/gows/gows"
	"github.com/urfave/cli"
)

func initCmd(c *cli.Context) error {
	if c.NArg() > 1 {
		return errors.New("More than one package argument specified")
	}
	packageName := ""
	if c.NArg() < 1 {
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
		packageName = c.Args()[0]
	}

	if err := gows.Init(packageName); err != nil {
		return err
	}

	log.Info("Successful init - gows is ready for use!")
	log.Infof(" Note: you should add %s to your .gitignore", config.WorkspaceConfigFilePath)

	return nil
}
