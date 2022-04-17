package cmd

import (
	"fmt"

	"github.com/harryzcy/mailbox-cli/internal/command"
	"github.com/spf13/cobra"
)

var commandGet = command.Get

// getCmd represents the get command
var getCmd = &cobra.Command{
	Use:   "get",
	Short: "Get an email by messageID",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println(len(args))
		if len(args) != 1 {
			cmd.PrintErrln("Please specify a messageID")
			osExit(1)
			return
		}
		messageID := args[0]

		verbose, _ := cmd.Flags().GetBool("verbose")
		result, err := commandGet(command.GetOptions{
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
	rootCmd.AddCommand(getCmd)
}
