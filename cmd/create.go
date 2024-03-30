package cmd

import (
	"github.com/harryzcy/mailbox-cli/internal/command"
	"github.com/spf13/cobra"
)

var commandCreate = command.Create

// createCmd represents the create command
var createCmd = &cobra.Command{
	Use:   "create",
	Short: "Create an email",
	Args:  cobra.ExactArgs(0),
	Run: func(cmd *cobra.Command, _ []string) {
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

		result, err := commandCreate(command.CreateOptions{
			APIID:    cmd.Flag("api-id").Value.String(),
			Region:   cmd.Flag("region").Value.String(),
			Endpoint: cmd.Flag("endpoint").Value.String(),
			Verbose:  verbose,

			Subject:      subject,
			From:         from,
			To:           to,
			Cc:           cc,
			Bcc:          bcc,
			ReplyTo:      replyTo,
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
	rootCmd.AddCommand(createCmd)
	createCmd.Flags().String("subject", "", "Subject")
	createCmd.Flags().StringArray("from", []string{}, "From")
	createCmd.Flags().StringArray("to", []string{}, "To")
	createCmd.Flags().StringArray("cc", []string{}, "Cc")
	createCmd.Flags().StringArray("bcc", []string{}, "Bcc")
	createCmd.Flags().StringArray("reply-to", []string{}, "Reply-To")
	createCmd.Flags().String("text", "", "Text")
	createCmd.Flags().String("html", "", "HTML")
	createCmd.Flags().String("generate-text", "", "Generate text from HTML (optional)")
	createCmd.Flags().String("file", "", "File")
	createCmd.Flags().Bool("send", false, "Send email immediately without using draft (optional)")
}
