package cmd

import (
	"os"
	"os/exec"
	"os/user"

	"github.com/pkg/errors"
)

var (
	DockerImage = "mmhealth"
)

// Responsible for passing any docker commands to the mmhealth container.
func runDockerCommand(cmdArgs []string) error {
	currentUser, err := user.Current()
	if err != nil {
		return errors.Wrap(err, "failed to get the current user")
	}

	pwd, err := os.Getwd()
	if err != nil {
		return errors.Wrap(err, "failed to get current directory")
	}

	dockerArgs := []string{
		"run",
		"--rm",
		"--platform=linux/amd64",
		"--volume",
		pwd + ":/files",
		"--user",
		currentUser.Uid + ":" + currentUser.Gid,
		DockerImage,
	}

	mergedArgs := append(dockerArgs, cmdArgs...)

	cmd := exec.Command("docker", mergedArgs...)
	stdout, _ := cmd.StdoutPipe()
	stderr, _ := cmd.StderrPipe()

	err = cmd.Start()
	if err != nil {
		return errors.Wrap(err, "failed to start the command")
	}

	go copyOutput(stdout)
	go copyOutput(stderr)

	err = cmd.Wait()
	if err != nil {
		return errors.Wrap(err, "failed to wait for the command to finish")
	}

	return nil
}
