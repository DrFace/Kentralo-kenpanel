package main

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "kenpanel-capsule",
	Short: "KenPanel App Capsule Portable Packager CLI",
	Long:  "Self-contained application snapshot and restore engine with secret redaction.",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("KenPanel App Capsule Packager v2.3.0")
		fmt.Println("Use 'kenpanel-capsule --help' for available commands.")
	},
}

func main() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}
