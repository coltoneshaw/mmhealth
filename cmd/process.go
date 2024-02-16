package cmd

import (
	"archive/zip"
	"fmt"
	"os"

	processpacket "github.com/coltoneshaw/healthcheck/processPacket"
	"github.com/spf13/cobra"
)

var ProcessCmd = &cobra.Command{
	Use:   "process",
	Short: "Process the support packet.",
	Long:  "Generates the output file from the support packet for any issues",
	RunE:  processPacketCmdF,
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

func processPacketCmdF(cmd *cobra.Command, args []string) error {
	inputFilePath, _ := cmd.Flags().GetString("file")

	fmt.Println("Processing the support packet: ", inputFilePath)

	// input file
	packetReader, err := os.Open(inputFilePath)
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

	packetConents, err := processpacket.UnzipToMemory(zipReader)

	if err != nil {
		return err
	}

	processpacket := processpacket.ProcessPacket{}

	err = processpacket.ProcessPacket(*packetConents)

	if err != nil {
		return err
	}

	return nil
}
