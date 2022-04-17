package cmd

import (
	"fmt"
	"os"

	"github.com/harryzcy/mailbox-cli/internal/email"
	"github.com/spf13/cobra"
)

// listCmd represents the list command
var listCmd = &cobra.Command{
	Use:   "list",
	Short: "List emails",
	Run: func(cmd *cobra.Command, args []string) {
		verbose, _ := cmd.Flags().GetBool("verbose")
		client := email.Client{
			APIID:   cmd.Flag("api-id").Value.String(),
			Region:  cmd.Flag("region").Value.String(),
			Verbose: verbose,
		}
		result, err := client.List(email.ListOptions{
			Type:       cmd.Flag("type").Value.String(),
			Year:       cmd.Flag("year").Value.String(),
			Month:      cmd.Flag("month").Value.String(),
			Order:      cmd.Flag("order").Value.String(),
			NextCursor: cmd.Flag("next-cursor").Value.String(),
		})
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		fmt.Println(result)
	},
}

func init() {
	rootCmd.AddCommand(listCmd)
	listCmd.Flags().String("type", "", "Type")
	listCmd.Flags().String("year", "", "Year")
	listCmd.Flags().String("month", "", "Month")
	listCmd.Flags().String("order", "", "Order")
	listCmd.Flags().String("next-cursor", "", "Next Cursor")
}
