package cli

import (
	"fmt"

	"github.com/bitrise-tools/gows/config"
	"github.com/urfave/cli"
)

func listWorkspacesCmd(c *cli.Context) error {
	gowsConfig, err := config.LoadGOWSConfigFromFile()
	if err != nil {
		return fmt.Errorf("Failed to load gows config: %s", err)
	}

	fmt.Println()
	fmt.Println("=== Registered gows [project -> workspace] path list ===")
	for projectPath, wfConfig := range gowsConfig.Workspaces {
		fmt.Printf(" * %s -> %s\n", projectPath, wfConfig.WorkspaceRootPath)
	}
	fmt.Println("========================================================")
	fmt.Println()

	return nil
}
