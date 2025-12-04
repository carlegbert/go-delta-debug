package main

import (
	"bufio"
	"fmt"
	"io"
	"math"
)

type TestFunc func([]string) bool

type DeltaDebugger struct {
	test         TestFunc
	inputReader  io.Reader
	outputWriter io.Writer
}

func Run(
	test TestFunc,
	inputReader io.Reader,
	outputWriter io.Writer,
) error {
	dd := &DeltaDebugger{
		test:         test,
		inputReader:  inputReader,
		outputWriter: outputWriter,
	}

	scanner := bufio.NewScanner(dd.inputReader)
	lines := []string{}
	lines, err := collectLines(scanner, lines)
	if err != nil {
		return err
	}

	n := 2
	for n <= len(lines) {
		testCases := getCases(lines, n)
		testCase, failed := getFirstFailingTestcase(dd, testCases)
		if failed {
			n = 2
			lines = testCase
			continue
		}

		nablas := getNablas(lines, n)
		testCase, failed = getFirstFailingTestcase(dd, nablas)
		if failed {
			n -= 1
			lines = testCase
			continue
		}

		n *= 2
	}

	err = writeFailedTestCase(dd, lines)
	if err != nil {
		return err
	}

	return nil
}

func collectLines(scanner *bufio.Scanner, lines []string) ([]string, error) {
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}
	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("error reading input: %w", err)
	}
	return lines, nil
}

func writeFailedTestCase(dd *DeltaDebugger, testCase []string) error {
	for _, line := range testCase {
		_, err := fmt.Fprintln(dd.outputWriter, line)
		if err != nil {
			return fmt.Errorf("error writing output: %w", err)
		}
	}
	return nil
}

func getCases(lines []string, n int) chan []string {
	cases := make(chan []string)
	chunkSize := len(lines) / n
	go func() {
		defer close(cases)
		for start := 0; start < len(lines); start += chunkSize {
			end := int(math.Min(float64(start)+float64(chunkSize), float64(len(lines))))
			cases <- lines[start:end]
		}
	}()

	return cases
}

func getNablas(lines []string, n int) chan []string {
	nablas := make(chan []string)
	chunkSize := len(lines) / n
	go func() {
		defer close(nablas)
		for start := 0; start < len(lines); start += chunkSize {
			end := int(math.Min(float64(start)+float64(chunkSize), float64(len(lines))))
			nablas <- append(lines[:start], lines[end:]...)
		}
	}()

	return nablas
}

func getFirstFailingTestcase(dd *DeltaDebugger, cases chan []string) ([]string, bool) {
	for testCase := range cases {
		if failed := dd.test(testCase); failed {
			return testCase, true
		}
	}

	return nil, false
}
