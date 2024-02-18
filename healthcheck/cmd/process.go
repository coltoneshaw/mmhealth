package cmd

import (
	"archive/zip"
	"os"
	"path/filepath"

	generate "github.com/coltoneshaw/mm-healthcheck/healthcheck/generate"
	"github.com/spf13/cobra"
)

var ProcessCmd = &cobra.Command{
	Use:   "process",
	Short: "Process the support packet.",
	Long:  "Generates the output file from the support packet for any issues",
	RunE:  generateCmdF,
}

func init() {
	ProcessCmd.Flags().StringP("file", "f", "", "the support packet file to process")
	ProcessCmd.Flags().Bool("debug", true, "Whether to show debug logs or not")

	if err := ProcessCmd.MarkFlagRequired("file"); err != nil {
		panic(err)
	}

	RootCmd.AddCommand(
		ProcessCmd,
	)
}

func generateCmdF(cmd *cobra.Command, args []string) error {
	inputFilePath, _ := cmd.Flags().GetString("file")

	// input file
	packetReader, err := os.Open(filepath.Join("/files", inputFilePath))
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

	generate := generate.ProcessPacket{}

	err = generate.ProcessPacket(*packetConents)

	if err != nil {
		return err
	}

	return nil
}
