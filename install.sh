#!/bin/bash

# Error handling is managed by the trap handler below

echo "=================================="
echo "Go Gin Project Installation Script"
echo "=================================="
echo ""

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

# Function to print colored messages
print_success() {
    echo -e "${GREEN}✓ $1${NC}"
}

print_error() {
    echo -e "${RED}✗ $1${NC}"
}

print_info() {
    echo -e "${YELLOW}ℹ $1${NC}"
}

# Detect OS
OS="$(uname -s)"
ARCH="$(uname -m)"

echo "Detected OS: $OS"
echo "Detected Architecture: $ARCH"
echo ""

# Default versions
GO_VERSION="1.25.4" # golang version
SQLC_VERSION="1.25.0" # sqlc version
MIGRATE_VERSION="4.17.0" # golang-migrate version

# Track what was installed during this run
INSTALLED_MAKE=false
INSTALLED_GO=false
INSTALLED_SQLC=false
INSTALLED_MIGRATE=false
INSTALLED_DOCKER=false
INSTALLED_DOCKER_COMPOSE=false
INSTALLED_AIR=false

# Track errors
INSTALLATION_FAILED=false

# Check if running as root
if [ "$EUID" -eq 0 ]; then 
    print_error "Please do not run this script as root"
    exit 1
fi

# Function to check if a command exists
command_exists() {
    command -v "$1" >/dev/null 2>&1
}

# Function to validate sudo access (caches credentials)
validate_sudo() {
    # Check if any operation will need sudo
    local needs_sudo=false
    
    if [ "$INSTALLED_DOCKER_COMPOSE" = true ] && [ -f "/usr/local/bin/docker-compose" ]; then
        needs_sudo=true
    fi
    if [ "$INSTALLED_DOCKER" = true ]; then
        needs_sudo=true
    fi
    if [ "$INSTALLED_MIGRATE" = true ] && [ -f "/usr/local/bin/migrate" ]; then
        needs_sudo=true
    fi
    if [ "$INSTALLED_GO" = true ] && [ -d "/usr/local/go" ]; then
        needs_sudo=true
    fi
    
    if [ "$needs_sudo" = true ]; then
        print_info "Some operations require sudo access. Please enter your password once:"
        sudo -v || {
            print_error "Failed to obtain sudo access"
            return 1
        }
        # Keep sudo alive in background
        while true; do sudo -n true; sleep 50; kill -0 "$$" || exit; done 2>/dev/null &
        SUDO_KEEPALIVE_PID=$!
    fi
}

# ================================
# Uninstall all installed dependencies
# ================================
uninstall_all() {
    echo ""
    echo "=================================="
    echo "Rolling Back Installations"
    echo "=================================="
    echo ""
    
    # Validate sudo access once at the beginning
    validate_sudo
    
    print_info "Cleaning up installed dependencies..."
    
    # Uninstall Air
    if [ "$INSTALLED_AIR" = true ]; then
        print_info "Removing Air..."
        if [ -f "$HOME/go/bin/air" ]; then
            rm -f "$HOME/go/bin/air"
            print_success "Air removed"
        fi
    fi
    
    # Uninstall Docker Compose
    if [ "$INSTALLED_DOCKER_COMPOSE" = true ]; then
        print_info "Removing Docker Compose..."
        if [ -f "/usr/local/bin/docker-compose" ]; then
            sudo rm -f /usr/local/bin/docker-compose
            print_success "Docker Compose removed"
        fi
    fi
    
    # Uninstall Docker
    if [ "$INSTALLED_DOCKER" = true ]; then
        print_info "Removing Docker..."
        case "$OS" in
            Linux*)
                if command_exists apt-get; then
                    sudo apt-get purge -y docker-ce docker-ce-cli containerd.io docker-buildx-plugin docker-compose-plugin
                    sudo rm -rf /var/lib/docker
                    sudo rm -rf /var/lib/containerd
                elif command_exists yum || command_exists dnf; then
                    sudo yum remove -y docker-ce docker-ce-cli containerd.io docker-buildx-plugin docker-compose-plugin
                    sudo rm -rf /var/lib/docker
                    sudo rm -rf /var/lib/containerd
                fi
                # Remove user from docker group
                sudo gpasswd -d $USER docker 2>/dev/null || true
                print_success "Docker removed"
                ;;
        esac
    fi
    
    # Uninstall golang-migrate
    if [ "$INSTALLED_MIGRATE" = true ]; then
        print_info "Removing golang-migrate..."
        if [ -f "/usr/local/bin/migrate" ]; then
            sudo rm -f /usr/local/bin/migrate
            print_success "golang-migrate removed"
        fi
    fi
    
    # Uninstall sqlc
    if [ "$INSTALLED_SQLC" = true ]; then
        print_info "Removing sqlc..."
        if [ -f "$HOME/go/bin/sqlc" ]; then
            rm -f "$HOME/go/bin/sqlc"
            print_success "sqlc removed"
        fi
    fi
    
    # Uninstall Go
    if [ "$INSTALLED_GO" = true ]; then
        print_info "Removing Go..."
        sudo rm -rf /usr/local/go
        
        # Remove Go environment variables from shell config
        SHELL_CONFIG=""
        if [ -n "$BASH_VERSION" ]; then
            SHELL_CONFIG="$HOME/.bashrc"
        elif [ -n "$ZSH_VERSION" ]; then
            SHELL_CONFIG="$HOME/.zshrc"
        else
            SHELL_CONFIG="$HOME/.profile"
        fi
        
        if [ -f "$SHELL_CONFIG" ]; then
            # Remove Go environment lines
            sed -i.bak '/# Go environment/d' "$SHELL_CONFIG" 2>/dev/null || true
            sed -i.bak '/export PATH=\$PATH:\/usr\/local\/go\/bin/d' "$SHELL_CONFIG" 2>/dev/null || true
            sed -i.bak '/export GOPATH=\$HOME\/go/d' "$SHELL_CONFIG" 2>/dev/null || true
            sed -i.bak '/export PATH=\$PATH:\$GOPATH\/bin/d' "$SHELL_CONFIG" 2>/dev/null || true
            rm -f "$SHELL_CONFIG.bak"
        fi
        
        print_success "Go removed"
    fi
    
    # Clean up temporary files
    rm -f go*.tar.gz migrate.tar.gz get-docker.sh 2>/dev/null || true
    
    # Kill sudo keepalive process if it exists
    if [ -n "${SUDO_KEEPALIVE_PID:-}" ]; then
        kill "$SUDO_KEEPALIVE_PID" 2>/dev/null || true
    fi
    
    echo ""
    print_success "Rollback completed. All installed dependencies have been removed."
    echo ""
}

# Error handler
handle_error() {
    local exit_code=$?
    if [ $exit_code -ne 0 ]; then
        print_error "Installation failed with exit code $exit_code"
        INSTALLATION_FAILED=true
        
        echo ""
        read -p "Do you want to rollback and uninstall all dependencies installed during this run? (y/N): " -n 1 -r
        echo
        if [[ $REPLY =~ ^[Yy]$ ]]; then
            uninstall_all
        else
            print_info "Keeping partial installation. You can manually run './install.sh uninstall' to clean up later."
        fi
        exit $exit_code
    fi
}

# Set up error trap
trap 'handle_error' ERR

# ================================
# Install Make
# ================================
install_make() {
    print_info "Checking for make..."
    if command_exists make; then
        print_success "make is already installed ($(make --version | head -n1))"
    else
        print_info "Installing make..."
        case "$OS" in
            Linux*)
                if command_exists apt-get; then
                    sudo apt-get update
                    sudo apt-get install -y build-essential
                elif command_exists yum; then
                    sudo yum groupinstall -y "Development Tools"
                elif command_exists dnf; then
                    sudo dnf groupinstall -y "Development Tools"
                elif command_exists pacman; then
                    sudo pacman -S --noconfirm base-devel
                else
                    print_error "Unable to install make. Please install manually."
                    return 1
                fi
                INSTALLED_MAKE=true
                ;;
            Darwin*)
                if command_exists brew; then
                    brew install make
                    INSTALLED_MAKE=true
                else
                    print_error "Homebrew not found. Please install Homebrew first."
                    return 1
                fi
                ;;
            *)
                print_error "Unsupported OS for automatic make installation"
                return 1
                ;;
        esac
        print_success "make installed successfully"
    fi
}

# ================================
# Install Golang
# ================================
install_golang() {
    print_info "Checking for Go..."
    if command_exists go; then
        CURRENT_GO_VERSION=$(go version | awk '{print $3}' | sed 's/go//')
        print_success "Go is already installed (version $CURRENT_GO_VERSION)"
        
        # Ask if user wants to update
        read -p "Do you want to reinstall Go $GO_VERSION? (y/N): " -n 1 -r
        echo
        if [[ ! $REPLY =~ ^[Yy]$ ]]; then
            return 0
        fi
    fi
    
    print_info "Installing Go $GO_VERSION..."
    
    # Determine the appropriate Go archive
    GO_ARCH=""
    case "$ARCH" in
        x86_64)
            GO_ARCH="amd64"
            ;;
        aarch64|arm64)
            GO_ARCH="arm64"
            ;;
        armv6l|armv7l)
            GO_ARCH="armv6l"
            ;;
        *)
            print_error "Unsupported architecture: $ARCH"
            return 1
            ;;
    esac
    
    GO_OS=""
    case "$OS" in
        Linux*)
            GO_OS="linux"
            ;;
        Darwin*)
            GO_OS="darwin"
            ;;
        *)
            print_error "Unsupported OS: $OS"
            return 1
            ;;
    esac
    
    GO_ARCHIVE="go${GO_VERSION}.${GO_OS}-${GO_ARCH}.tar.gz"
    GO_URL="https://go.dev/dl/${GO_ARCHIVE}"
    
    # Download Go
    print_info "Downloading Go from ${GO_URL}..."
    wget -q --show-progress "${GO_URL}" || {
        print_error "Failed to download Go. Please check your internet connection."
        return 1
    }
    
    # Remove old Go installation
    print_info "Removing old Go installation (if exists)..."
    sudo rm -rf /usr/local/go
    
    # Extract and install
    print_info "Installing Go to /usr/local/go..."
    sudo tar -C /usr/local -xzf "${GO_ARCHIVE}"
    rm "${GO_ARCHIVE}"
    
    # Set up environment variables
    print_info "Setting up environment variables..."
    
    # Determine shell config file
    SHELL_CONFIG=""
    if [ -n "$BASH_VERSION" ]; then
        SHELL_CONFIG="$HOME/.bashrc"
    elif [ -n "$ZSH_VERSION" ]; then
        SHELL_CONFIG="$HOME/.zshrc"
    else
        SHELL_CONFIG="$HOME/.profile"
    fi
    
    # Add Go to PATH if not already there
    if ! grep -q "/usr/local/go/bin" "$SHELL_CONFIG" 2>/dev/null; then
        echo "" >> "$SHELL_CONFIG"
        echo "# Go environment" >> "$SHELL_CONFIG"
        echo "export PATH=\$PATH:/usr/local/go/bin" >> "$SHELL_CONFIG"
        echo "export GOPATH=\$HOME/go" >> "$SHELL_CONFIG"
        echo "export PATH=\$PATH:\$GOPATH/bin" >> "$SHELL_CONFIG"
    fi
    
    # Apply changes to current session
    export PATH=$PATH:/usr/local/go/bin
    export GOPATH=$HOME/go
    export PATH=$PATH:$GOPATH/bin
    
    INSTALLED_GO=true
    
    print_success "Go $GO_VERSION installed successfully"
    go version
}

# ================================
# Install sqlc
# ================================
install_sqlc() {
    print_info "Checking for sqlc..."
    if command_exists sqlc; then
        print_success "sqlc is already installed ($(sqlc version))"
        return 0
    fi
    
    print_info "Installing sqlc $SQLC_VERSION..."
    
    # Install via Go
    if command_exists go; then
        go install github.com/sqlc-dev/sqlc/cmd/sqlc@latest
        INSTALLED_SQLC=true
        print_success "sqlc installed successfully"
        sqlc version
    else
        print_error "Go is required to install sqlc. Please install Go first."
        return 1
    fi
}

# ================================
# Install golang-migrate
# ================================
install_migrate() {
    print_info "Checking for migrate..."
    if command_exists migrate; then
        print_success "migrate is already installed ($(migrate -version 2>&1 | head -n1))"
        return 0
    fi
    
    print_info "Installing golang-migrate..."
    
    case "$OS" in
        Linux*)
            # Install using binary release
            MIGRATE_ARCH=""
            case "$ARCH" in
                x86_64)
                    MIGRATE_ARCH="amd64"
                    ;;
                aarch64|arm64)
                    MIGRATE_ARCH="arm64"
                    ;;
                *)
                    print_error "Unsupported architecture for migrate: $ARCH"
                    return 1
                    ;;
            esac
            
            MIGRATE_URL="https://github.com/golang-migrate/migrate/releases/download/v${MIGRATE_VERSION}/migrate.linux-${MIGRATE_ARCH}.tar.gz"
            print_info "Downloading migrate from ${MIGRATE_URL}..."
            
            wget -q --show-progress "${MIGRATE_URL}" -O migrate.tar.gz || {
                print_error "Failed to download migrate"
                return 1
            }
            
            tar -xzf migrate.tar.gz
            sudo mv migrate /usr/local/bin/migrate
            sudo chmod +x /usr/local/bin/migrate
            rm migrate.tar.gz
            INSTALLED_MIGRATE=true
            ;;
            
        Darwin*)
            if command_exists brew; then
                brew install golang-migrate
                INSTALLED_MIGRATE=true
            else
                print_error "Homebrew not found. Please install Homebrew first."
                return 1
            fi
            ;;
            
        *)
            print_error "Unsupported OS for migrate installation"
            return 1
            ;;
    esac
    
    print_success "golang-migrate installed successfully"
    migrate -version
}

# ================================
# Install Docker (if not present)
# ================================
install_docker() {
    print_info "Checking for Docker..."
    if command_exists docker; then
        print_success "Docker is already installed ($(docker --version))"
    else
        print_info "Docker is not installed."
        read -p "Do you want to install Docker? (y/N): " -n 1 -r
        echo
        if [[ $REPLY =~ ^[Yy]$ ]]; then
            case "$OS" in
                Linux*)
                    # Install Docker on Linux
                    curl -fsSL https://get.docker.com -o get-docker.sh
                    sudo sh get-docker.sh
                    rm get-docker.sh
                    INSTALLED_DOCKER=true
                    
                    # Add user to docker group
                    sudo usermod -aG docker $USER
                    print_success "Docker installed. Please log out and log back in for group changes to take effect."
                    ;;
                Darwin*)
                    print_info "Please install Docker Desktop from https://www.docker.com/products/docker-desktop"
                    ;;
                *)
                    print_error "Unsupported OS for Docker installation"
                    ;;
            esac
        fi
    fi
}

# ================================
# Install Docker Compose (if not present)
# ================================
install_docker_compose() {
    print_info "Checking for Docker Compose..."
    if command_exists docker-compose || docker compose version >/dev/null 2>&1; then
        print_success "Docker Compose is already installed"
    else
        print_info "Docker Compose is not installed."
        read -p "Do you want to install Docker Compose? (y/N): " -n 1 -r
        echo
        if [[ $REPLY =~ ^[Yy]$ ]]; then
            case "$OS" in
                Linux*)
                    DOCKER_COMPOSE_VERSION="2.23.3"
                    sudo curl -L "https://github.com/docker/compose/releases/download/v${DOCKER_COMPOSE_VERSION}/docker-compose-$(uname -s)-$(uname -m)" -o /usr/local/bin/docker-compose
                    sudo chmod +x /usr/local/bin/docker-compose
                    INSTALLED_DOCKER_COMPOSE=true
                    print_success "Docker Compose installed successfully"
                    ;;
                Darwin*)
                    print_info "Docker Compose comes with Docker Desktop on macOS"
                    ;;
                *)
                    print_error "Unsupported OS for Docker Compose installation"
                    ;;
            esac
        fi
    fi
}

# ================================
# Install Air (hot reload for Go)
# ================================
install_air() {
    print_info "Checking for Air (hot reload tool)..."
    if command_exists air; then
        print_success "Air is already installed"
    else
        print_info "Installing Air for hot reload..."
        if command_exists go; then
            go install github.com/air-verse/air@latest
            INSTALLED_AIR=true
            print_success "Air installed successfully"
        else
            print_error "Go is required to install Air"
            return 1
        fi
    fi
}

# ================================
# Install project Go dependencies
# ================================
install_go_dependencies() {
    print_info "Installing Go project dependencies..."
    if [ -f "go.mod" ]; then
        go mod download
        go mod tidy
        print_success "Go dependencies installed successfully"
    else
        print_error "go.mod not found in current directory"
        return 1
    fi
}

# ================================
# Verify installations
# ================================
verify_installations() {
    echo ""
    echo "=================================="
    echo "Verifying Installations"
    echo "=================================="
    echo ""
    
    local all_good=true
    
    # Check Make
    if command_exists make; then
        print_success "make: $(make --version | head -n1)"
    else
        print_error "make: Not installed"
        all_good=false
    fi
    
    # Check Go
    if command_exists go; then
        print_success "go: $(go version)"
    else
        print_error "go: Not installed"
        all_good=false
    fi
    
    # Check sqlc
    if command_exists sqlc; then
        print_success "sqlc: $(sqlc version)"
    else
        print_error "sqlc: Not installed"
        all_good=false
    fi
    
    # Check migrate
    if command_exists migrate; then
        print_success "migrate: $(migrate -version 2>&1 | head -n1)"
    else
        print_error "migrate: Not installed"
        all_good=false
    fi
    
    # Check Docker
    if command_exists docker; then
        print_success "docker: $(docker --version)"
    else
        print_info "docker: Not installed (optional)"
    fi
    
    # Check Air
    if command_exists air; then
        print_success "air: Installed"
    else
        print_info "air: Not installed (optional)"
    fi
    
    echo ""
    if [ "$all_good" = true ]; then
        print_success "All required tools are installed!"
        echo ""
        echo "Next steps:"
        echo "1. Copy .env.example to .env: cp .env.example .env"
        echo "2. Update database credentials in .env"
        echo "3. Start the development environment: make dev"
        echo "   OR start just the database: make noapp"
        echo "4. Run migrations: make migrate_up"
        echo "5. Generate sqlc code: make sqlc"
    else
        print_error "Some tools failed to install. Please check the errors above."
        return 1
    fi
}

# ================================
# Main installation flow
# ================================
main_install() {
    # Validate sudo access early (will prompt only if needed)
    print_info "Checking system requirements..."
    sudo -v 2>/dev/null && {
        # Keep sudo alive in background if we got it
        while true; do sudo -n true; sleep 50; kill -0 "$$" || exit; done 2>/dev/null &
        SUDO_KEEPALIVE_PID=$!
    }
    
    # Install dependencies
    install_make
    install_golang
    install_sqlc
    install_migrate
    install_docker
    install_docker_compose
    install_air
    
    # If we're in the project directory, install Go dependencies
    if [ -f "go.mod" ]; then
        install_go_dependencies
    fi
    
    # Verify everything
    verify_installations
    
    # Kill sudo keepalive process if it exists
    if [ -n "${SUDO_KEEPALIVE_PID:-}" ]; then
        kill "$SUDO_KEEPALIVE_PID" 2>/dev/null || true
    fi
    
    echo ""
    echo "=================================="
    echo "Installation Complete!"
    echo "=================================="
    echo ""
    print_info "Please restart your terminal or run: source ~/.bashrc (or ~/.zshrc)"
    echo ""
}

# ================================
# Main entry point
# ================================
main() {
    # Check for command line arguments
    if [ $# -gt 0 ]; then
        case "$1" in
            uninstall|remove|cleanup)
                # Set all flags to true for manual uninstall (remove everything)
                INSTALLED_MAKE=true
                INSTALLED_GO=true
                INSTALLED_SQLC=true
                INSTALLED_MIGRATE=true
                INSTALLED_DOCKER=true
                INSTALLED_DOCKER_COMPOSE=true
                INSTALLED_AIR=true
                uninstall_all
                exit 0
                ;;
            --help|-h|help)
                echo "Usage: $0 [COMMAND]"
                echo ""
                echo "Commands:"
                echo "  (no args)           Install all dependencies"
                echo "  uninstall           Uninstall all dependencies"
                echo "  help                Show this help message"
                echo ""
                exit 0
                ;;
            *)
                print_error "Unknown command: $1"
                echo "Run '$0 help' for usage information."
                exit 1
                ;;
        esac
    fi
    
    # Run installation
    main_install
}

# Run main function with all arguments
main "$@"
