# Persona

> A simple CLI tool to manage multiple Git identities and switch between them seamlessly.

---

## ✨ Features

Persona provides a beautiful, interactive TUI (Terminal User Interface) for managing your Git profiles. All commands feature a consistent, modern interface built with [Bubble Tea](https://github.com/charmbracelet/bubbletea).

- 🎨 **Beautiful TUI** - Consistent, modern interface across all commands
- 🔄 **Quick Switching** - Switch between profiles locally or globally
- 📝 **Profile Management** - Add, edit, list, and delete profiles
- 🔍 **Interactive Selection** - Search and filter profiles with ease
- ⚡ **Fast & Lightweight** - Built with Go for performance

---

## 🚀 Commands

### `persona` (or `persona --help`)

Display the current active profile and Git identity information.

```bash
persona
```

### `persona init`

Initialize Persona and create a default profile from your current global Git configuration.

```bash
persona init
```

This command:

- Creates the config directory (`~/.config/persona/`)
- Reads your current global Git user.name and user.email
- Creates a default profile with these values
- Sets the default profile as active

### `persona add [profile]`

Add a new Git profile interactively. Opens a beautiful form to enter:

- Profile name
- Git user name
- Git email
- Signing key (optional)
- Description (optional)

```bash
# Interactive form
persona add

# Pre-fill profile name
persona add work-profile
```

### `persona list`

Display all configured profiles in a styled list view with their details.

```bash
persona list
```

Shows:

- Profile name (with ✅ indicator for active profile)
- User name
- Email
- Signing key (if set)
- Description (if set)
- Total profile count

### `persona edit [profile]`

Edit an existing profile. Opens an interactive form pre-filled with current values.

```bash
# Select profile from list
persona edit

# Edit specific profile directly
persona edit work-profile
```

### `persona delete [profile]`

Delete a profile from your configuration. Active profiles cannot be deleted.

```bash
# Select profile from list
persona delete

# Delete specific profile directly
persona delete old-profile
```

### `persona use [profile]`

Switch between Git profiles and update your Git config automatically.

```bash
# Interactive profile selector
persona use
persona use -g          # Apply globally
persona use --global    # Apply globally

# Switch directly to a profile
persona use work-profile
persona use work-profile -g
persona use work-profile --global
```

**Options:**

- `-g, --global` - Apply the profile globally (to all repositories)

---

## 📋 Keyboard Shortcuts

All TUI commands support the following shortcuts:

- **Esc** or **Ctrl+C** - Cancel/Exit
- **Tab** or **↓** - Next field/item
- **Shift+Tab** or **↑** - Previous field/item
- **Enter** - Submit/Select
- **Type to filter** - In profile selectors, type to filter profiles
- **q** - Quit (in list/info views)

---

## 🎯 Usage Examples

### Initial Setup

```bash
# 1. Initialize Persona
persona init

# 2. Add additional profiles
persona add work-profile
persona add personal-profile

# 3. Switch to a profile
persona use work-profile
```

### Daily Workflow

```bash
# Check current profile
persona

# List all profiles
persona list

# Switch profiles
persona use work-profile      # Local (current repo)
persona use work-profile -g   # Global (all repos)

# Edit a profile
persona edit work-profile
```

---

## 📁 Configuration

Persona stores its configuration in:

```
~/.config/persona/config.yaml
```

The config file contains:

- List of all profiles
- Currently active profile
- Profile details (name, user, email, signing key, description)

---

## 🛠️ Installation

```bash
# Clone the repository
git clone <repository-url>
cd persona

# Build
go build

# Install (optional)
go install
```

---

## 📝 Notes

- Active profiles cannot be deleted. Switch to another profile first.
- Profiles are stored locally in your config directory.
- The `use` command updates your Git config immediately.
- All commands feature a consistent, beautiful TUI interface.
