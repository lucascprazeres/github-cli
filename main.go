package main

import (
	"github.com/lucascprazeres/github-cli/cmd"
	"github.com/lucascprazeres/github-cli/internal/config/env"
)

func main() {
	env.Setup()

	cmd.RootCommand().Execute()
}
