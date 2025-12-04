package main

import (
	"bufio"
	"fmt"
	"math"
	"os"
)

type TestFunc func([]string) bool

type DeltaDebugger struct {
	test             TestFunc
	initialInputFile string
	outputFile       string
}

func Run(
	test TestFunc,
	initialInputFile string,
	outputFile string,
) error {
	dd := &DeltaDebugger{
		test:             test,
		initialInputFile: initialInputFile,
		outputFile:       outputFile,
	}

	file, err := os.Open(dd.initialInputFile)
	if err != nil {
		return fmt.Errorf("error opening initial input file: %w", err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	lines := []string{}
	lines, err = collectLines(scanner, lines)
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
			err = writeFailedTestCase(dd, testCase)
			if err != nil {
				return err
			}

			continue
		}

		nablas := getNablas(lines, n)
		testCase, failed = getFirstFailingTestcase(dd, nablas)
		if failed {
			n -= 1
			lines = testCase
			err = writeFailedTestCase(dd, testCase)
			if err != nil {
				return err
			}

			continue
		}

		n *= 2
	}

	return nil
}

func collectLines(scanner *bufio.Scanner, lines []string) ([]string, error) {
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}
	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("error reading initial input file: %w", err)
	}
	return lines, nil
}

func writeFailedTestCase(dd *DeltaDebugger, testCase []string) error {
	outputFile, err := os.Create(dd.outputFile)
	if err != nil {
		return fmt.Errorf("error creating output file: %w", err)
	}
	defer outputFile.Close()
	for _, line := range testCase {
		outputFile.WriteString(line + "\n")
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
