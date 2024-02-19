package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var RootCmd = &cobra.Command{
	Use:   "healthcli",
	Short: "Mattermost health cli for managing health checks.",
}

func Execute() {
	RootCmd.CompletionOptions.HiddenDefaultCmd = true

	if err := RootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
