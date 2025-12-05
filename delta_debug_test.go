package deltadebug_test

import (
	"bytes"
	"strings"
	"testing"

	deltadebug "github.com/carlegbert/go-delta-debug"
	"slices"
)

func TestRun(t *testing.T) {
	containsFail := func(lines []string) bool {
		return slices.Contains(lines, "FAIL")
	}
	containsTwoItems := func(lines []string) bool {
		return slices.Contains(lines, "a") && slices.Contains(lines, "b")
	}

	testCases := []struct {
		name     string
		testFunc deltadebug.TestFunc
		lines    []string
		expected []string
	}{
		{
			name:     "Single failing line",
			testFunc: containsFail,
			lines:    []string{"FAIL"},
			expected: []string{"FAIL"},
		},
		{
			name:     "Two failing lines",
			testFunc: containsFail,
			lines:    []string{"FAIL", "FAIL"},
			expected: []string{"FAIL"},
		},
		{
			name:     "Multiple lines with one failure",
			testFunc: containsFail,
			lines: []string{
				"PASS", "PASS", "PASS", "PASS", "PASS", "FAIL",
				"PASS", "PASS", "PASS", "PASS", "PASS", "PASS",
			},
			expected: []string{"FAIL"},
		},
		{
			name:     "Multi-line failure",
			testFunc: containsFail,
			lines:    []string{"a", "b"},
			expected: []string{"a", "b"},
		},
		{
			name:     "Multi-line failure with other lines",
			testFunc: containsTwoItems,
			lines:    []string{"x", "x", "x", "a", "x", "b"},
			expected: []string{"a", "b"},
		},
	}

	for _, tt := range testCases {
		t.Run(tt.name, func(t *testing.T) {
			input := strings.NewReader(strings.Join(tt.lines, "\n") + "\n")
			var output bytes.Buffer
			err := deltadebug.Run(tt.testFunc, input, &output)
			if err != nil {
				t.Fatalf("Run returned error: %v", err)
			}
			outputStr := output.String()
			expected := strings.Join(tt.expected, "\n") + "\n"
			if outputStr != expected {
				t.Errorf("Expected output %q, got %q", expected, outputStr)
			}
		})
	}

}
