package cmd

import (
	"sip-parser/pkg/process"

	"github.com/spf13/cobra"
)

var reverseCsvFilePath string

// getCostCmd represents the `get_cost` command
var reverseCsvCmd = &cobra.Command{
	Use:   "reverse_csv",
	Short: "reverse_csv",
	Long:  "reverse_csv",
	Run: func(cmd *cobra.Command, args []string) {
		process.ReverseCsv(reverseCsvFilePath)
		//process.CalculateSipCostTest(costFilePath)
	},
}

func init() {
	reverseCsvCmd.Flags().StringVarP(&reverseCsvFilePath, "file", "f", "", "CSV File")
	reverseCsvCmd.MarkFlagRequired("file")
}
