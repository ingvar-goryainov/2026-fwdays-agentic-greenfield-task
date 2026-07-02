package main

import (
	"os"

	"github.com/ingvar-goryainov/agents-lint/cmd/agents-lint/cmd"
)

func main() {
	if err := cmd.Execute(); err != nil {
		os.Exit(1)
	}
}
