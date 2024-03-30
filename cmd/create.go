package cmd

import (
	"github-cli/internal/commands"

	"github.com/spf13/cobra"
)

func CreateCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   `create`,
		Short: `create issues on GitHub`,
		RunE: func(cmd *cobra.Command, args []string) error {
			repo := args[0]
			title, _ := cmd.Flags().GetString("title")
			body, _ := cmd.Flags().GetString("body")

			return commands.Create(repo, title, body)
		},
	}

	cmd.PersistentFlags().StringP("title", "t", "", "title of the issue")
	cmd.PersistentFlags().StringP("body", "b", "", "body of the issue")

	return cmd
}
