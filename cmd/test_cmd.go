package cmd

import (
	"github.com/spf13/cobra"
	"sip-parser/pkg/process"
)

// getCostCmd represents the `get_cost` command
var testCmd = &cobra.Command{
	Use:   "test",
	Short: "test",
	Long:  "test",
	Run: func(cmd *cobra.Command, args []string) {
		process.TestFunc()
	},
}
