package main

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "kenpanel-mailbridge",
	Short: "KenPanel MailBridge IMAP Migration Engine CLI",
	Long:  "Flag-preserving streaming IMAP sync engine with zero message body logging.",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("KenPanel MailBridge IMAP Migration Engine v2.3.0")
		fmt.Println("Use 'kenpanel-mailbridge --help' for available commands.")
	},
}

func main() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}
