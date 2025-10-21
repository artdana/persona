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

		editedProfile := editProfile(*selectedProfile)
		if editedProfile == nil {
			fmt.Println("Edit cancelled.")
			return
		}

		if err := updateProfile(*selectedProfile, *editedProfile); err != nil {
			fmt.Printf("❌ Failed to update profile: %s\n", err)
			return
		}

		fmt.Printf("✅ Profile '%s' updated successfully\n", editedProfile.Name)
	},
}

func init() {
	rootCmd.AddCommand(editCmd)
}

// editProfile interactively edits a profile's fields
func editProfile(profile persona.Profile) *persona.Profile {
	fmt.Printf("\n📝 Editing profile: %s\n", profile.Name)
	fmt.Println("Press Enter to keep current value, or type new value")
	fmt.Println()

	newName := promptForEdit("Profile name", profile.Name)
	if newName == "" {
		fmt.Println("❌ Profile name cannot be empty")
		return nil
	}

	newUser := promptForEdit("Git user name", profile.User)
	if newUser == "" {
		fmt.Println("❌ User name cannot be empty")
		return nil
	}

	newEmail := promptForEdit("Git email", profile.Email)
	if newEmail == "" {
		fmt.Println("❌ Email cannot be empty")
		return nil
	}

	newSigningKey := promptForEdit("Signing key", profile.SigningKey)

	newDescription := promptForEdit("Description", profile.Description)

	return &persona.Profile{
		Name:        newName,
		User:        newUser,
		Email:       newEmail,
		SigningKey:  newSigningKey,
		Description: newDescription,
	}
}

// promptForEdit prompts the user to edit a field with current value as default
func promptForEdit(fieldName, currentValue string) string {
	prompt := fmt.Sprintf("%s [%s]: ", fieldName, currentValue)
	fmt.Print(prompt)
	reader := bufio.NewReader(os.Stdin)
	input, _ := reader.ReadString('\n')
	input = strings.TrimSpace(input)
	
	if input == "" {
		return currentValue
	}
	
	return input
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
