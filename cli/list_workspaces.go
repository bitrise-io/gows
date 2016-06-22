package cli

import (
	"fmt"
	"os"

	log "github.com/Sirupsen/logrus"
	"github.com/bitrise-io/go-utils/colorstring"
	"github.com/bitrise-tools/gows/config"
	"github.com/urfave/cli"
)

func listWorkspacesCmd(c *cli.Context) error {
	gowsConfig, err := config.LoadGOWSConfigFromFile()
	if err != nil {
		return fmt.Errorf("Failed to load gows config: %s", err)
	}

	currWorkDir, err := os.Getwd()
	if err != nil {
		log.Debugf("Failed to get current working directory: %s", err)
	}

	fmt.Println()
	fmt.Println("=== Registered gows [project -> workspace] path list ===")
	for projectPath, wfConfig := range gowsConfig.Workspaces {
		if projectPath == currWorkDir {
			fmt.Println(colorstring.Greenf(" * %s -> %s", projectPath, wfConfig.WorkspaceRootPath))
		} else {
			fmt.Printf(" * %s -> %s\n", projectPath, wfConfig.WorkspaceRootPath)
		}
	}
	fmt.Println("========================================================")
	fmt.Println()

	return nil
}
