package cmd

import (
	"github.com/spf13/cobra"
)

// rootCmd represents the base command
var rootCmd = &cobra.Command{
	Use:   "SIP Parser",
	Short: "A CLI tool to process pcap files and folders",
	Long:  "A CLI tool with multiple commands to process pcap files and folders.",
}

// Execute runs the root command
func Execute() {
	cobra.CheckErr(rootCmd.Execute())
}

func init() {
	// Add subcommands to the root command
	rootCmd.AddCommand(loadCmd)
	rootCmd.AddCommand(getCostCmd)
	rootCmd.AddCommand(testCmd)
	rootCmd.AddCommand(filterCmd)
}
