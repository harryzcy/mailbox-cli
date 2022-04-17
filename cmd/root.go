package cmd

import (
	"os"

	"github.com/spf13/cobra"
)

var osExit = os.Exit

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "mailbox-cli",
	Short: "Handle mailbox APIs from the command line.",
	Run: func(cmd *cobra.Command, args []string) {
		cmd.Help()
	},
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		osExit(1)
	}
}

func init() {
	rootCmd.PersistentFlags().String("api-id", "", "API ID")
	rootCmd.PersistentFlags().String("region", "", "Region")
	rootCmd.PersistentFlags().String("endpoint", "", "Endpoint")
	rootCmd.PersistentFlags().BoolP("verbose", "v", false, "Verbose mode")
}
