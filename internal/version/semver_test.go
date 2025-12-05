package version

import (
	"testing"
)

func TestParseSemVer(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		want    *SemVer
		wantErr bool
	}{
		{
			name:  "simple version",
			input: "1.2.3",
			want:  &SemVer{Major: 1, Minor: 2, Patch: 3},
		},
		{
			name:  "with v prefix",
			input: "v1.2.3",
			want:  &SemVer{Major: 1, Minor: 2, Patch: 3},
		},
		{
			name:  "with prerelease",
			input: "1.2.3-alpha.1",
			want:  &SemVer{Major: 1, Minor: 2, Patch: 3, Prerelease: "alpha.1"},
		},
		{
			name:  "with build",
			input: "1.2.3+build.123",
			want:  &SemVer{Major: 1, Minor: 2, Patch: 3, Build: "build.123"},
		},
		{
			name:  "with prerelease and build",
			input: "1.2.3-alpha.1+build.123",
			want:  &SemVer{Major: 1, Minor: 2, Patch: 3, Prerelease: "alpha.1", Build: "build.123"},
		},
		{
			name:    "invalid version",
			input:   "invalid",
			wantErr: true,
		},
		{
			name:    "incomplete version",
			input:   "1.2",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ParseSemVer(tt.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("ParseSemVer() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if err != nil {
				return
			}
			if got.Major != tt.want.Major || got.Minor != tt.want.Minor || got.Patch != tt.want.Patch {
				t.Errorf("ParseSemVer() = %v, want %v", got, tt.want)
			}
			if got.Prerelease != tt.want.Prerelease {
				t.Errorf("ParseSemVer() prerelease = %v, want %v", got.Prerelease, tt.want.Prerelease)
			}
			if got.Build != tt.want.Build {
				t.Errorf("ParseSemVer() build = %v, want %v", got.Build, tt.want.Build)
			}
		})
	}
}

func TestSemVerString(t *testing.T) {
	tests := []struct {
		name   string
		semver *SemVer
		want   string
	}{
		{
			name:   "simple version",
			semver: &SemVer{Major: 1, Minor: 2, Patch: 3},
			want:   "1.2.3",
		},
		{
			name:   "with prerelease",
			semver: &SemVer{Major: 1, Minor: 2, Patch: 3, Prerelease: "alpha.1"},
			want:   "1.2.3-alpha.1",
		},
		{
			name:   "with build",
			semver: &SemVer{Major: 1, Minor: 2, Patch: 3, Build: "build.123"},
			want:   "1.2.3+build.123",
		},
		{
			name:   "with prerelease and build",
			semver: &SemVer{Major: 1, Minor: 2, Patch: 3, Prerelease: "alpha.1", Build: "build.123"},
			want:   "1.2.3-alpha.1+build.123",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.semver.String(); got != tt.want {
				t.Errorf("SemVer.String() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestSemVerBump(t *testing.T) {
	tests := []struct {
		name    string
		semver  *SemVer
		part    string
		want    string
		wantErr bool
	}{
		{
			name:   "bump major",
			semver: &SemVer{Major: 1, Minor: 2, Patch: 3},
			part:   "major",
			want:   "2.0.0",
		},
		{
			name:   "bump minor",
			semver: &SemVer{Major: 1, Minor: 2, Patch: 3},
			part:   "minor",
			want:   "1.3.0",
		},
		{
			name:   "bump patch",
			semver: &SemVer{Major: 1, Minor: 2, Patch: 3},
			part:   "patch",
			want:   "1.2.4",
		},
		{
			name:   "bump prerelease from empty",
			semver: &SemVer{Major: 1, Minor: 2, Patch: 3},
			part:   "pre",
			want:   "1.2.3-alpha.1",
		},
		{
			name:   "bump prerelease with number",
			semver: &SemVer{Major: 1, Minor: 2, Patch: 3, Prerelease: "alpha.1"},
			part:   "pre",
			want:   "1.2.3-alpha.2",
		},
		{
			name:   "bump build from empty",
			semver: &SemVer{Major: 1, Minor: 2, Patch: 3},
			part:   "build",
			want:   "1.2.3+1",
		},
		{
			name:   "bump build with number",
			semver: &SemVer{Major: 1, Minor: 2, Patch: 3, Build: "5"},
			part:   "build",
			want:   "1.2.3+6",
		},
		{
			name:    "invalid part",
			semver:  &SemVer{Major: 1, Minor: 2, Patch: 3},
			part:    "invalid",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.semver.Bump(tt.part)
			if (err != nil) != tt.wantErr {
				t.Errorf("SemVer.Bump() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if err != nil {
				return
			}
			if got.String() != tt.want {
				t.Errorf("SemVer.Bump() = %v, want %v", got.String(), tt.want)
			}
		})
	}
}
