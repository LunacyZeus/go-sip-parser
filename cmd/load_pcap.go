package cmd

import (
	"github.com/spf13/cobra"
	"sip-parser/pkg/process"
)

var loadFilePath string

// loadCmd represents the `load` command
var loadCmd = &cobra.Command{
	Use:   "load",
	Short: "Loads Pcap files or folders",
	Long:  "Loads a Pcap file or processes all .pcap files in a folder.",
	Run: func(cmd *cobra.Command, args []string) {
		process.LoadPcap(loadFilePath)
	},
}

func init() {
	loadCmd.Flags().StringVarP(&loadFilePath, "file", "f", "", "Pcap File or folder path to process (required)")
	loadCmd.MarkFlagRequired("file")
}
