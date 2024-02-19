package cmd

import (
	"os"

	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"gopkg.in/yaml.v2"
)

type DockerCompose struct {
	Version  string   `yaml:"version"`
	Services Services `yaml:"services"`
}

type Services struct {
	MMHealthCheck MMHealthCheck `yaml:"mm-healthcheck"`
}

type MMHealthCheck struct {
	Image   string   `yaml:"image"`
	Volumes []string `yaml:"volumes"`
}

var InitCmd = &cobra.Command{
	Use:   "init",
	Short: "Initialize a new healthcheck environment",
	Long:  `This command generates the docker-compose file for the healthcheck environment inside of the directory you're currently in.`,
	RunE:  initCmdF,
}

func init() {
	RootCmd.AddCommand(InitCmd)
}

func initCmdF(cmd *cobra.Command, args []string) error {
	dc := DockerCompose{
		Version: "3",
		Services: Services{
			MMHealthCheck: MMHealthCheck{
				Image:   "ghcr.io/coltoneshaw/mm-healthcheck:latest",
				Volumes: []string{".:/files"},
			},
		},
	}

	data, err := yaml.Marshal(&dc)
	if err != nil {
		return errors.Wrap(err, "failed to marshal docker-compose.yaml")
	}

	err = os.WriteFile("docker-compose.yaml", data, 0644)
	if err != nil {
		return errors.Wrap(err, "failed to write docker-compose.yaml")
	}
	return nil
}
