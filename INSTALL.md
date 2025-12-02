# Installation Script

This script automates the installation of all dependencies required for the Go Gin project.

## Features

✅ **Automatic Installation** of:
- Make (build tool)
- Golang (v1.25.4)
- sqlc (SQL compiler)
- golang-migrate (database migration tool)
- Docker (optional)
- Docker Compose (optional)
- Air (hot reload for development)

✅ **Smart Error Handling**:
- Tracks all installations during the session
- Automatically offers to rollback on errors
- Prevents partial installations from breaking your system
- **Single sudo password prompt** - enter once at the start, not repeatedly

✅ **OS Support**:
- Linux (Debian/Ubuntu, RHEL/CentOS/Fedora, Arch)
- macOS (via Homebrew)

## Usage

### Install Dependencies

```bash
./install.sh
```

This will:
1. Check for existing installations
2. Install missing dependencies
3. Configure environment variables
4. Verify all installations
5. Provide next steps

### Uninstall Dependencies

If something goes wrong during installation, you'll be prompted to automatically rollback. You can also manually uninstall:

```bash
./install.sh uninstall
```

This removes:
- All installed tools (Go, sqlc, migrate, Air, Docker, Docker Compose)
- Environment variables from shell config
- Temporary files

### Get Help

```bash
./install.sh help
```

## What Happens on Error?

If installation fails:

1. **Error Detection**: The script detects the failure immediately
2. **Prompt**: You'll be asked if you want to rollback
3. **Cleanup**: If you choose yes, all dependencies installed *during this run* will be removed
4. **Clean State**: Your system returns to the state before running the script

Example:
```
✗ Installation failed with exit code 1

Do you want to rollback and uninstall all dependencies installed during this run? (y/N):
```

## After Installation

1. **Restart your terminal** or run:
   ```bash
   source ~/.bashrc  # or ~/.zshrc if using zsh
   ```

2. **Set up your project**:
   ```bash
   cp .env.example .env
   # Edit .env with your database credentials
   ```

3. **Start development**:
   ```bash
   make dev          # Full development environment
   # OR
   make noapp        # Just the database
   ```

4. **Run migrations**:
   ```bash
   make migrate_up
   ```

5. **Generate SQL code**:
   ```bash
   make sqlc
   ```

## Troubleshooting

### Installation fails but no rollback

If you skipped the rollback prompt, you can manually clean up:

```bash
./install.sh uninstall
```

### Environment variables not working

Restart your terminal or manually source your shell config:

```bash
source ~/.bashrc  # or ~/.zshrc
```

### Permission errors

Don't run the script as root:
```bash
# ✗ Wrong
sudo ./install.sh

# ✓ Correct
./install.sh
```

The script will request sudo only when needed for specific operations.

### Sudo password prompt

The script will ask for your sudo password **once** at the beginning if needed. The credentials are cached throughout the installation/uninstallation process, so you won't be prompted repeatedly.

If you want to completely avoid the sudo password prompt (advanced users only), you can configure passwordless sudo for your user:

```bash
# Edit sudoers file (be careful!)
sudo visudo

# Add this line at the end (replace 'yourusername' with your actual username):
yourusername ALL=(ALL) NOPASSWD: ALL
```

⚠️ **Warning**: This reduces system security. Only do this on development machines.

## What Gets Tracked?

The script only removes tools that were **installed during the current run**. Pre-existing installations are never touched during rollback.

For example:
- If Go was already installed → Rollback won't remove it
- If Go was installed during this run → Rollback will remove it

## Manual Uninstall

The `./install.sh uninstall` command removes **everything**, regardless of when it was installed. Use this for a complete cleanup.
