package cmd

import (
	"os"
	"os/exec"
	"os/user"

	"github.com/pkg/errors"
)

var DockerImageProd = "ghcr.io/coltoneshaw/mmhealth"
var DockerImageDev = "mmhealth"

var DockerImage = func() string {
	if BuildVersion == "(devel)" {
		return DockerImageDev
	}
	return DockerImageProd + ":" + BuildVersion
}()

// Responsible for passing any docker commands to the mmhealth container.
func runDockerCommand(cmdArgs []string, additionalDockerArgs []string) error {
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
		"--volume", pwd + ":/files",
		"--user", currentUser.Uid + ":" + currentUser.Gid,
	}

	if len(additionalDockerArgs) > 0 {
		// adding the docker args right here to keep them before the image but not conflicting
		// mainy used for volumes
		dockerArgs = append(dockerArgs, additionalDockerArgs...)
	}

	dockerArgs = append(dockerArgs, DockerImage)
	mergedArgs := append(dockerArgs, cmdArgs...)

	return runCommand("docker", mergedArgs)
}

func runCommand(name string, args []string) error {
	cmd := exec.Command(name, args...)
	stdout, _ := cmd.StdoutPipe()
	stderr, _ := cmd.StderrPipe()

	err := cmd.Start()
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
