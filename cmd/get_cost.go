package cmd

import (
	"sip-parser/pkg/process"

	"github.com/spf13/cobra"
)

var costFilePath string
var costThreads int

// getCostCmd represents the `get_cost` command
var getCostCmd = &cobra.Command{
	Use:   "get_cost",
	Short: "Calculates cost for file processing",
	Long:  "Calculates the cost of processing a  csv file with sip data",
	Run: func(cmd *cobra.Command, args []string) {
		process.CalculateSipCost(costFilePath, costThreads)
		//process.CalculateSipCostTest(costFilePath)
	},
}

func init() {
	getCostCmd.Flags().StringVarP(&costFilePath, "file", "f", "", "CSV File")
	getCostCmd.Flags().IntVarP(&costThreads, "threads", "t", 20, "telnet threads")
	getCostCmd.MarkFlagRequired("file")
}
