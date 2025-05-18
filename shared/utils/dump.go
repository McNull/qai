package utils

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"github.com/mattn/go-colorable"
	"github.com/neilotoole/jsoncolor"
)

// DumpInColor controls whether Dump and DumpString use colored output
var DumpInColor = true

// Dump method that accepts anything and dumps the content as a JSON output to the console with colors.
func Dump(v any) {
	str := DumpString(v)

	// Print the string to stdout
	if DumpInColor && jsoncolor.IsColorTerminal(os.Stdout) {
		out := colorable.NewColorable(os.Stdout) // Needed for Windows
		out.Write([]byte(str))
	} else {
		os.Stdout.Write([]byte(str))
	}

	fmt.Println() // Add a newline at the end
}

func DumpString(v any) string {
	var buf bytes.Buffer

	if DumpInColor {
		// Create encoder with color support
		enc := jsoncolor.NewEncoder(&buf)
		// Set colors similar to jq
		clrs := jsoncolor.DefaultColors()
		enc.SetColors(clrs)
		enc.SetIndent("", "  ") // Add indentation

		if err := enc.Encode(v); err != nil {
			return "Error: " + err.Error()
		}
	} else {
		// Use standard JSON encoder without colors
		enc := json.NewEncoder(&buf)
		enc.SetIndent("", "  ")

		if err := enc.Encode(v); err != nil {
			return "Error: " + err.Error()
		}
	}

	// Remove trailing newline that the encoder adds
	return strings.TrimSuffix(buf.String(), "\n")
}
