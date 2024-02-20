package cmd

import (
	"errors"

	"github.com/spf13/cobra"
)

type DockerCompose struct {
	Version  string   `yaml:"version"`
	Services Services `yaml:"services"`
}

type Services struct {
	MMHealthCheck MMHealthCheck `yaml:"mmhealth"`
}

type MMHealthCheck struct {
	Image   string   `yaml:"image"`
	Volumes []string `yaml:"volumes"`
}

var InitCmd = &cobra.Command{
	Use:   "init",
	Short: "Inits everything needed for mmhealth",
	Long:  `Downloads the docker image needed`,
	RunE:  initCmdF,
}

func init() {
	RootCmd.AddCommand(InitCmd)
}

func initCmdF(cmd *cobra.Command, args []string) error {

	if BuildVersion == "(devel)" {
		return errors.New("Not a supported command in dev mode. Run `make buildDocker` instead")
	}
	return runCommand("docker", []string{"pull", DockerImage})
}
