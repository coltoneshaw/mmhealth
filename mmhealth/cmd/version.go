package cmd

import (
	"fmt"

	mmhealth "github.com/coltoneshaw/mmhealth/mmhealth"
	"github.com/spf13/cobra"
)

var VersionCmd = &cobra.Command{
	Use:   "version",
	Short: "Prints the version of healthcheck.",
	Args:  cobra.NoArgs,
	RunE:  versionCmdF,
}

func init() {
	RootCmd.AddCommand(VersionCmd)
}

// TODO - have this pull a version from the env var
func versionCmdF(cmd *cobra.Command, args []string) error {

	fmt.Println("mmhealth " + mmhealth.BuildVersion + " -- " + mmhealth.BuildCommit)

	return nil
}
