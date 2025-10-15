package cmd

import (
	"bufio"
	"fmt"
	"os"
	"persona/internal/persona"
	"strings"

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
		} else {
			profileName = promptForInput("Profile name: ")
		}

		if profileName == "" {
			fmt.Println("❌ Profile name cannot be empty")
			return
		}

		// Check if profile already exists
		if profileExists(profileName) {
			fmt.Printf("❌ Profile '%s' already exists\n", profileName)
			return
		}

		// Create new profile
		profile := createProfile(profileName)
		if profile == nil {
			return
		}

		// Save profile to config
		if err := saveProfile(*profile); err != nil {
			fmt.Printf("❌ Failed to save profile: %s\n", err)
			return
		}

		fmt.Printf("✅ Profile '%s' added successfully\n", profileName)
	},
}

func init() {
	rootCmd.AddCommand(addCmd)
}

// promptForInput prompts the user for input and returns the trimmed result
func promptForInput(prompt string) string {
	fmt.Print(prompt)
	reader := bufio.NewReader(os.Stdin)
	input, _ := reader.ReadString('\n')
	return strings.TrimSpace(input)
}

// createProfile interactively creates a new profile
func createProfile(name string) *persona.Profile {
	user := promptForInput("Git user name: ")
	if user == "" {
		fmt.Println("❌ User name cannot be empty")
		return nil
	}

	email := promptForInput("Git email: ")
	if email == "" {
		fmt.Println("❌ Email cannot be empty")
		return nil
	}

	signingKey := promptForInput("Signing key (optional): ")
	description := promptForInput("Description (optional): ")

	return &persona.Profile{
		Name:        name,
		User:        user,
		Email:       email,
		SigningKey:  signingKey,
		Description: description,
	}
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