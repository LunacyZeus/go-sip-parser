package cmd

import (
	"sip-parser/pkg/process"

	"github.com/spf13/cobra"
)

var convertCsvFilePath string

// getCostCmd represents the `get_cost` command
var convertCsvCmd = &cobra.Command{
	Use:   "convert_csv",
	Short: "convert_csv",
	Long:  "convert_csv",
	Run: func(cmd *cobra.Command, args []string) {
		process.ConvertCsv(convertCsvFilePath)
		//process.CalculateSipCostTest(costFilePath)
	},
}

func init() {
	convertCsvCmd.Flags().StringVarP(&convertCsvFilePath, "file", "f", "", "CSV File")
	convertCsvCmd.MarkFlagRequired("file")
}
