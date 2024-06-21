package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
)

var technicalTestCmd = &cobra.Command{
	Use:   "technical",
	Short: "Technical Test: RESTful API for a blog platform",
	Long:  `This command starts a server for a RESTful API for a blog platform.`,
	Run: func(cmd *cobra.Command, args []string) {
		address, _ := cmd.Flags().GetString("address")
		port, _ := cmd.Flags().GetInt("port")
		fmt.Printf("Starting server at %s:%d\n", address, port)
	},
}

var address string
var port int

func init() {
	rootCmd.AddCommand(technicalTestCmd)

	technicalTestCmd.Flags().StringVarP(&address, "address", "a", "127.0.0.1", "bind address for the server")
	technicalTestCmd.Flags().IntVarP(&port, "port", "p", 8080, "port for the server")
}

func startServer(address string, port int) {
}
