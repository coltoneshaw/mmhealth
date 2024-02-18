package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var RootCmd = &cobra.Command{
	Use:   "healthcheck",
	Short: "Mattermost healthcheck tool for parsing the support packets.",
}

func Execute() {
	RootCmd.CompletionOptions.HiddenDefaultCmd = true

	if err := RootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
