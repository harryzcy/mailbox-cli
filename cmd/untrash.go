package cmd

import (
	"github.com/harryzcy/mailbox-cli/internal/command"
	"github.com/spf13/cobra"
)

var commandUntrash = command.Untrash

// untrashCmd represents the untrash command
var untrashCmd = &cobra.Command{
	Use:   "untrash",
	Short: "Untrash an email",
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) != 1 {
			cmd.PrintErrln("Please specify a messageID")
			osExit(1)
			return
		}
		messageID := args[0]

		verbose, _ := cmd.Flags().GetBool("verbose")
		result, err := commandUntrash(command.UntrashOptions{
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
	rootCmd.AddCommand(untrashCmd)
}
