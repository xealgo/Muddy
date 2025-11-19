#!/bin/bash

set -e  # Exit on any error

DESIRED_BUFFER_SIZE=7500000

# Detect operating system
detect_os() {
    case "$(uname -s)" in
        Linux*)     echo "Linux";;
        Darwin*)    echo "macOS";;
        CYGWIN*)    echo "Windows";;
        MINGW*)     echo "Windows";;
        MSYS*)      echo "Windows";;
        *)          echo "Unknown";;
    esac
}

OS=$(detect_os)

echo "🔧 Setting up UDP buffer limits for $OS..."

setup_linux() {
    echo "📊 Checking Linux UDP buffer limits..."
    
    # Get current values
    CURRENT_RMEM=$(sysctl -n net.core.rmem_max 2>/dev/null || echo "0")
    CURRENT_WMEM=$(sysctl -n net.core.wmem_max 2>/dev/null || echo "0")
    
    echo "Current rmem_max: $CURRENT_RMEM"
    echo "Current wmem_max: $CURRENT_WMEM"
    
    NEED_UPDATE=false
    
    if [ "$CURRENT_RMEM" -lt "$DESIRED_BUFFER_SIZE" ]; then
        echo "⚠️  rmem_max too small (need $DESIRED_BUFFER_SIZE)"
        NEED_UPDATE=true
    fi
    
    if [ "$CURRENT_WMEM" -lt "$DESIRED_BUFFER_SIZE" ]; then
        echo "⚠️  wmem_max too small (need $DESIRED_BUFFER_SIZE)"
        NEED_UPDATE=true
    fi
    
    if [ "$NEED_UPDATE" = true ]; then
        echo "🔐 Updating UDP buffer limits (requires sudo)..."
        sudo sysctl -w net.core.rmem_max=$DESIRED_BUFFER_SIZE
        sudo sysctl -w net.core.wmem_max=$DESIRED_BUFFER_SIZE
        echo "✅ Updated successfully!"
        echo ""
        echo "💡 To make these changes permanent, add to /etc/sysctl.conf:"
        echo "   net.core.rmem_max = $DESIRED_BUFFER_SIZE"
        echo "   net.core.wmem_max = $DESIRED_BUFFER_SIZE"
    else
        echo "✅ UDP buffer limits are already sufficient"
    fi
}

setup_macos() {
    echo "📊 Checking macOS UDP buffer limits..."
    
    # macOS uses different sysctl parameters
    CURRENT_MAX=$(sysctl -n kern.ipc.maxsockbuf 2>/dev/null || echo "0")
    
    echo "Current maxsockbuf: $CURRENT_MAX"
    
    if [ "$CURRENT_MAX" -lt "$DESIRED_BUFFER_SIZE" ]; then
        echo "⚠️  maxsockbuf too small (need $DESIRED_BUFFER_SIZE)"
        echo "🔐 Updating UDP buffer limits (requires sudo)..."
        sudo sysctl -w kern.ipc.maxsockbuf=$DESIRED_BUFFER_SIZE
        echo "✅ Updated successfully!"
        echo ""
        echo "💡 To make permanent, add to /etc/sysctl.conf:"
        echo "   kern.ipc.maxsockbuf = $DESIRED_BUFFER_SIZE"
    else
        echo "✅ UDP buffer limits are already sufficient"
    fi
}

setup_windows() {
    echo "📊 Windows UDP buffer management..."
    echo "ℹ️  Windows manages UDP buffers differently than Unix systems"
    echo "   Buffer sizes are typically set at the application level"
    echo "   Your Go application will handle buffer sizing automatically"
    echo ""
    echo "🪟 Windows-specific optimizations:"
    echo "   - Consider adjusting Windows network adapter settings"
    echo "   - Use 'netsh int udp set global' for system-wide UDP tuning"
    echo "   - Application-level buffer setting is usually sufficient"
    echo ""
    echo "✅ No system-level changes required for Windows"
}

# Main execution
case $OS in
    "Linux")
        setup_linux
        ;;
    "macOS")
        setup_macos
        ;;
    "Windows")
        setup_windows
        ;;
    "Unknown")
        echo "❌ Unsupported operating system: $(uname -s)"
        echo "   Manual UDP buffer configuration may be required"
        exit 1
        ;;
esac

echo ""
echo "🚀 UDP setup complete for $OS!"