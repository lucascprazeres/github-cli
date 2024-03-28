package cmd

import (
	"fmt"

	"github.com/lucascprazeres/github-cli/internal/issues"
	"github.com/lucascprazeres/github-cli/internal/logging"
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
	is, _ := cmd.Flags().GetString("is")
	if is != "" && is != "open" && is != "closed" {
		return fmt.Errorf("the only valid options for 'is' are 'open' and 'closed'")
	}

	author, _ := cmd.Flags().GetString("author")

	result, err := issues.List(args[0], is, author)
	if err != nil {
		return err
	}

	logging.Prompt("%v issues found:\n", result.TotalCount)

	if result.TotalCount == 0 {
		return nil
	}

	logging.Info("[#%-5s]\t %15s\t %.55s\n", "ISSUE", "AUTHOR", "TITLE")
	for _, item := range result.Items {
		logging.Success("[#%-5d]\t %15.15s\t %.55s\n", item.Number, item.User.Login, item.Title)
	}

	return nil
}
