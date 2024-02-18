package processpacket

import (
	"fmt"
	"os"
	"os/exec"

	"github.com/pkg/errors"
)

func ReportToPDF(inputFilePath, outputFilePath string) error {
	cmdArgs := []string{
		"--template=template.tex",
		inputFilePath,
		"-o",
		outputFilePath,
	}

	pandoc := exec.Command("pandoc", cmdArgs...)
	pandoc.Stdout = os.Stdout
	pandoc.Stderr = os.Stderr
	err := pandoc.Run()
	if err != nil {
		return errors.Wrap(err, "failed to generate the pandoc report")
	}

	fmt.Println("PDF generated successfully.")
	return nil
}
