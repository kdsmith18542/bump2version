package version

import (
	"testing"
	"time"
)

func TestParseCalVer(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		want    *CalVer
		wantErr bool
	}{
		{
			name:  "simple date",
			input: "2024.01.15",
			want:  &CalVer{Year: 2024, Month: 1, Day: 15},
		},
		{
			name:  "with v prefix",
			input: "v2024.01.15",
			want:  &CalVer{Year: 2024, Month: 1, Day: 15},
		},
		{
			name:  "with release number",
			input: "2024.01.15.2",
			want:  &CalVer{Year: 2024, Month: 1, Day: 15, Release: 2},
		},
		{
			name:  "single digit month and day",
			input: "2024.1.5",
			want:  &CalVer{Year: 2024, Month: 1, Day: 5},
		},
		{
			name:    "invalid format",
			input:   "invalid",
			wantErr: true,
		},
		{
			name:    "incomplete date",
			input:   "2024.01",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ParseCalVer(tt.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("ParseCalVer() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if err != nil {
				return
			}
			if got.Year != tt.want.Year || got.Month != tt.want.Month || got.Day != tt.want.Day {
				t.Errorf("ParseCalVer() date = %d.%d.%d, want %d.%d.%d",
					got.Year, got.Month, got.Day, tt.want.Year, tt.want.Month, tt.want.Day)
			}
			if got.Release != tt.want.Release {
				t.Errorf("ParseCalVer() release = %d, want %d", got.Release, tt.want.Release)
			}
		})
	}
}

func TestCalVerString(t *testing.T) {
	tests := []struct {
		name   string
		calver *CalVer
		want   string
	}{
		{
			name:   "simple date",
			calver: &CalVer{Year: 2024, Month: 1, Day: 15},
			want:   "2024.01.15",
		},
		{
			name:   "with release",
			calver: &CalVer{Year: 2024, Month: 1, Day: 15, Release: 2},
			want:   "2024.01.15.2",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.calver.String(); got != tt.want {
				t.Errorf("CalVer.String() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestCalVerBump(t *testing.T) {
	now := time.Now()
	today := &CalVer{Year: now.Year(), Month: int(now.Month()), Day: now.Day()}

	tests := []struct {
		name      string
		calver    *CalVer
		part      string
		checkDate bool
		wantErr   bool
	}{
		{
			name:      "bump date",
			calver:    &CalVer{Year: 2024, Month: 1, Day: 15},
			part:      "date",
			checkDate: true,
		},
		{
			name:   "bump release",
			calver: &CalVer{Year: 2024, Month: 1, Day: 15, Release: 0},
			part:   "release",
		},
		{
			name:   "bump build (alias)",
			calver: &CalVer{Year: 2024, Month: 1, Day: 15, Release: 1},
			part:   "build",
		},
		{
			name:    "invalid part",
			calver:  &CalVer{Year: 2024, Month: 1, Day: 15},
			part:    "invalid",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.calver.Bump(tt.part)
			if (err != nil) != tt.wantErr {
				t.Errorf("CalVer.Bump() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if err != nil {
				return
			}
			gotCalVer := got.(*CalVer)
			if tt.checkDate {
				if gotCalVer.Year != today.Year || gotCalVer.Month != today.Month || gotCalVer.Day != today.Day {
					t.Errorf("CalVer.Bump(date) = %v, want today's date", got.String())
				}
			}
			if tt.part == "release" || tt.part == "build" {
				if gotCalVer.Release != tt.calver.Release+1 {
					t.Errorf("CalVer.Bump(%s) release = %d, want %d",
						tt.part, gotCalVer.Release, tt.calver.Release+1)
				}
			}
		})
	}
}
