package output

import (
	"bytes"
	"encoding/json"
	"testing"
)

func TestPrinter_Println(t *testing.T) {
	tests := []struct {
		name     string
		jsonMode bool
		format   string
		args     []interface{}
		want     string
	}{
		{
			name:     "human mode prints message",
			jsonMode: false,
			format:   "Version: %s",
			args:     []interface{}{"1.0.0"},
			want:     "Version: 1.0.0\n",
		},
		{
			name:     "json mode suppresses output",
			jsonMode: true,
			format:   "Version: %s",
			args:     []interface{}{"1.0.0"},
			want:     "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var buf bytes.Buffer
			p := &Printer{
				writer:   &buf,
				jsonMode: tt.jsonMode,
			}
			p.Println(tt.format, tt.args...)
			if got := buf.String(); got != tt.want {
				t.Errorf("Println() output = %q, want %q", got, tt.want)
			}
		})
	}
}

func TestPrinter_Error(t *testing.T) {
	var buf bytes.Buffer
	p := &Printer{
		errWriter: &buf,
	}
	p.Error("something went wrong: %s", "details")
	want := "Error: something went wrong: details\n"
	if got := buf.String(); got != want {
		t.Errorf("Error() output = %q, want %q", got, want)
	}
}

func TestPrinter_Warning(t *testing.T) {
	var buf bytes.Buffer
	p := &Printer{
		errWriter: &buf,
	}
	p.Warning("this is a warning: %s", "info")
	want := "Warning: this is a warning: info\n"
	if got := buf.String(); got != want {
		t.Errorf("Warning() output = %q, want %q", got, want)
	}
}

func TestBumpResult_JSON(t *testing.T) {
	result := BumpResult{
		OldVersion:   "1.0.0",
		NewVersion:   "1.0.1",
		Part:         "patch",
		ChangedFiles: []string{"package.json", "VERSION"},
		GitCommit:    "abc1234",
		GitTag:       "v1.0.1",
	}

	data, err := json.Marshal(result)
	if err != nil {
		t.Fatalf("Failed to marshal BumpResult: %v", err)
	}

	var unmarshaled map[string]interface{}
	if err := json.Unmarshal(data, &unmarshaled); err != nil {
		t.Fatalf("Failed to unmarshal: %v", err)
	}

	// Check field names match spec
	if _, ok := unmarshaled["old_version"]; !ok {
		t.Error("Expected 'old_version' field in JSON")
	}
	if _, ok := unmarshaled["new_version"]; !ok {
		t.Error("Expected 'new_version' field in JSON")
	}
	if _, ok := unmarshaled["files_changed"]; !ok {
		t.Error("Expected 'files_changed' field in JSON")
	}
	if _, ok := unmarshaled["git_committed"]; !ok {
		t.Error("Expected 'git_committed' field in JSON")
	}
	if _, ok := unmarshaled["git_tagged"]; !ok {
		t.Error("Expected 'git_tagged' field in JSON")
	}
}

func TestShowResult_JSON(t *testing.T) {
	result := ShowResult{
		Version: "2.0.0",
	}

	data, err := json.Marshal(result)
	if err != nil {
		t.Fatalf("Failed to marshal ShowResult: %v", err)
	}

	var unmarshaled map[string]interface{}
	if err := json.Unmarshal(data, &unmarshaled); err != nil {
		t.Fatalf("Failed to unmarshal: %v", err)
	}

	if v, ok := unmarshaled["version"]; !ok || v != "2.0.0" {
		t.Errorf("Expected version=2.0.0, got %v", v)
	}
}

func TestValidationResult_JSON(t *testing.T) {
	result := ValidationResult{
		Valid: false,
		Errors: []ValidationError{
			{Field: "project.current_version", Message: "required field"},
		},
		Warnings: []string{"no files configured"},
	}

	data, err := json.Marshal(result)
	if err != nil {
		t.Fatalf("Failed to marshal ValidationResult: %v", err)
	}

	var unmarshaled map[string]interface{}
	if err := json.Unmarshal(data, &unmarshaled); err != nil {
		t.Fatalf("Failed to unmarshal: %v", err)
	}

	if v, ok := unmarshaled["valid"]; !ok || v != false {
		t.Errorf("Expected valid=false, got %v", v)
	}
}
