package cmd

import (
	"github.com/harryzcy/mailbox-cli/internal/command"
	"github.com/spf13/cobra"
)

var commandDelete = command.Delete

// deleteCmd represents the delete command
var deleteCmd = &cobra.Command{
	Use:   "delete messageID",
	Short: "Delete an email",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		messageID := args[0]

		verbose, err := cmd.Flags().GetBool("verbose")
		if err != nil {
			cmd.PrintErrln(err)
			osExit(1)
		}

		result, err := commandDelete(command.DeleteOptions{
			APIID:    cmd.Flag("api-id").Value.String(),
			Region:   cmd.Flag("region").Value.String(),
			Endpoint: cmd.Flag("endpoint").Value.String(),
			Verbose:  verbose,

			MessageID: messageID,
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
	rootCmd.AddCommand(deleteCmd)
}
