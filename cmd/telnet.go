package cmd

import (
	"github.com/spf13/cobra"
	"sip-parser/pkg/process"
)

var csvFilePath string

// telnetCmd represents the `load` command
var telnetCmd = &cobra.Command{
	Use:   "telnet",
	Short: "telnet",
	Long:  "telnet.",
	Run: func(cmd *cobra.Command, args []string) {
		process.StartTelnet(csvFilePath)
	},
}

func init() {
	telnetCmd.Flags().StringVarP(&csvFilePath, "file", "f", "", "Csv file")
	telnetCmd.MarkFlagRequired("file")
}
