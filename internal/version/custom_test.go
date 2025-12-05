package version

import (
	"testing"
)

func TestParseCustom(t *testing.T) {
	tests := []struct {
		name     string
		version  string
		template string
		wantStr  string
		wantErr  bool
	}{
		{
			name:     "simple version",
			version:  "1.2.3",
			template: "",
			wantStr:  "1.2.3",
		},
		{
			name:     "version with prefix",
			version:  "v1.2.3",
			template: "",
			wantStr:  "v1.2.3",
		},
		{
			name:     "version with suffix",
			version:  "1.2.3-beta",
			template: "",
			wantStr:  "1.2.3-beta",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ParseCustom(tt.version, tt.template)
			if (err != nil) != tt.wantErr {
				t.Errorf("ParseCustom() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if err != nil {
				return
			}
			if got.String() != tt.wantStr {
				t.Errorf("ParseCustom() = %v, want %v", got.String(), tt.wantStr)
			}
		})
	}
}

func TestCustomVersionString(t *testing.T) {
	tests := []struct {
		name    string
		version string
		want    string
	}{
		{
			name:    "simple version",
			version: "1.2.3",
			want:    "1.2.3",
		},
		{
			name:    "version with prefix",
			version: "release-1.0.0",
			want:    "release-1.0.0",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cv := &CustomVersion{version: tt.version}
			if got := cv.String(); got != tt.want {
				t.Errorf("CustomVersion.String() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestCustomVersionBump(t *testing.T) {
	tests := []struct {
		name    string
		version string
		part    string
		want    string
		wantErr bool
	}{
		{
			name:    "bump major",
			version: "1.2.3",
			part:    "major",
			want:    "2.0.0",
		},
		{
			name:    "bump minor",
			version: "1.2.3",
			part:    "minor",
			want:    "1.3.0",
		},
		{
			name:    "bump patch",
			version: "1.2.3",
			part:    "patch",
			want:    "1.2.4",
		},
		{
			name:    "bump major with prefix",
			version: "v1.2.3",
			part:    "major",
			want:    "v2.0.0",
		},
		{
			name:    "bump build increments last numeric suffix",
			version: "1.2.3",
			part:    "build",
			want:    "1.2.4", // The last .3 is treated as suffix
		},
		{
			name:    "bump build with explicit suffix",
			version: "1.2.3.5",
			part:    "build",
			want:    "1.2.3.6",
		},
		{
			name:    "bump release alias",
			version: "1.0.0",
			part:    "release",
			want:    "1.0.1", // The last .0 is treated as suffix
		},
		{
			name:    "invalid part",
			version: "1.2.3",
			part:    "invalid",
			wantErr: true,
		},
		{
			name:    "bump major on non-semver",
			version: "release",
			part:    "major",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cv := &CustomVersion{version: tt.version}
			got, err := cv.Bump(tt.part)
			if (err != nil) != tt.wantErr {
				t.Errorf("CustomVersion.Bump() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if err != nil {
				return
			}
			if got.String() != tt.want {
				t.Errorf("CustomVersion.Bump() = %v, want %v", got.String(), tt.want)
			}
		})
	}
}

func TestBumpSemverPart(t *testing.T) {
	tests := []struct {
		name    string
		version string
		part    string
		want    string
		wantErr bool
	}{
		{
			name:    "bump major",
			version: "1.2.3",
			part:    "major",
			want:    "2.0.0",
		},
		{
			name:    "bump minor",
			version: "1.2.3",
			part:    "minor",
			want:    "1.3.0",
		},
		{
			name:    "bump patch",
			version: "1.2.3",
			part:    "patch",
			want:    "1.2.4",
		},
		{
			name:    "version with prefix",
			version: "version-1.2.3",
			part:    "patch",
			want:    "version-1.2.4",
		},
		{
			name:    "no semver pattern",
			version: "release",
			part:    "major",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := bumpSemverPart(tt.version, tt.part)
			if (err != nil) != tt.wantErr {
				t.Errorf("bumpSemverPart() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if err != nil {
				return
			}
			if got != tt.want {
				t.Errorf("bumpSemverPart() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestBumpNumericSuffix(t *testing.T) {
	tests := []struct {
		name    string
		version string
		want    string
		wantErr bool
	}{
		{
			name:    "semver-like version",
			version: "1.2.3",
			want:    "1.2.4", // .3 is treated as suffix
		},
		{
			name:    "with explicit suffix",
			version: "1.2.3.5",
			want:    "1.2.3.6",
		},
		{
			name:    "suffix at 99",
			version: "1.0.0.99",
			want:    "1.0.0.100",
		},
		{
			name:    "no dots",
			version: "release",
			want:    "release.1",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := bumpNumericSuffix(tt.version)
			if (err != nil) != tt.wantErr {
				t.Errorf("bumpNumericSuffix() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if err != nil {
				return
			}
			if got != tt.want {
				t.Errorf("bumpNumericSuffix() = %v, want %v", got, tt.want)
			}
		})
	}
}
