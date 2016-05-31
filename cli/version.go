package cli

import (
	"fmt"

	"github.com/bitrise-tools/gows/version"
	"github.com/urfave/cli"
)

func versionCmd(c *cli.Context) error {
	fmt.Printf("%s\n", version.VERSION)
	return nil
}
