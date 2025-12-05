// Package version provides version parsing and bumping functionality.
package version

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
)

// SemVer represents a semantic version following the SemVer 2.0.0 specification.
type SemVer struct {
	Major      int
	Minor      int
	Patch      int
	Prerelease string
	Build      string
}

var semverRegex = regexp.MustCompile(`^v?(\d+)\.(\d+)\.(\d+)(?:-([0-9A-Za-z\-\.]+))?(?:\+([0-9A-Za-z\-\.]+))?$`)

// ParseSemVer parses a version string into a SemVer struct.
func ParseSemVer(version string) (*SemVer, error) {
	matches := semverRegex.FindStringSubmatch(version)
	if matches == nil {
		return nil, fmt.Errorf("invalid semantic version: %s", version)
	}

	major, _ := strconv.Atoi(matches[1])
	minor, _ := strconv.Atoi(matches[2])
	patch, _ := strconv.Atoi(matches[3])

	return &SemVer{
		Major:      major,
		Minor:      minor,
		Patch:      patch,
		Prerelease: matches[4],
		Build:      matches[5],
	}, nil
}

// String returns the string representation of the version.
func (v *SemVer) String() string {
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%d.%d.%d", v.Major, v.Minor, v.Patch))

	if v.Prerelease != "" {
		sb.WriteString("-")
		sb.WriteString(v.Prerelease)
	}

	if v.Build != "" {
		sb.WriteString("+")
		sb.WriteString(v.Build)
	}

	return sb.String()
}

// Bump increments the specified part of the version.
// Valid parts are: major, minor, patch, pre, build.
func (v *SemVer) Bump(part string) (Version, error) {
	newVersion := &SemVer{
		Major:      v.Major,
		Minor:      v.Minor,
		Patch:      v.Patch,
		Prerelease: v.Prerelease,
		Build:      v.Build,
	}

	switch strings.ToLower(part) {
	case "major":
		newVersion.Major++
		newVersion.Minor = 0
		newVersion.Patch = 0
		newVersion.Prerelease = ""
		newVersion.Build = ""
	case "minor":
		newVersion.Minor++
		newVersion.Patch = 0
		newVersion.Prerelease = ""
		newVersion.Build = ""
	case "patch":
		newVersion.Patch++
		newVersion.Prerelease = ""
		newVersion.Build = ""
	case "pre", "prerelease":
		newVersion.Prerelease = bumpPrerelease(newVersion.Prerelease)
	case "build":
		newVersion.Build = bumpBuild(newVersion.Build)
	default:
		return nil, fmt.Errorf("invalid version part: %s (valid: major, minor, patch, pre, build)", part)
	}

	return newVersion, nil
}

// bumpPrerelease increments the prerelease version.
func bumpPrerelease(prerelease string) string {
	if prerelease == "" {
		return "alpha.1"
	}

	// Try to find a numeric suffix to increment
	parts := strings.Split(prerelease, ".")
	if len(parts) > 0 {
		lastPart := parts[len(parts)-1]
		if num, err := strconv.Atoi(lastPart); err == nil {
			parts[len(parts)-1] = strconv.Itoa(num + 1)
			return strings.Join(parts, ".")
		}
	}

	// If no numeric suffix, append .1
	return prerelease + ".1"
}

// bumpBuild increments the build metadata.
func bumpBuild(build string) string {
	if build == "" {
		return "1"
	}

	// Try to parse as a number
	if num, err := strconv.Atoi(build); err == nil {
		return strconv.Itoa(num + 1)
	}

	// Try to find a numeric suffix to increment
	parts := strings.Split(build, ".")
	if len(parts) > 0 {
		lastPart := parts[len(parts)-1]
		if num, err := strconv.Atoi(lastPart); err == nil {
			parts[len(parts)-1] = strconv.Itoa(num + 1)
			return strings.Join(parts, ".")
		}
	}

	// If no numeric suffix, append .1
	return build + ".1"
}

// Copy creates a deep copy of the SemVer.
func (v *SemVer) Copy() *SemVer {
	return &SemVer{
		Major:      v.Major,
		Minor:      v.Minor,
		Patch:      v.Patch,
		Prerelease: v.Prerelease,
		Build:      v.Build,
	}
}
