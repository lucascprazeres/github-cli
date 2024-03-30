package cmd

import (
	"github-cli/internal/commands"

	"github.com/spf13/cobra"
)

func AuthCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   `auth`,
		Short: `authenticate on github via OAuth2`,
		RunE: func(cmd *cobra.Command, args []string) error {
			return commands.Auth()
		},
	}

	return cmd
}
