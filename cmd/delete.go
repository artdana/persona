package cmd

import (
	"bufio"
	"fmt"
	"os"
	"persona/internal/persona"
	"persona/internal/tui"
	"strings"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var deleteCmd = &cobra.Command{
	Use:   "delete [profile]",
	Short: "Delete a profile",
	Long:  "Delete a git profile from the configuration. Active profiles cannot be deleted.",
	Run: func(cmd *cobra.Command, args []string) {
		var profiles []persona.Profile
		if err := viper.UnmarshalKey("profiles", &profiles); err != nil || len(profiles) == 0 {
			fmt.Println("❌ No profiles found.")
			return
		}

		activeProfile := viper.GetString("active_profile")

		var deletableProfiles []persona.Profile
		for _, profile := range profiles {
			if profile.Name != activeProfile {
				deletableProfiles = append(deletableProfiles, profile)
			}
		}

		if len(deletableProfiles) == 0 {
			fmt.Println("❌ No profiles available for deletion.")
			if activeProfile != "" {
				fmt.Printf("Active profile '%s' cannot be deleted. Switch to another profile first.\n", activeProfile)
			}
			return
		}

		var selectedProfile *persona.Profile

		if len(args) > 0 {
			profileName := args[0]
			
			if profileName == activeProfile {
				fmt.Printf("❌ Cannot delete active profile '%s'. Switch to another profile first.\n", profileName)
				return
			}

			for _, profile := range deletableProfiles {
				if profile.Name == profileName {
					selectedProfile = &profile
					break
				}
			}

			if selectedProfile == nil {
				fmt.Printf("❌ Profile '%s' not found or not available for deletion.\n", profileName)
				fmt.Println("Available profiles for deletion:")
				for _, profile := range deletableProfiles {
					fmt.Printf("  - %s\n", profile.Name)
				}
				return
			}
		} else {
			selected := tui.StartProfileSelector(deletableProfiles, "")
			if selected == nil {
				return
			}
			selectedProfile = selected
		}

		if !confirmDeletion(selectedProfile.Name) {
			fmt.Println("Deletion cancelled.")
			return
		}

		if err := deleteProfile(selectedProfile.Name); err != nil {
			fmt.Printf("❌ Failed to delete profile: %s\n", err)
			return
		}

		fmt.Printf("✅ Profile '%s' deleted successfully\n", selectedProfile.Name)
	},
}

func init() {
	rootCmd.AddCommand(deleteCmd)
}

func confirmDeletion(profileName string) bool {
	fmt.Printf("Are you sure you want to delete profile '%s'? (y/N): ", profileName)
	reader := bufio.NewReader(os.Stdin)
	response, _ := reader.ReadString('\n')
	response = strings.TrimSpace(strings.ToLower(response))
	return response == "y" || response == "yes"
}

func deleteProfile(profileName string) error {
	var profiles []persona.Profile
	if err := viper.UnmarshalKey("profiles", &profiles); err != nil {
		return fmt.Errorf("failed to load profiles: %w", err)
	}

	var updatedProfiles []persona.Profile
	found := false
	for _, profile := range profiles {
		if profile.Name != profileName {
			updatedProfiles = append(updatedProfiles, profile)
		} else {
			found = true
		}
	}

	if !found {
		return fmt.Errorf("profile '%s' not found", profileName)
	}

	viper.Set("profiles", updatedProfiles)

	if err := viper.WriteConfig(); err != nil {
		return fmt.Errorf("failed to save config: %w", err)
	}

	return nil
}
