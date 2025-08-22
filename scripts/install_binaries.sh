#!/bin/bash

set -e

GCS_BUCKET="https://storage.googleapis.com/claude-code-dist-86c565f3-f756-42ad-8dfa-d59b1c096819/claude-code-releases"
BASE_DOWNLOAD_DIR="$(pwd)/binaries"

# Function to get version
get_version() {
    local target="${1:-latest}"
    
    # Validate target if provided
    if [[ ! "$target" =~ ^(stable|latest|[0-9]+\.[0-9]+\.[0-9]+(-[^[:space:]]+)?)$ ]]; then
        echo "Error: Invalid target '$target'. Use stable|latest|VERSION" >&2
        return 1
    fi
    
    # Get version from GCS bucket
    curl -fsSLk "$GCS_BUCKET/$target"
}

# Check if this script is called for version only
if [ "$1" = "get_version" ]; then
    get_version "$2"
    exit $?
fi

# Parse command line arguments for install mode
TARGET="${1:-latest}"  # Default to latest if not provided

# Validate target if provided
if [[ ! "$TARGET" =~ ^(stable|latest|[0-9]+\.[0-9]+\.[0-9]+(-[^[:space:]]+)?)$ ]]; then
    echo "Usage: $0 [stable|latest|VERSION] or $0 get_version [stable|latest|VERSION]" >&2
    exit 1
fi

# Get version
version=$(get_version "$TARGET")
echo "Downloading version: $version"

DOWNLOAD_DIR="$BASE_DOWNLOAD_DIR/$version"
mkdir -p "$DOWNLOAD_DIR"

# Check for required dependencies
if ! command -v curl >/dev/null 2>&1; then
    echo "curl is required but not installed" >&2
    exit 1
fi

# Check if jq is available (optional)
HAS_JQ=false
if command -v jq >/dev/null 2>&1; then
    HAS_JQ=true
fi

# Simple JSON parser for extracting checksum when jq is not available
get_checksum_from_manifest() {
    local json="$1"
    local platform="$2"
    # Normalize JSON to single line and extract checksum
    json=$(echo "$json" | tr -d '\n\r\t' | sed 's/ \+/ /g')
    # Extract checksum for platform using bash regex
    if [[ $json =~ \"$platform\"[^}]*\"checksum\"[[:space:]]*:[[:space:]]*\"([a-f0-9]{64})\" ]]; then
        echo "${BASH_REMATCH[1]}"
        return 0
    fi
    return 1
}

# Define all supported platforms with their file extensions
declare -A platforms=(
    ["darwin-x64"]="claude"
    ["darwin-arm64"]="claude"
    ["linux-x64"]="claude"
    ["linux-arm64"]="claude"
    ["win32-x64"]="claude.exe"
)

# Download manifest once
echo "Fetching manifest..."
manifest_json=$(curl -fsSLk "$GCS_BUCKET/$version/manifest.json")

# Function to download and verify a platform binary
download_platform() {
    local platform=$1
    local manifest_json=$2
    local HAS_JQ=$3

    echo "Processing platform: $platform"

    # Extract checksum for current platform
    if [ "$HAS_JQ" = true ]; then
        checksum=$(echo "$manifest_json" | jq -r ".platforms[\"$platform\"].checksum // empty")
    else
        checksum=$(get_checksum_from_manifest "$manifest_json" "$platform")
    fi

    # Validate checksum format (SHA256 = 64 hex characters)
    if [ -z "$checksum" ] || [[ ! "$checksum" =~ ^[a-f0-9]{64}$ ]]; then
        echo "Warning: Platform $platform not found in manifest, skipping..." >&2
        return 1
    fi

    # Get the correct filename and extension for this platform
    filename="${platforms[$platform]}"

    # Set binary path with appropriate extension
    if [[ "$platform" == win32-* ]]; then
        binary_path="$DOWNLOAD_DIR/claude-$version-$platform.exe"
    else
        binary_path="$DOWNLOAD_DIR/claude-$version-$platform"
    fi

    # Download binary
    echo "Downloading $platform binary..."

    if ! curl -fsSLk -o "$binary_path" "$GCS_BUCKET/$version/$platform/$filename"; then
        echo "Warning: Download failed for $platform, skipping..." >&2
        rm -f "$binary_path"
        return 1
    fi

    # Verify checksum
    echo "Verifying checksum for $platform..."

    # Pick the right checksum tool
    if command -v sha256sum >/dev/null 2>&1; then
        actual=$(sha256sum "$binary_path" | cut -d' ' -f1)
    elif command -v shasum >/dev/null 2>&1; then
        actual=$(shasum -a 256 "$binary_path" | cut -d' ' -f1)
    else
        echo "Warning: No SHA256 tool found, skipping checksum verification for $platform" >&2
        chmod +x "$binary_path" 2>/dev/null || true  # chmod may fail on some filesystems for .exe files
        echo "Downloaded (checksum not verified): $binary_path"
        return 0
    fi

    if [ "$actual" != "$checksum" ]; then
        echo "Warning: Checksum verification failed for $platform, removing file..." >&2
        rm -f "$binary_path"
        return 1
    fi

    # Set executable permissions (will be ignored for Windows .exe files on non-Windows systems)
    chmod +x "$binary_path" 2>/dev/null || true
    echo "Successfully downloaded and verified: $binary_path"
    return 0
}

# Download binaries for all platforms in parallel
echo "Starting parallel downloads for all platforms..."
for platform in "${!platforms[@]}"; do
    download_platform "$platform" "$manifest_json" "$HAS_JQ" &
done

# Wait for all background processes to finish
wait
echo "All downloads completed!"

echo ""
echo "Download completed!"
echo "Files saved to: $DOWNLOAD_DIR"
echo ""
ls -la "$DOWNLOAD_DIR" | grep claude || echo "No claude files found"
