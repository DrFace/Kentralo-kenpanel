package main

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "kenpanel-transplant",
	Short: "KenPanel Transplant Migration Engine CLI",
	Long:  "Vendor-neutral migration parser and streaming importer for cPanel/WHM and Plesk.",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("KenPanel Transplant Migration Engine v2.3.0")
		fmt.Println("Use 'kenpanel-transplant --help' for available commands.")
	},
}

func main() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}
