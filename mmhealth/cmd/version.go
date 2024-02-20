package cmd

import (
	"fmt"
	"runtime/debug"

	"github.com/spf13/cobra"
)

var GitCommit string
var GitVersion string

var BuildCommit = func() string {
	if GitCommit != "" {
		return GitCommit
	}

	if info, ok := debug.ReadBuildInfo(); ok {
		for _, setting := range info.Settings {
			if setting.Key == "vcs.revision" {
				return setting.Value
			}
		}
	}

	return ""
}()

var BuildVersion = func() string {
	if GitVersion != "" {
		return GitVersion
	}

	if info, ok := debug.ReadBuildInfo(); ok {
		return info.Main.Version
	}
	return ""
}()

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

	fmt.Println("mmhealth " + BuildVersion + " -- " + BuildCommit)

	return nil
}
