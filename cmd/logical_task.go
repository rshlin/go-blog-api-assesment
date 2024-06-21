package cmd

import (
	"fmt"
	"github.com/rshlin/go-blog-api-assesment/logical"
	"github.com/spf13/cobra"
	"log"
)

var logicalTestCmd = &cobra.Command{
	Use:   "logical",
	Short: "Logical Test: Decoding a message",
	Long:  `This command takes a string of digits and returns the number of ways it can be decoded back into its original message.`,
	Run: func(cmd *cobra.Command, args []string) {
		input, _ := cmd.Flags().GetString("digits")
		fmt.Printf("Number of ways to decode message: %d\n", logical.NumDecodings(input))
	},
}

var digits string

func init() {
	rootCmd.AddCommand(logicalTestCmd)

	logicalTestCmd.Flags().StringVarP(&digits, "digits", "d", "", "string of digits to decode (required)")
	if err := logicalTestCmd.MarkFlagRequired("digits"); err != nil {
		log.Fatalf("Error marking digits flag as required: %v", err)
	}
}
