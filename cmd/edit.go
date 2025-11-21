package cmd

import (
	"fmt"
	"persona/internal/persona"
	"persona/internal/tui"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var editCmd = &cobra.Command{
	Use:   "edit [profile]",
	Short: "Edit a profile",
	Long:  "Edit an existing git profile. Select a profile and modify its fields.",
	Run: func(cmd *cobra.Command, args []string) {
		var profiles []persona.Profile
		if err := viper.UnmarshalKey("profiles", &profiles); err != nil || len(profiles) == 0 {
			fmt.Println("❌ No profiles found. Run `persona add` first.")
			return
		}

		var selectedProfile *persona.Profile

		if len(args) > 0 {
			profileName := args[0]
			
			for _, profile := range profiles {
				if profile.Name == profileName {
					selectedProfile = &profile
					break
				}
			}

			if selectedProfile == nil {
				fmt.Printf("❌ Profile '%s' not found.\n", profileName)
				fmt.Println("Available profiles:")
				for _, profile := range profiles {
					fmt.Printf("  - %s\n", profile.Name)
				}
				return
			}
		} else {
			selected := tui.StartProfileSelector(profiles, viper.GetString("active_profile"))
			if selected == nil {
				return
			}
			selectedProfile = selected
		}

		// Create form fields with current profile values
		fields := tui.CreateProfileFormFields(selectedProfile, true)

		// Start form TUI
		values := tui.StartForm("Edit Profile", fields, func(vals map[string]string) bool {
			// Get new values from form (form returns original value if field wasn't changed)
			editedProfile := &persona.Profile{
				Name:        vals["Profile name"],
				User:        vals["Git user name"],
				Email:       vals["Git email"],
				SigningKey:  vals["Signing key"],
				Description: vals["Description"],
			}

			if err := updateProfile(*selectedProfile, *editedProfile); err != nil {
				fmt.Printf("❌ Failed to update profile: %s\n", err)
				return false
			}

			fmt.Printf("✅ Profile '%s' updated successfully\n", editedProfile.Name)
			return true
		})

		if values == nil {
			fmt.Println("Edit cancelled.")
		}
	},
}

func init() {
	rootCmd.AddCommand(editCmd)
}


// updateProfile updates a profile in the config
func updateProfile(oldProfile, newProfile persona.Profile) error {
	var profiles []persona.Profile
	if err := viper.UnmarshalKey("profiles", &profiles); err != nil {
		return fmt.Errorf("failed to load profiles: %w", err)
	}

	found := false
	for i, profile := range profiles {
		if profile.Name == oldProfile.Name {
			profiles[i] = newProfile
			found = true
			break
		}
	}

	if !found {
		return fmt.Errorf("profile '%s' not found", oldProfile.Name)
	}

	if oldProfile.Name != newProfile.Name {
		for _, profile := range profiles {
			if profile.Name == newProfile.Name && profile != newProfile {
				return fmt.Errorf("profile name '%s' already exists", newProfile.Name)
			}
		}
	}

	viper.Set("profiles", profiles)

	if err := viper.WriteConfig(); err != nil {
		return fmt.Errorf("failed to save config: %w", err)
	}

	return nil
}
