package cmd

import (
	"github.com/lucascprazeres/github-cli/internal/commands"
	"github.com/spf13/cobra"
)

func AuthCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   `auth`,
		Short: `auth issues on your own repositories`,
		RunE: func(cmd *cobra.Command, args []string) error {
			return commands.Auth()
		},
	}

	return cmd
}
