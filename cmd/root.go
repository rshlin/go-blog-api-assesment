package cmd

import (
	"os"

	"github.com/spf13/cobra"
)

var RootCmd = &cobra.Command{
	Use:   "app",
	Short: "Backend Developer Skill Assessment",
	Long: `This application contains two commands to complete the Backend Developer Skill Assessment:
- logical: Takes a string of digits and returns the number of ways it can be decoded back into its original message.
- technical: Starts a server to handle a RESTful API for a blog platform.`,
}

func Execute() {
	err := RootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}
