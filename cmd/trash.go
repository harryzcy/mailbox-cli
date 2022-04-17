package cmd

import (
	"github.com/harryzcy/mailbox-cli/internal/command"
	"github.com/spf13/cobra"
)

var commandTrash = command.Trash

// trashCmd represents the trash command
var trashCmd = &cobra.Command{
	Use:   "trash",
	Short: "Trash an email",
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) != 1 {
			cmd.PrintErrln("Please specify a messageID")
			osExit(1)
			return
		}
		messageID := args[0]

		verbose, _ := cmd.Flags().GetBool("verbose")
		result, err := commandTrash(command.TrashOptions{
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
	rootCmd.AddCommand(trashCmd)
}
