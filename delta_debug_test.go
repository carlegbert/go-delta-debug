package main

import (
	"bytes"
	"strings"
	"testing"
)

func TestRun(t *testing.T) {
	testFunc := func(lines []string) bool {
		for _, line := range lines {
			if line == "FAIL" {
				return true
			}
		}
		return false
	}

	input := strings.NewReader(strings.Join([]string{
		"PASS", "PASS", "PASS", "PASS", "PASS",
		"FAIL",
		"PASS", "PASS", "PASS", "PASS", "PASS", "PASS",
	}, "\n") + "\n")

	var output bytes.Buffer

	err := Run(testFunc, input, &output)
	if err != nil {
		t.Fatalf("Run returned error: %v", err)
	}

	outputStr := output.String()
	expected := "FAIL\n"
	if outputStr != expected {
		t.Errorf("Expected output %q, got %q", expected, outputStr)
	}
}
