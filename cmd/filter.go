package cmd

import (
	"github.com/spf13/cobra"
	"sip-parser/pkg/process"
)

var filterFilePath string
var callId string

// loadCmd represents the `load` command
var filterCmd = &cobra.Command{
	Use:   "filter",
	Short: "Loads filter Pcap files or folders",
	Long:  "Loads filter a Pcap file or processes all .pcap files in a folder.",
	Run: func(cmd *cobra.Command, args []string) {
		process.FilterPcapFileOrFolder(filterFilePath, callId)
	},
}

func init() {
	filterCmd.Flags().StringVarP(&filterFilePath, "file", "f", "", "Pcap File or folder path to process (required)")
	filterCmd.Flags().StringVarP(&callId, "callId", "c", "", "callId")
	filterCmd.MarkFlagRequired("file")
}
