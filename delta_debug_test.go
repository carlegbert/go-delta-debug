package deltadebug_test

import (
	"bytes"
	"strings"
	"testing"

	deltadebug "github.com/carlegbert/go-delta-debug"
	"slices"
)

func TestRun(t *testing.T) {
	testFunc := func(lines []string) bool {
		return slices.Contains(lines, "FAIL")
	}

	input := strings.NewReader(strings.Join([]string{
		"PASS", "PASS", "PASS", "PASS", "PASS",
		"FAIL",
		"PASS", "PASS", "PASS", "PASS", "PASS", "PASS",
	}, "\n") + "\n")

	var output bytes.Buffer

	err := deltadebug.Run(testFunc, input, &output)
	if err != nil {
		t.Fatalf("Run returned error: %v", err)
	}

	outputStr := output.String()
	expected := "FAIL\n"
	if outputStr != expected {
		t.Errorf("Expected output %q, got %q", expected, outputStr)
	}
}
