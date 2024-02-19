package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var (
	BuildHash = "dev mode"
	Version   = "0.1.0"
)

var VersionCmd = &cobra.Command{
	Use:   "version",
	Short: "Prints the version of healthcheck.",
	Args:  cobra.NoArgs,
	Run:   versionCmdF,
}

func init() {
	RootCmd.AddCommand(VersionCmd)
}

// TODO - have this pull a version from the env var
func versionCmdF(cmd *cobra.Command, args []string) {
	fmt.Println("healthcheck " + Version + " -- " + BuildHash)
}
