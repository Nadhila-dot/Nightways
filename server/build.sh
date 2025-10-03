#!/bin/bash

# Build script for vela app across multiple platforms
# Requires CGO toolchains for cross-compilation

set -e

# Output directory
BUILD_DIR="./builds"
mkdir -p "$BUILD_DIR"

# Function to build for a specific platform
build_for_platform() {
    local GOOS=$1
    local GOARCH=$2
    local OUTPUT_DIR="$BUILD_DIR/$GOOS-$GOARCH"
    local BINARY_NAME="vela-app"
    
    if [ "$GOOS" = "windows" ]; then
        BINARY_NAME="vela-app.exe"
    fi
    
    mkdir -p "$OUTPUT_DIR"
    
    echo "Building for $GOOS/$GOARCH..."
    
    # Set environment variables for cross-compilation
    export GOOS=$GOOS
    export GOARCH=$GOARCH
    export CGO_ENABLED=1
    
    # Platform-specific C compiler settings
    if [ "$GOOS" = "windows" ]; then
        export CC="x86_64-w64-mingw32-gcc"
        export CXX="x86_64-w64-mingw32-g++"
    elif [ "$GOOS" = "linux" ]; then
        # Use Zig for Linux cross-compilation
        export CC="zig cc -target x86_64-linux-gnu"
        export CXX="zig c++ -target x86_64-linux-gnu"
    elif [ "$GOOS" = "darwin" ]; then
        export CC="clang"
        export CXX="clang++"
    fi
    
    # Build with appropriate tags for webview_go
    if go build -tags webkit2_41 -o "$OUTPUT_DIR/$BINARY_NAME" .; then
        echo "✅ Built successfully: $OUTPUT_DIR/$BINARY_NAME"
    else
        echo "❌ Failed to build for $GOOS/$GOARCH"
        return 1
    fi
}

# Build for each platform
echo "Starting multi-platform build..."

# macOS (Intel and Apple Silicon)
build_for_platform "darwin" "amd64"
build_for_platform "darwin" "arm64"

# Windows
build_for_platform "windows" "amd64"

echo "Build complete! Check the '$BUILD_DIR' directory for binaries."