package cmd

import (
	"fmt"
	"persona/internal/persona"

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

		fmt.Println("📋 Available Profiles:")
		fmt.Println()

		for i, profile := range profiles {
			activeIndicator := ""
			if profile.Name == activeProfile {
				activeIndicator = " ✅ (active)"
			}

			fmt.Printf("%d. %s%s\n", i+1, profile.Name, activeIndicator)
			fmt.Printf("   User: %s\n", profile.User)
			fmt.Printf("   Email: %s\n", profile.Email)
			
			if profile.SigningKey != "" {
				fmt.Printf("   Signing Key: %s\n", profile.SigningKey)
			}
			
			if profile.Description != "" {
				fmt.Printf("   Description: %s\n", profile.Description)
			}
			
			if i < len(profiles)-1 {
				fmt.Println()
			}
		}

		fmt.Println()
		fmt.Printf("Total: %d profile(s)\n", len(profiles))
	},
}

func init() {
	rootCmd.AddCommand(listCmd)
}
