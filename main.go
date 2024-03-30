package main

import (
	"github-cli/cmd"
	"github-cli/internal/config/env"
)

func main() {
	env.Setup()

	cmd.RootCommand().Execute()
}
