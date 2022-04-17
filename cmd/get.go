package cmd

import (
	"fmt"
	"os"

	"github.com/harryzcy/mailbox-cli/internal/email"
	"github.com/spf13/cobra"
)

// getCmd represents the get command
var getCmd = &cobra.Command{
	Use:   "get",
	Short: "Get an email by messageID",
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) != 1 {
			fmt.Println("Please specify a messageID")
			os.Exit(1)
		}
		messageID := args[0]

		client := email.Client{
			APIID:  cmd.Flag("api-id").Value.String(),
			Region: cmd.Flag("region").Value.String(),
		}
		result, err := client.Get(email.GetOptions{
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
	rootCmd.AddCommand(getCmd)
}
