// Package output provides output formatting utilities.
package output

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
)

// Printer handles output formatting.
type Printer struct {
	writer    io.Writer
	errWriter io.Writer
	jsonMode  bool
}

// NewPrinter creates a new Printer.
func NewPrinter(jsonMode bool) *Printer {
	return &Printer{
		writer:    os.Stdout,
		errWriter: os.Stderr,
		jsonMode:  jsonMode,
	}
}

// Print prints a message to stdout.
func (p *Printer) Print(format string, args ...interface{}) {
	if p.jsonMode {
		return
	}
	fmt.Fprintf(p.writer, format, args...)
}

// Println prints a message to stdout with a newline.
func (p *Printer) Println(format string, args ...interface{}) {
	if p.jsonMode {
		return
	}
	fmt.Fprintf(p.writer, format+"\n", args...)
}

// Error prints an error message to stderr.
func (p *Printer) Error(format string, args ...interface{}) {
	fmt.Fprintf(p.errWriter, "Error: "+format+"\n", args...)
}

// Warning prints a warning message to stderr.
func (p *Printer) Warning(format string, args ...interface{}) {
	fmt.Fprintf(p.errWriter, "Warning: "+format+"\n", args...)
}

// JSON outputs data as JSON.
func (p *Printer) JSON(data interface{}) error {
	encoder := json.NewEncoder(p.writer)
	encoder.SetIndent("", "  ")
	return encoder.Encode(data)
}

// ShowResult represents the result of a show command.
type ShowResult struct {
	Version string `json:"version"`
}

// BumpResult represents the result of a bump command.
type BumpResult struct {
	OldVersion   string   `json:"old_version"`
	NewVersion   string   `json:"new_version"`
	Part         string   `json:"part"`
	ChangedFiles []string `json:"changed_files,omitempty"`
	GitCommit    string   `json:"git_commit,omitempty"`
	GitTag       string   `json:"git_tag,omitempty"`
}

// ValidationResult represents the result of a validate command.
type ValidationResult struct {
	Valid    bool              `json:"valid"`
	Errors   []ValidationError `json:"errors,omitempty"`
	Warnings []string          `json:"warnings,omitempty"`
}

// ValidationError represents a validation error.
type ValidationError struct {
	Field   string `json:"field"`
	Message string `json:"message"`
}
