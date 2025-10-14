package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var helloCmd = &cobra.Command{
	Use:   "hello",
	Short: "Persona says hello.",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Hello from Persona!")
	},
}

func init() {
	rootCmd.AddCommand(helloCmd)
}