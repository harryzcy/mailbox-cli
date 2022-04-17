package cmd

import (
	"fmt"
	"os"

	"github.com/harryzcy/mailbox-cli/internal/email"
	"github.com/spf13/cobra"
)

// sendCmd represents the send command
var sendCmd = &cobra.Command{
	Use:   "send",
	Short: "Send an saved email",
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
		result, err := client.Send(email.SendOptions{
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
	rootCmd.AddCommand(sendCmd)
}
