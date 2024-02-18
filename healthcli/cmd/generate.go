package cmd

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"

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
	inputFileName, err := cmd.Flags().GetString("inputFile")

	if err != nil {
		return errors.Wrap(err, "failed to get input file")
	}
	outputFileName, err := cmd.Flags().GetString("outputFile")

	if err != nil {
		return errors.Wrap(err, "failed to get output file")
	}

	inputFilePath := filepath.Join("/files", inputFileName)
	outputFilePath := filepath.Join("/files", outputFileName)

	if _, err := os.Stat(inputFilePath); os.IsNotExist(err) {
		return errors.Wrap(err, "report file does not exist")
	}

	cmdArgs := []string{
		"--template=template.tex",
		inputFilePath,
		"-o",
		outputFilePath,
	}

	pandoc := exec.Command("pandoc", cmdArgs...)
	pandoc.Stdout = os.Stdout
	pandoc.Stderr = os.Stderr
	err = pandoc.Run()
	if err != nil {
		return errors.Wrap(err, "failed to generate the pandoc report")
	}

	fmt.Println("PDF generated successfully.")
	return nil
}
