package cmd

import (
	"fmt"
	"persona/internal/persona"
	"persona/internal/tui"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var listCmd = &cobra.Command{
	Use:   "list",
	Short: "List all profiles",
	Long:  "Display all configured git profiles with their details",
	Run: func(cmd *cobra.Command, args []string) {
		var profiles []persona.Profile
		if err := viper.UnmarshalKey("profiles", &profiles); err != nil || len(profiles) == 0 {
			fmt.Println("❌ No profiles found. Run `persona add` first.")
			return
		}

		activeProfile := viper.GetString("active_profile")
		tui.StartList(profiles, activeProfile, "Available Profiles")
	},
}

func init() {
	rootCmd.AddCommand(listCmd)
}
