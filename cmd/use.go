package cmd

import (
	"fmt"
	"os"
	"persona/internal/git"
	"persona/internal/persona"
	"persona/internal/tui"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var useCmd = &cobra.Command{
    Use:   "use [profile]",
    Short: "Select and switch between profiles",
    Run: func(cmd *cobra.Command, args []string) {
        global, _ := cmd.Flags().GetBool("global")
        
        if len(args) > 0 {
			setActiveProfile(args[0], global)
			return
		}
		
		var profiles []persona.Profile
        if err := viper.UnmarshalKey("profiles", &profiles); err != nil || len(profiles) == 0 {
            println("❌ No profiles found. Run `persona add` first.")
            return
        }

        selected := tui.StartProfileSelector(profiles, viper.GetString("active_profile"))
        if selected != nil {
            setActiveProfile(selected.Name, global)
        }
    },
}

func init() {
    useCmd.Flags().BoolP("global", "g", false, "Set profile globally for all repositories")
    rootCmd.AddCommand(useCmd)
}

// load profiles from config
func loadProfiles() []persona.Profile {
	var profiles []persona.Profile
	if err := viper.UnmarshalKey("profiles", &profiles); err != nil {
		fmt.Printf("❌ Failed to load profiles: %s\n", err)
		return nil
	}
	return profiles
}

// check if profile exists
func profileExists(name string) bool {
	profiles := loadProfiles()
	for _, profile := range profiles {
		if profile.Name == name {
			return true
		}
	}
	return false
}

// set active profile
func setActiveProfile(name string, global bool) {
	profileExists := profileExists(name)

	if !profileExists {
		fmt.Printf("❌ Profile '%s' not found. Available profiles:\n", name)
		for _, profile := range loadProfiles() {
			fmt.Printf("  - %s\n", profile.Name)
		}
		return
	}

	profiles := loadProfiles()
	var selectedProfile *persona.Profile
	for _, profile := range profiles {
		if profile.Name == name {
			selectedProfile = &profile
			break
		}
	}

	viper.Set("active_profile", name)
	err := viper.WriteConfig()
	if err != nil {
		if os.IsNotExist(err) {
			err = viper.SafeWriteConfig()
		}
		if err != nil {
			fmt.Printf("❌ Failed to save config: %s", err)
			return
		}
	}

	if global {
		if err := git.ApplyProfileGlobal(*selectedProfile); err != nil {
			fmt.Printf("❌ Failed to apply global git config: %s\n", err)
			return
		}
	} else {
		if err := git.ApplyProfileLocal(*selectedProfile); err != nil {
			fmt.Printf("❌ Failed to apply local git config: %s\n", err)
			return
		}
	}
	
	scope := "locally"
	if global {
		scope = "globally"
	}
	fmt.Printf("✅ Switched to profile: %s (%s)\n", name, scope)
}
