#!/usr/bin/env bash
set -euo pipefail

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
BINARY_NAME="butter"
BINARY_PATH="/usr/local/bin/$BINARY_NAME"
EXTENSION_DIR="$SCRIPT_DIR/butter-extension"
VSIX_OUTPUT="$SCRIPT_DIR/butter-extension.vsix"

GREEN='\033[0;32m'
YELLOW='\033[1;33m'
RED='\033[0;31m'
NC='\033[0m'

info()  { echo -e "${GREEN}[INFO]${NC} $1"; }
warn()  { echo -e "${YELLOW}[WARN]${NC} $1"; }
error() { echo -e "${RED}[ERROR]${NC} $1"; }

install_binary() {
    info "Building Butter compiler binary..."
    
    if ! command -v go &>/dev/null; then
        error "Go is not installed. Please install Go from https://go.dev/dl/"
        exit 1
    fi
    
    cd "$SCRIPT_DIR"
    go build -o "$BINARY_NAME" main.go
    
    if [ ! -f "$BINARY_NAME" ]; then
        error "Build failed — no binary produced"
        exit 1
    fi
    
    INSTALL_DIR="/usr/local/bin"
    if [ ! -w "$INSTALL_DIR" ]; then
        INSTALL_DIR="$HOME/.local/bin"
        mkdir -p "$INSTALL_DIR"
    fi
    
    cp "$BINARY_NAME" "$INSTALL_DIR/$BINARY_NAME"
    chmod +x "$INSTALL_DIR/$BINARY_NAME"
    
    info "Butter compiler installed to $INSTALL_DIR/$BINARY_NAME"
    if [ -x "$INSTALL_DIR/$BINARY_NAME" ]; then
        info "$("$INSTALL_DIR/$BINARY_NAME" --version)"
    fi
    
    if [[ ":$PATH:" != *":$INSTALL_DIR:"* ]]; then
        warn "$INSTALL_DIR is not in PATH. Add it to your shell profile:"
        warn "  export PATH=\"\$PATH:$INSTALL_DIR\""
    fi
}

install_extension() {
    info "Installing VS Code Butter extension..."
    
    if ! command -v code &>/dev/null; then
        warn "VS Code 'code' CLI not found. Skipping extension installation."
        warn "To install manually, open butter-extension/ in VS Code and press F5."
        return
    fi
    
    if command -v vsce &>/dev/null; then
        cd "$EXTENSION_DIR"
        vsce package --out "$VSIX_OUTPUT" 2>/dev/null || true
        cd "$SCRIPT_DIR"
        
        if [ -f "$VSIX_OUTPUT" ]; then
            code --install-extension "$VSIX_OUTPUT" --force
            rm -f "$VSIX_OUTPUT"
            info "VS Code Butter extension installed successfully."
            return
        fi
    fi
    
    EXTENSIONS_DIR="${HOME}/.vscode/extensions/butter-extension"
    mkdir -p "$EXTENSIONS_DIR"
    cp -r "$EXTENSION_DIR"/* "$EXTENSIONS_DIR/"
    info "VS Code Butter extension copied to $EXTENSIONS_DIR"
    warn "Restart VS Code for the extension to take effect."
}

update_binary() {
    info "Updating Butter compiler binary..."
    install_binary
}

update_extension() {
    info "Updating VS Code Butter extension..."
    install_extension
}

case "${1:-install}" in
    install)
        install_binary
        install_extension
        info "Butter installation complete."
        ;;
    update)
        update_binary
        update_extension
        info "Butter update complete."
        ;;
    binary)
        install_binary
        ;;
    extension)
        install_extension
        ;;
    *)
        echo "Usage: $0 [install|update|binary|extension]"
        echo ""
        echo "  install   (default) Build and install compiler + VS Code extension"
        echo "  update    Rebuild compiler and reinstall extension"
        echo "  binary    Build and install only the compiler binary"
        echo "  extension Install only the VS Code extension"
        exit 1
        ;;
esac
