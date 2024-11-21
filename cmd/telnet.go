package cmd

import (
	"github.com/spf13/cobra"
	"sip-parser/pkg/process"
)

var cip string
var cport string
var ani string
var dnis string

// telnetCmd represents the `load` command
var telnetCmd = &cobra.Command{
	Use:   "telnet",
	Short: "telnet",
	Long:  "telnet.",
	Run: func(cmd *cobra.Command, args []string) {
		params := process.CallSimulationParams{
			CallerIp:   cip,
			CallerPort: cport,
			Ani:        ani,
			Dnis:       dnis,
		}
		process.StartTelnet(params)
	},
}

func init() {
	telnetCmd.Flags().StringVarP(&cip, "cip", "", "", "CallerIp")
	telnetCmd.Flags().StringVarP(&cport, "cport", "", "", "CallerPort")
	telnetCmd.Flags().StringVarP(&ani, "ani", "", "", "Ani")
	telnetCmd.Flags().StringVarP(&dnis, "dnis", "", "", "Dnis")

	telnetCmd.MarkFlagRequired("cip")
	telnetCmd.MarkFlagRequired("cport")
	telnetCmd.MarkFlagRequired("ani")
	telnetCmd.MarkFlagRequired("dnis")
}
