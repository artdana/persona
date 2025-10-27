package cmd

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"persona/internal/persona"
	"strings"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var initCmd = &cobra.Command{
	Use:   "init",
	Short: "Initialize persona and create default profile from git config",
	Long:  "Initialize persona by creating the config file and setting up a default profile from your current global git configuration.",
	Run: func(cmd *cobra.Command, args []string) {
		if configExists() {
			fmt.Println("❌ Persona is already initialized.")
			fmt.Println("   If you want to add a new profile, use `persona add`")
			return
		}

		gitUserName, err := exec.Command("git", "config", "--global", "user.name").Output()
		if err != nil {
			fmt.Println("❌ Failed to read git global user.name")
			return
		}
		userName := strings.TrimSpace(string(gitUserName))
		if userName == "" {
			fmt.Println("❌ Git global user.name is not configured")
			fmt.Println("   Please run: git config --global user.name \"Your Name\"")
			return
		}

		gitUserEmail, err := exec.Command("git", "config", "--global", "user.email").Output()
		if err != nil {
			fmt.Println("❌ Failed to read git global user.email")
			return
		}
		userEmail := strings.TrimSpace(string(gitUserEmail))
		if userEmail == "" {
			fmt.Println("❌ Git global user.email is not configured")
			fmt.Println("   Please run: git config --global user.email \"your.email@example.com\"")
			return
		}

		if err := ensureConfigDir(); err != nil {
			fmt.Printf("❌ Failed to create config directory: %s\n", err)
			return
		}

		defaultProfile := persona.Profile{
			Name:       "default",
			User:       userName,
			Email:      userEmail,
			SigningKey: "",
		}

		home, err := os.UserHomeDir()
		if err != nil {
			fmt.Printf("❌ Failed to get home directory: %s\n", err)
			return
		}

		configPath := filepath.Join(home, ".config", "persona", "config.yaml")
		viper.SetConfigFile(configPath)
		viper.SetConfigType("yaml")

		profiles := []persona.Profile{defaultProfile}
		viper.Set("profiles", profiles)
		viper.Set("active_profile", "default")

		if err := viper.WriteConfig(); err != nil {
			fmt.Printf("❌ Failed to write config: %s\n", err)
			return
		}

		fmt.Println("✅ Persona initialized successfully!")
		fmt.Printf("   Created default profile from your git config:\n")
		fmt.Printf("   User: %s\n", userName)
		fmt.Printf("   Email: %s\n", userEmail)
		fmt.Println("\n   You can now use `persona add` to create additional profiles.")
	},
}

func init() {
	rootCmd.AddCommand(initCmd)
}

func configExists() bool {
	home, err := os.UserHomeDir()
	if err != nil {
		return false
	}
	configPath := filepath.Join(home, ".config", "persona", "config.yaml")
	_, err = os.Stat(configPath)
	return err == nil
}

func ensureConfigDir() error {
	home, err := os.UserHomeDir()
	if err != nil {
		return err
	}
	configDir := filepath.Join(home, ".config", "persona")
	return os.MkdirAll(configDir, 0755)
}
