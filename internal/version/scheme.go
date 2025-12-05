// Package version provides version parsing and bumping functionality.
package version

// Scheme represents a versioning scheme.
type Scheme string

const (
	// SchemeSemVer represents semantic versioning.
	SchemeSemVer Scheme = "semver"
	// SchemeCalVer represents calendar versioning.
	SchemeCalVer Scheme = "calver"
	// SchemeCustom represents a custom versioning scheme.
	SchemeCustom Scheme = "custom"
)

// Version is the interface that all version types must implement.
type Version interface {
	// String returns the string representation of the version.
	String() string
	// Bump increments the specified part of the version.
	Bump(part string) (Version, error)
}

// Parse parses a version string based on the scheme.
func Parse(version string, scheme Scheme) (Version, error) {
	switch scheme {
	case SchemeSemVer:
		return ParseSemVer(version)
	case SchemeCalVer:
		return ParseCalVer(version)
	default:
		// Default to SemVer
		return ParseSemVer(version)
	}
}
