// Package version provides version parsing and bumping functionality.
package version

import (
	"fmt"
	"regexp"
	"strconv"
	"time"
)

// CalVer represents a calendar version.
// Format: YYYY.MM.DD or YYYY.MM.DD.N
type CalVer struct {
	Year    int
	Month   int
	Day     int
	Release int // Optional release number for multiple releases per day
}

var calverRegex = regexp.MustCompile(`^v?(\d{4})\.(\d{1,2})\.(\d{1,2})(?:\.(\d+))?$`)

// ParseCalVer parses a calendar version string into a CalVer struct.
func ParseCalVer(version string) (*CalVer, error) {
	matches := calverRegex.FindStringSubmatch(version)
	if matches == nil {
		return nil, fmt.Errorf("invalid calendar version: %s", version)
	}

	year, _ := strconv.Atoi(matches[1])
	month, _ := strconv.Atoi(matches[2])
	day, _ := strconv.Atoi(matches[3])
	release := 0
	if matches[4] != "" {
		release, _ = strconv.Atoi(matches[4])
	}

	return &CalVer{
		Year:    year,
		Month:   month,
		Day:     day,
		Release: release,
	}, nil
}

// String returns the string representation of the calendar version.
func (v *CalVer) String() string {
	if v.Release > 0 {
		return fmt.Sprintf("%d.%02d.%02d.%d", v.Year, v.Month, v.Day, v.Release)
	}
	return fmt.Sprintf("%d.%02d.%02d", v.Year, v.Month, v.Day)
}

// Bump increments the specified part of the calendar version.
// Valid parts are: date, build/release.
func (v *CalVer) Bump(part string) (Version, error) {
	switch part {
	case "date":
		now := time.Now()
		newVersion := &CalVer{
			Year:    now.Year(),
			Month:   int(now.Month()),
			Day:     now.Day(),
			Release: 0,
		}
		// If same date, increment release
		if v.Year == newVersion.Year && v.Month == newVersion.Month && v.Day == newVersion.Day {
			newVersion.Release = v.Release + 1
		}
		return newVersion, nil
	case "build", "release":
		newVersion := &CalVer{
			Year:    v.Year,
			Month:   v.Month,
			Day:     v.Day,
			Release: v.Release + 1,
		}
		return newVersion, nil
	default:
		return nil, fmt.Errorf("invalid CalVer part: %s (valid: date, build, release)", part)
	}
}

// Copy creates a deep copy of the CalVer.
func (v *CalVer) Copy() *CalVer {
	return &CalVer{
		Year:    v.Year,
		Month:   v.Month,
		Day:     v.Day,
		Release: v.Release,
	}
}
