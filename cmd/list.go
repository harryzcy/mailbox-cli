package cmd

import (
	"github.com/harryzcy/mailbox-cli/internal/command"
	"github.com/spf13/cobra"
)

var commandList = command.List

// listCmd represents the list command
var listCmd = &cobra.Command{
	Use:   "list",
	Short: "List emails",
	Args:  cobra.ExactArgs(0),
	Run: func(cmd *cobra.Command, _ []string) {
		verbose, err := cmd.Flags().GetBool("verbose")
		if err != nil {
			cmd.PrintErrln(err)
			osExit(1)
		}

		result, err := commandList(command.ListOptions{
			APIID:    cmd.Flag("api-id").Value.String(),
			Region:   cmd.Flag("region").Value.String(),
			Endpoint: cmd.Flag("endpoint").Value.String(),
			Verbose:  verbose,

			Type:       cmd.Flag("type").Value.String(),
			Year:       cmd.Flag("year").Value.String(),
			Month:      cmd.Flag("month").Value.String(),
			Order:      cmd.Flag("order").Value.String(),
			NextCursor: cmd.Flag("next-cursor").Value.String(),
		})
		if err != nil {
			cmd.PrintErrln(err)
			osExit(1)
			return
		}

		cmd.Println(result)
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
