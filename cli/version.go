package cli

import (
	"fmt"

	"github.com/bitrise-tools/gows/version"
	"github.com/codegangsta/cli"
)

func versionCmd(c *cli.Context) error {
	fmt.Printf("%s\n", version.VERSION)
	return nil
}
