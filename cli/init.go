package cli

import (
	"errors"
	"fmt"

	log "github.com/Sirupsen/logrus"
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
	} else {
		packageName = c.Args()[0]
	}

	return gows.Init(packageName)
}
