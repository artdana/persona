package cmd

import (
	"fmt"
	"os"
	"persona/internal/persona"
	"persona/internal/tui"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var addCmd = &cobra.Command{
	Use:   "add [profile]",
	Short: "Add a new profile",
	Long:  "Add a new git profile with name, user, email, signing key, and description",
	Run: func(cmd *cobra.Command, args []string) {
		var profileName string
		if len(args) > 0 {
			profileName = args[0]
		}

		// Check if profile already exists
		if profileName != "" && profileExists(profileName) {
			fmt.Printf("❌ Profile '%s' already exists\n", profileName)
			return
		}

		// Create form fields
		fields := tui.CreateProfileFormFields(nil, false)
		if profileName != "" {
			fields[0].Value = profileName
			fields[0].Placeholder = profileName
		}

		// Start form TUI
		values := tui.StartForm("Add New Profile", fields, func(vals map[string]string) bool {
			name := vals["Profile name"]
			if name == "" {
				return false
			}

			// Check if profile already exists
			if profileExists(name) {
				fmt.Printf("❌ Profile '%s' already exists\n", name)
				return false
			}

			// Create profile
			profile := &persona.Profile{
				Name:        name,
				User:        vals["Git user name"],
				Email:       vals["Git email"],
				SigningKey:  vals["Signing key"],
				Description: vals["Description"],
			}

			// Save profile
			if err := saveProfile(*profile); err != nil {
				fmt.Printf("❌ Failed to save profile: %s\n", err)
				return false
			}

			fmt.Printf("✅ Profile '%s' added successfully\n", name)
			return true
		})

		if values == nil {
			fmt.Println("Add cancelled.")
		}
	},
}

func init() {
	rootCmd.AddCommand(addCmd)
}

// saveProfile saves a profile to the config file
func saveProfile(profile persona.Profile) error {
	// Load existing profiles
	var profiles []persona.Profile
	if err := viper.UnmarshalKey("profiles", &profiles); err != nil {
		// If profiles key doesn't exist, start with empty slice
		profiles = []persona.Profile{}
	}

	// Add new profile
	profiles = append(profiles, profile)

	// Set the profiles in viper
	viper.Set("profiles", profiles)

	// Write config to file
	if err := viper.WriteConfig(); err != nil {
		if os.IsNotExist(err) {
			// Create config file if it doesn't exist
			return viper.SafeWriteConfig()
		}
		return err
	}

	return nil
}