package cmd

import (
	"os"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "app",
	Short: "Backend Developer Skill Assessment",
	Long: `This application contains two commands to complete the Backend Developer Skill Assessment:
- logicalTest: Takes a string of digits and returns the number of ways it can be decoded back into its original message.
- technicalTest: Starts a server to handle a RESTful API for a blog platform.`,
}

func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}
