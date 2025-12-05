// Package output provides JSON formatting utilities.
package output

import (
	"encoding/json"
	"io"
	"os"
)

// JSONPrinter handles JSON output formatting.
type JSONPrinter struct {
	writer io.Writer
}

// NewJSONPrinter creates a new JSONPrinter.
func NewJSONPrinter() *JSONPrinter {
	return &JSONPrinter{
		writer: os.Stdout,
	}
}

// JSON outputs data as JSON.
func (jp *JSONPrinter) JSON(data interface{}) error {
	encoder := json.NewEncoder(jp.writer)
	encoder.SetIndent("", "  ")
	return encoder.Encode(data)
}