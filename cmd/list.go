package cmd

import (
	"github.com/lucascprazeres/github-cli/internal/issues"
	"github.com/lucascprazeres/github-cli/internal/logging"
	"github.com/spf13/cobra"
)

func ListCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   `list`,
		Short: `list issues of a repository`,
		Args:  ValidateArgs,
		Run: func(cmd *cobra.Command, args []string) {
			is, _ := cmd.Flags().GetString("is")
			if len(is) > 0 && is != "open" && is != "closed" {
				logging.Error("the only valid options for 'is' are 'open' and 'closed'\n")
				return
			}

			author, _ := cmd.Flags().GetString("author")

			result, err := issues.List(args[0], is, author)
			if err != nil {
				logging.Error("%s\n", err)
				return
			}

			logging.Prompt("%v issues found:\n", result.TotalCount)

			if result.TotalCount == 0 {
				return
			}

			logging.Info("[#%-5s]\t %15s\t %.55s\n", "ISSUE", "AUTHOR", "TITLE")
			for _, item := range result.Items {
				logging.Success("[#%-5d]\t %15.15s\t %.55s\n", item.Number, item.User.Login, item.Title)
			}
		},
	}

	cmd.PersistentFlags().String("is", "", "Filter issues by state")
	cmd.PersistentFlags().String("author", "", "Filter issues by author")

	return cmd
}
