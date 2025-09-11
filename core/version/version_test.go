package version

import (
	"strings"
	"testing"
)

func TestParseVersion_BasicAndPreRelease(t *testing.T) {
	cases := []struct {
		in                  string
		major, minor, patch int
		pre, build          string
	}{
		{"v1.2.3", 1, 2, 3, "", ""},
		{"1.2.3", 1, 2, 3, "", ""},
		{"1.2.3-alpha.1+build.123", 1, 2, 3, "alpha.1", "build.123"},
		{"1.2.3+dirty", 1, 2, 3, "", "dirty"},
	}
	for _, tc := range cases {
		v, err := ParseVersion(tc.in)
		if err != nil {
			t.Fatalf("ParseVersion(%q) unexpected error: %v", tc.in, err)
		}
		if v.Major != tc.major || v.Minor != tc.minor || v.Patch != tc.patch || v.Pre != tc.pre || v.Build != tc.build {
			t.Errorf("ParseVersion(%q) => %+v", tc.in, v)
		}
	}
}

func TestParseVersion_SpecialAndErrors(t *testing.T) {
	// special values
	for _, in := range []string{"dev", "unknown", ""} {
		v, err := ParseVersion(in)
		if err != nil {
			t.Fatalf("ParseVersion(%q) error: %v", in, err)
		}
		if v.Major != 0 || v.Pre != "dev" {
			t.Errorf("ParseVersion(%q) => %+v, want Pre=dev", in, v)
		}
	}
	// errors
	bad := []string{"x.1.2", "1.x.2", "1.2.x"}
	for _, in := range bad {
		if _, err := ParseVersion(in); err == nil {
			t.Errorf("ParseVersion(%q) expected error", in)
		}
	}
}

func TestCompare_And_IsNewer(t *testing.T) {
	mk := func(s string) *SemanticVersion { v, _ := ParseVersion(s); return v }
	if mk("1.2.3").Compare(mk("1.2.2")) <= 0 {
		t.Error("1.2.3 should be greater than 1.2.2")
	}
	if mk("1.0.0-alpha").Compare(mk("1.0.0")) >= 0 {
		t.Error("pre-release should be less than release")
	}
	if mk("1.0.0-alpha.2").Compare(mk("1.0.0-alpha.1")) <= 0 {
		t.Error("alpha.2 should be greater than alpha.1 lexicographically")
	}
	// build metadata should not affect ordering
	if mk("1.2.3+dirty").Compare(mk("1.2.3")) != 0 {
		t.Error("build metadata must be ignored in comparison")
	}
	// IsNewer wrapper
	if !mk("1.0.1").IsNewer(mk("1.0.0")) {
		t.Error("1.0.1 should be newer than 1.0.0")
	}
}

func TestIsNewerVersion_Helper(t *testing.T) {
	// dev should prompt update
	if !IsNewerVersion("dev", "v1.0.0") {
		t.Error("dev current should indicate update available")
	}
	// equal versions
	if IsNewerVersion("v1.2.3", "1.2.3") {
		t.Error("equal versions should not indicate update")
	}
	// dirty build should be treated as same base version
	if IsNewerVersion("1.2.3+dirty", "1.2.3") {
		t.Error("dirty build should not indicate update vs same version")
	}
	// non-semver fallback: different strings => update
	if !IsNewerVersion("feature-branch", "1.0.0") {
		t.Error("non-semver comparison should fallback to string inequality")
	}
}

func TestGetAndGetVersion_ReflectsVariables(t *testing.T) {
	oldV, oldT, oldC := Version, BuildTime, GitCommit
	defer func() {
		Version, BuildTime, GitCommit = oldV, oldT, oldC
	}()
	Version = "1.2.3"
	BuildTime = "2025-09-02T00:00:00Z"
	GitCommit = "abcdef"
	info := Get()
	if info.Version != Version || info.BuildTime != BuildTime || info.GitCommit != GitCommit {
		t.Errorf("Get() mismatch: %+v", info)
	}
	if got := GetVersion(); got != Version {
		t.Errorf("GetVersion()=%s, want %s", got, Version)
	}
	if !strings.HasPrefix(info.GoVersion, "go") && info.GoVersion != "unknown" {
		t.Errorf("unexpected go version: %s", info.GoVersion)
	}
}
