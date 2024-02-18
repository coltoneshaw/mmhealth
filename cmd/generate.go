package cmd

import (
	"fmt"
	"os"
	"os/exec"
	"os/user"

	"github.com/pkg/errors"
	"github.com/spf13/cobra"
)

var PdfCmd = &cobra.Command{
	Use:   "pdf",
	Short: "Generate PDF",
	Long:  `This command generates a PDF using Docker.`,
	RunE:  generatePdfCmdF,
}

func init() {
	PdfCmd.Flags().StringP("inputFile", "i", "report.md", "The input file, which should be the result of ./healtcheck process")
	PdfCmd.Flags().StringP("outputFile", "o", "report.pdf", "The output file name for the PDF.")

	RootCmd.AddCommand(PdfCmd)
}

func generatePdfCmdF(cmd *cobra.Command, args []string) error {
	inputFilePath, _ := cmd.Flags().GetString("inputFile")
	outputFilePath, _ := cmd.Flags().GetString("outputFile")

	if _, err := os.Stat(inputFilePath); os.IsNotExist(err) {
		return errors.Wrap(err, "report file does not exist")
	}

	pwd, err := os.Getwd()
	if err != nil {
		return errors.Wrap(err, "failed to get current directory")
	}

	currentUser, err := user.Current()
	if err != nil {
		fmt.Println("Error:", err)
		return errors.Wrap(err, "failed to get current user")
	}

	cmdArgs := []string{
		"run",
		"--rm",
		"--platform=linux/amd64",
		"--volume",
		pwd + ":/data",
		"--user",
		currentUser.Uid + ":" + currentUser.Gid,
		"ghcr.io/coltoneshaw/mm-healthcheck:latest",
		"--template=/data/template/template.tex",
		inputFilePath,
		"-o",
		outputFilePath,
	}

	// used to debug
	// fmt.Println("Running command: docker", strings.Join(cmdArgs, " "))

	runDocker := exec.Command("docker", cmdArgs...)
	runDocker.Stdout = os.Stdout
	runDocker.Stderr = os.Stderr
	err = runDocker.Run()
	if err != nil {
		return errors.Wrap(err, "failed to generate report with docker container")
	}

	fmt.Println("PDF generated successfully.")
	return nil
}
