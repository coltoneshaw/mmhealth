package cmd

import (
	"archive/zip"
	"os"
	"path/filepath"

	generate "github.com/coltoneshaw/mm-healthcheck/healthcheck/generate"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"gopkg.in/yaml.v2"
)

var ProcessCmd = &cobra.Command{
	Use:   "generate",
	Short: "Generate the entire health report from the support packet.",
	Long:  "Generates the entire health report from the support packet, and outputting a pdf file.",
	RunE:  generateCmdF,
}

func init() {
	ProcessCmd.Flags().StringP("packet", "p", "", "the support packet file to process")
	ProcessCmd.Flags().StringP("output", "o", "healthcheck-report.pdf", "the output file name for the PDF.")

	ProcessCmd.Flags().Bool("debug", true, "Whether to show debug logs or not")

	if err := ProcessCmd.MarkFlagRequired("packet"); err != nil {
		panic(err)
	}

	RootCmd.AddCommand(
		ProcessCmd,
	)
}

func generateCmdF(cmd *cobra.Command, args []string) error {
	supportPacketFile, _ := cmd.Flags().GetString("packet")
	outputFileName, _ := cmd.Flags().GetString("output")

	outputFilePath := filepath.Join("/files", outputFileName)
	reportFilePath := filepath.Join("/files", "healthcheck-report.md")
	// input file
	packetReader, err := os.Open(filepath.Join("/files", supportPacketFile))
	if err != nil {
		return err
	}
	defer packetReader.Close()

	zipFileInfo, err := packetReader.Stat()
	if err != nil {
		return err
	}

	zipReader, err := zip.NewReader(packetReader, zipFileInfo.Size())
	if err != nil || zipReader.File == nil {
		return err
	}

	packetConents, err := generate.UnzipToMemory(zipReader)

	if err != nil {
		return err
	}

	g := generate.ProcessPacket{}

	report, err := g.ProcessPacket(*packetConents)
	if err != nil {
		return err
	}

	err = saveMarkdownReportToFile(reportFilePath, report)
	if err != nil {
		return err
	}

	err = generate.ReportToPDF(reportFilePath, outputFilePath)
	if err != nil {
		return err
	}

	err = deleteMarkdownReportFile(reportFilePath)
	if err != nil {
		return err
	}

	return nil
}

func saveMarkdownReportToFile(reportFilePath string, results generate.CheckResults) error {
	// Convert the results to YAML
	data, err := yaml.Marshal(results)
	if err != nil {
		return errors.Wrap(err, "failed to marshal results to yaml")
	}

	file, err := os.Create(reportFilePath)
	if err != nil {
		return errors.Wrap(err, "failed to create report.md")
	}
	defer file.Close()

	markdown := "---\n" + string(data) + "---\n"

	err = os.WriteFile(reportFilePath, []byte(markdown), 0644)
	if err != nil {
		return errors.Wrap(err, "failed to write to report.md")
	}
	return nil
}

func deleteMarkdownReportFile(reportFilePath string) error {
	err := os.Remove(reportFilePath)
	if err != nil {
		return errors.Wrap(err, "failed to delete report.md")
	}
	return nil
}
