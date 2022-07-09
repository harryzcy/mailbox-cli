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
	Run: func(cmd *cobra.Command, args []string) {
		verbose, _ := cmd.Flags().GetBool("verbose")
		subject, _ := cmd.Flags().GetString("subject")
		from, _ := cmd.Flags().GetStringArray("from")
		to, _ := cmd.Flags().GetStringArray("to")
		cc, _ := cmd.Flags().GetStringArray("cc")
		bcc, _ := cmd.Flags().GetStringArray("bcc")
		replyTo, _ := cmd.Flags().GetStringArray("reply-to")
		text, _ := cmd.Flags().GetString("text")
		html, _ := cmd.Flags().GetString("html")
		file, _ := cmd.Flags().GetString("file")
		generateText, _ := cmd.Flags().GetString("generate-text")
		send, _ := cmd.Flags().GetBool("send")

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
