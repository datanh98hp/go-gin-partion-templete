#!/bin/bash

set -e  # Exit on error

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

# Check if running as root
if [ "$EUID" -eq 0 ]; then 
    print_error "Please do not run this script as root"
    exit 1
fi

# Function to check if a command exists
command_exists() {
    command -v "$1" >/dev/null 2>&1
}

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
                ;;
            Darwin*)
                if command_exists brew; then
                    brew install make
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
            ;;
            
        Darwin*)
            if command_exists brew; then
                brew install golang-migrate
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
            go install github.com/cosmtrek/air@latest
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
main() {
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
    
    echo ""
    echo "=================================="
    echo "Installation Complete!"
    echo "=================================="
    echo ""
    print_info "Please restart your terminal or run: source ~/.bashrc (or ~/.zshrc)"
    echo ""
}

# Run main function
main
