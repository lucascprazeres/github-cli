package main

import "github.com/lucascprazeres/github-cli/cmd"

func main() {
	rootCmd := cmd.RootCommand()
	rootCmd.Execute()
}
