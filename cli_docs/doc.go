package main

import (
	"log"

	"github.com/nikhilsbhat/terragen/cmd"
	"github.com/spf13/cobra/doc"
)

//go:generate go run github.com/nikhilsbhat/terragen/cli_docs
func main() {
	commands := cmd.SetTerragenCommands()
	err := doc.GenMarkdownTree(commands, "doc")
	if err != nil {
		log.Fatal(err)
	}
}
