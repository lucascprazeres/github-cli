package cmd

import (
	"github-cli/internal/commands"

	"github.com/spf13/cobra"
)

func ListCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   `list`,
		Short: `list issues of a repository`,
		Args:  ValidateArgs,
		RunE:  listRun,
	}

	cmd.PersistentFlags().String("is", "", "Filter issues by state")
	cmd.PersistentFlags().String("author", "", "Filter issues by author")

	return cmd
}

func listRun(cmd *cobra.Command, args []string) error {
	repo := args[0]
	is, _ := cmd.Flags().GetString("is")
	author, _ := cmd.Flags().GetString("author")

	err := commands.List(repo, is, author)
	if err != nil {
		return err
	}

	return nil
}
