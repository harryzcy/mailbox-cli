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

		verbose, err := cmd.Flags().GetBool("verbose")
		if err != nil {
			cmd.PrintErrln(err)
			osExit(1)
		}

		subject, err := cmd.Flags().GetString("subject")
		if err != nil {
			cmd.PrintErrln(err)
			osExit(1)
		}
		from, err := cmd.Flags().GetStringArray("from")
		if err != nil {
			cmd.PrintErrln(err)
			osExit(1)
		}
		to, err := cmd.Flags().GetStringArray("to")
		if err != nil {
			cmd.PrintErrln(err)
			osExit(1)
		}
		cc, err := cmd.Flags().GetStringArray("cc")
		if err != nil {
			cmd.PrintErrln(err)
			osExit(1)
		}
		bcc, err := cmd.Flags().GetStringArray("bcc")
		if err != nil {
			cmd.PrintErrln(err)
			osExit(1)
		}
		replyTo, err := cmd.Flags().GetStringArray("reply-to")
		if err != nil {
			cmd.PrintErrln(err)
			osExit(1)
		}
		body, err := cmd.Flags().GetString("body")
		if err != nil {
			cmd.PrintErrln(err)
			osExit(1)
		}
		text, err := cmd.Flags().GetString("text")
		if err != nil {
			cmd.PrintErrln(err)
			osExit(1)
		}
		html, err := cmd.Flags().GetString("html")
		if err != nil {
			cmd.PrintErrln(err)
			osExit(1)
		}
		file, err := cmd.Flags().GetString("file")
		if err != nil {
			cmd.PrintErrln(err)
			osExit(1)
		}
		generateText, err := cmd.Flags().GetString("generate-text")
		if err != nil {
			cmd.PrintErrln(err)
			osExit(1)
		}
		send, err := cmd.Flags().GetBool("send")
		if err != nil {
			cmd.PrintErrln(err)
			osExit(1)
		}

		result, err := commandSave(command.SaveOptions{
			APIID:    cmd.Flag("api-id").Value.String(),
			Region:   cmd.Flag("region").Value.String(),
			Endpoint: cmd.Flag("endpoint").Value.String(),
			Verbose:  verbose,

			MessageID:    messageID,
			Subject:      subject,
			From:         from,
			To:           to,
			Cc:           cc,
			Bcc:          bcc,
			ReplyTo:      replyTo,
			Body:         body,
			Text:         text,
			HTML:         html,
			GenerateText: generateText,
			Send:         send,

			File: file,
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
	saveCmd.Flags().String("generate-text", "", "Generate text from HTML (optional)")
	saveCmd.Flags().String("file", "", "File")
	saveCmd.Flags().Bool("send", false, "Send email immediately without using draft (optional)")
}
