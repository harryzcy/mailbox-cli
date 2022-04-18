package cmd

import (
	"github.com/harryzcy/mailbox-cli/internal/command"
	"github.com/spf13/cobra"
)

var commandSave = command.Save

// saveCmd represents the save command
var saveCmd = &cobra.Command{
	Use:   "save messageID",
	Short: "Save a draft email",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		messageID := args[0]

		verbose, _ := cmd.Flags().GetBool("verbose")
		subject, _ := cmd.Flags().GetString("subject")
		from, _ := cmd.Flags().GetStringArray("from")
		to, _ := cmd.Flags().GetStringArray("to")
		cc, _ := cmd.Flags().GetStringArray("cc")
		bcc, _ := cmd.Flags().GetStringArray("bcc")
		replyTo, _ := cmd.Flags().GetStringArray("reply-to")
		body, _ := cmd.Flags().GetString("body")
		text, _ := cmd.Flags().GetString("text")
		html, _ := cmd.Flags().GetString("html")
		file, _ := cmd.Flags().GetString("file")

		result, err := commandSave(command.SaveOptions{
			APIID:    cmd.Flag("api-id").Value.String(),
			Region:   cmd.Flag("region").Value.String(),
			Endpoint: cmd.Flag("endpoint").Value.String(),
			Verbose:  verbose,

			MessageID: messageID,
			Subject:   subject,
			From:      from,
			To:        to,
			Cc:        cc,
			Bcc:       bcc,
			ReplyTo:   replyTo,
			Body:      body,
			Text:      text,
			HTML:      html,
			File:      file,
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
	rootCmd.AddCommand(saveCmd)
	saveCmd.Flags().String("subject", "", "Subject")
	saveCmd.Flags().StringArray("from", []string{}, "From")
	saveCmd.Flags().StringArray("to", []string{}, "To")
	saveCmd.Flags().StringArray("cc", []string{}, "Cc")
	saveCmd.Flags().StringArray("bcc", []string{}, "Bcc")
	saveCmd.Flags().StringArray("reply-to", []string{}, "Reply-To")
	saveCmd.Flags().String("body", "", "Body")
	saveCmd.Flags().String("text", "", "Text")
	saveCmd.Flags().String("html", "", "HTML")
	saveCmd.Flags().String("file", "", "File")
}
