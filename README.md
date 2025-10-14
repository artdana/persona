# Persona
> A simple CLI tool to manage multiple Git identities and switch between them seamlessly.

---

⚠️ **Work in Progress**

This project is still under active development.  
Currently, only the `persona use` command is implemented — which allows you to switch between your saved Git profiles locally or globally.

---

## 🚀 Current Status

✅ `persona use`  
> Switch between Git profiles and update your system Git config automatically.

Example:
```bash
# to choose from the list of profiles that you have set up
persona use
persona use -g
persona use --global

# to switch directly to the profile
persona use <profile-name>
persona use <profile-name> --g
persona use <profile-name> --global
