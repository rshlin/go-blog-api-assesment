package main

import (
	"github.com/rshlin/go-blog-api-assesment/cmd"
	"github.com/spf13/cobra/doc"
	"log"
)

//go:generate go run gen_cli_doc.go
func cliDocGen() {
	if err := doc.GenMarkdownTree(cmd.RootCmd, "./"); err != nil {
		log.Fatalf("failed to generate markdown: %v", err)
	}
}

func main() {
	cliDocGen()
}
