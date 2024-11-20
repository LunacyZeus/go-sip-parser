package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
)

// getCostCmd represents the `get_cost` command
var testCmd = &cobra.Command{
	Use:   "test",
	Short: "test",
	Long:  "test",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("hello")
	},
}
