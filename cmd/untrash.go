package cmd

import (
	"fmt"
	"os"

	"github.com/harryzcy/mailbox-cli/internal/email"
	"github.com/spf13/cobra"
)

// untrashCmd represents the untrash command
var untrashCmd = &cobra.Command{
	Use:   "untrash",
	Short: "Untrash an email",
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) != 1 {
			fmt.Println("Please specify a messageID")
			os.Exit(1)
		}
		messageID := args[0]

		verbose, _ := cmd.Flags().GetBool("verbose")
		client := email.Client{
			APIID:   cmd.Flag("api-id").Value.String(),
			Region:  cmd.Flag("region").Value.String(),
			Verbose: verbose,
		}
		result, err := client.Untrash(email.UntrashOptions{
			MessageID: messageID,
		})
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		fmt.Println(result)
	},
}

func init() {
	rootCmd.AddCommand(untrashCmd)
}
