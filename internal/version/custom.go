// Package version provides custom version parsing and bumping functionality.
package version

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
	"time"
)

// CustomVersion represents a custom version string with parsed components
type CustomVersion struct {
	version string
}

// ParseCustom parses a custom version string based on a template.
// The template can contain placeholders like {major}, {minor}, {patch}, {year}, {month}, {day}.
func ParseCustom(version string, template string) (Version, error) {
	// For now, basic implementation just wraps the version
	// A more sophisticated implementation would parse the template and extract values
	return &CustomVersion{version: version}, nil
}

// String returns the string representation of the custom version.
func (cv *CustomVersion) String() string {
	return cv.version
}

// Bump increments the specified part of the custom version.
func (cv *CustomVersion) Bump(part string) (Version, error) {
	var bumpedVersion string

	switch part {
	case "major", "minor", "patch":
		// For SemVer-like parts, try to parse and increment numbers
		bumped, err := bumpSemverPart(cv.version, part)
		if err != nil {
			return nil, err
		}
		bumpedVersion = bumped
	case "date":
		// Replace date placeholders with current date
		now := time.Now()
		bumpedVersion = strings.Replace(cv.version, "YYYY", fmt.Sprintf("%04d", now.Year()), -1)
		bumpedVersion = strings.Replace(bumpedVersion, "MM", fmt.Sprintf("%02d", int(now.Month())), -1)
		bumpedVersion = strings.Replace(bumpedVersion, "DD", fmt.Sprintf("%02d", now.Day()), -1)
	case "build", "release":
		// Increment numeric suffix if it exists
		bumped, err := bumpNumericSuffix(cv.version)
		if err != nil {
			return nil, err
		}
		bumpedVersion = bumped
	default:
		// Try to handle other custom parts by looking for patterns in the version string
		return nil, fmt.Errorf("invalid custom version part: %s", part)
	}

	return &CustomVersion{version: bumpedVersion}, nil
}

// bumpSemverPart attempts to bump major/minor/patch from a version string
func bumpSemverPart(version string, part string) (string, error) {
	// Look for patterns like x.y.z and increment the appropriate part
	re := regexp.MustCompile(`(\d+)\.(\d+)\.(\d+)`)
	matches := re.FindStringSubmatch(version)
	if matches == nil {
		return "", fmt.Errorf("version does not contain semver pattern x.y.z: %s", version)
	}

	// Parse the numbers
	major, _ := strconv.Atoi(matches[1])
	minor, _ := strconv.Atoi(matches[2])
	patch, _ := strconv.Atoi(matches[3])

	// Increment the appropriate part
	switch part {
	case "major":
		major++
		minor = 0
		patch = 0
	case "minor":
		minor++
		patch = 0
	case "patch":
		patch++
	}

	// Replace in the version string
	newVersion := re.ReplaceAllString(version, fmt.Sprintf("%d.%d.%d", major, minor, patch))
	return newVersion, nil
}

// bumpNumericSuffix increments a numeric suffix (like x.y.z.1 to x.y.z.2)
func bumpNumericSuffix(version string) (string, error) {
	// Look for patterns like \.\d+$ (a dot followed by digits at the end)
	re := regexp.MustCompile(`(\.)(\d+)$`)
	matches := re.FindStringSubmatch(version)
	if matches == nil {
		// If no numeric suffix found, append .1
		return version + ".1", nil
	}

	// Parse the number and increment it
	num, err := strconv.Atoi(matches[2])
	if err != nil {
		return "", fmt.Errorf("could not parse numeric suffix: %v", err)
	}
	
	newNum := num + 1
	newVersion := re.ReplaceAllString(version, fmt.Sprintf(".%d", newNum))
	return newVersion, nil
}