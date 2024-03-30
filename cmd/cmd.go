package cmd

import (
	"errors"

	"github-cli/internal/logging"

	"github.com/spf13/cobra"
)

func RootCommand() *cobra.Command {
	rootCmd := &cobra.Command{
		Use:   `ghi`,
		Short: `GitHub Issues CLI`,
		Run: func(cmd *cobra.Command, args []string) {
			logging.Success("Bem-vindo ao Github CLI! 🚀\n")
		},
	}

	rootCmd.AddCommand(ListCommand())
	rootCmd.AddCommand(AuthCommand())
	rootCmd.AddCommand(CreateCommand())

	return rootCmd
}

func ValidateArgs(cmd *cobra.Command, args []string) error {
	if len(args) < 1 {
		return errors.New("cli: at least one argument must be informed")
	}
	return nil
}
