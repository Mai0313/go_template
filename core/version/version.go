package version

import (
	"fmt"
	"runtime/debug"
	"strconv"
	"strings"
)

// These variables will be set at build time via -ldflags
var (
	// Version is the current version of the application
	Version = "dev"
	// BuildTime is when the binary was built
	BuildTime = "unknown"
	// GitCommit is the git commit hash
	GitCommit = "unknown"
)

// Info holds version information
type Info struct {
	Version   string `json:"version"`
	BuildTime string `json:"build_time"`
	GitCommit string `json:"git_commit"`
	GoVersion string `json:"go_version"`
}

// SemanticVersion represents a semantic version
type SemanticVersion struct {
	Major int
	Minor int
	Patch int
	Pre   string // pre-release identifier (e.g., "alpha.1", "beta.2")
	Build string // build metadata (e.g., "dirty", "20220101.abcdef")
}

// Get returns version information
func Get() Info {
	info := Info{
		Version:   Version,
		BuildTime: BuildTime,
		GitCommit: GitCommit,
		GoVersion: getGoVersion(),
	}

	// If version is still "dev", try to get it from build info (for go install)
	if info.Version == "dev" {
		if buildInfo, ok := debug.ReadBuildInfo(); ok {
			if buildInfo.Main.Version != "(devel)" && buildInfo.Main.Version != "" {
				info.Version = buildInfo.Main.Version
			}
		}
	}

	return info
}

// GetVersion returns just the version string
func GetVersion() string {
	return Get().Version
}

// IsNewerVersion checks whether a newer version is available (simplified)
func IsNewerVersion(current, latest string) bool {
	// Remove 'v' prefix
	current = strings.TrimPrefix(current, "v")
	latest = strings.TrimPrefix(latest, "v")

	// If current is "dev" or "unknown", treat as needing update
	if current == "dev" || current == "unknown" {
		return true
	}

	// Handle "dirty" versions (e.g., "0.2.3+dirty")
	if strings.Contains(current, "+") {
		current = strings.Split(current, "+")[0]
	}

	// If versions are equal, no update needed
	if current == latest {
		return false
	}

	// Try semantic comparison
	currentVer, err1 := ParseVersion(current)
	latestVer, err2 := ParseVersion(latest)

	if err1 == nil && err2 == nil {
		return latestVer.IsNewer(currentVer)
	}

	// If parsing fails, fall back to string comparison
	return current != latest
}

// ParseVersion parses a semantic version string
func ParseVersion(version string) (*SemanticVersion, error) {
	// Remove 'v' prefix
	version = strings.TrimPrefix(version, "v")

	if version == "dev" || version == "unknown" || version == "" {
		return &SemanticVersion{0, 0, 0, "dev", ""}, nil
	}

	sv := &SemanticVersion{}

	// Split out build metadata (after '+')
	if idx := strings.Index(version, "+"); idx != -1 {
		sv.Build = version[idx+1:]
		version = version[:idx]
	}

	// Split out pre-release (after '-')
	if idx := strings.Index(version, "-"); idx != -1 {
		sv.Pre = version[idx+1:]
		version = version[:idx]
	}

	// Parse major version
	parts := strings.Split(version, ".")
	if len(parts) < 1 {
		return nil, fmt.Errorf("invalid version format: %s", version)
	}

	var err error
	if sv.Major, err = strconv.Atoi(parts[0]); err != nil {
		return nil, fmt.Errorf("invalid major version: %s", parts[0])
	}

	// Parse minor version
	if len(parts) >= 2 {
		if sv.Minor, err = strconv.Atoi(parts[1]); err != nil {
			return nil, fmt.Errorf("invalid minor version: %s", parts[1])
		}
	}

	// Parse patch version
	if len(parts) >= 3 {
		if sv.Patch, err = strconv.Atoi(parts[2]); err != nil {
			return nil, fmt.Errorf("invalid patch version: %s", parts[2])
		}
	}

	return sv, nil
}

// Compare compares two versions
// Return: -1 (less), 0 (equal), 1 (greater)
func (sv *SemanticVersion) Compare(other *SemanticVersion) int {
	// Compare major
	if sv.Major != other.Major {
		if sv.Major < other.Major {
			return -1
		}
		return 1
	}

	// Compare minor
	if sv.Minor != other.Minor {
		if sv.Minor < other.Minor {
			return -1
		}
		return 1
	}

	// Compare patch
	if sv.Patch != other.Patch {
		if sv.Patch < other.Patch {
			return -1
		}
		return 1
	}

	// Compare pre-release
	if sv.Pre == "" && other.Pre != "" {
		return 1 // Release > pre-release
	}
	if sv.Pre != "" && other.Pre == "" {
		return -1 // Pre-release < release
	}
	if sv.Pre != other.Pre {
		return strings.Compare(sv.Pre, other.Pre)
	}

	// Build metadata does not affect precedence
	return 0
}

// IsNewer checks whether the current version is newer than another
func (sv *SemanticVersion) IsNewer(other *SemanticVersion) bool {
	return sv.Compare(other) > 0
}

func getGoVersion() string {
	if buildInfo, ok := debug.ReadBuildInfo(); ok {
		return buildInfo.GoVersion
	}
	return "unknown"
}
