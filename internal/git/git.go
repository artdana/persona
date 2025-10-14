package git

import (
	"fmt"
	"os/exec"
	"persona/internal/persona"
)

// change global git config
func ApplyProfileGlobal(profile persona.Profile) error {
	if err := exec.Command("git", "config", "--global", "user.name", profile.User).Run(); err != nil {
		return fmt.Errorf("failed to set git config: %s", err)
	}
	if err := exec.Command("git", "config", "--global", "user.email", profile.Email).Run(); err != nil {
		return fmt.Errorf("failed to set git config: %s", err)
	}
	return nil
}


// change local repo git config
func ApplyProfileLocal(profile persona.Profile) error {
	if err := exec.Command("git", "config", "user.name", profile.User).Run(); err != nil {
		return fmt.Errorf("failed to set git config: %s", err)
	}
	if err := exec.Command("git", "config", "user.email", profile.Email).Run(); err != nil {
		return fmt.Errorf("failed to set git config: %s", err)
	}
	return nil
}