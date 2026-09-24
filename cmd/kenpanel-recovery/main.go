package main

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "kenpanel-recovery",
	Short: "KenPanel Recovery and Disaster Kit CLI",
	Long:  "Emergency recovery, offline Disaster Kit runbooks, and RecoveryOS utilities.",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("KenPanel Recovery and Disaster Kit CLI v2.3.0")
		fmt.Println("Use 'kenpanel-recovery --help' for available commands.")
	},
}

func main() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}
