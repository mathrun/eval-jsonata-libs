package main

import (
	"fmt"
	"log"
	"os"

	"cybus.io/jsonata/eval/model"
	blues "cybus.io/jsonata/eval/runner/blues"
	xiatechs "cybus.io/jsonata/eval/runner/xiatechs"

	"github.com/goccy/go-json"
)

func countPassed(results model.TestSuiteResult) int {
	var count int
	for _, result := range results.Results {
		if result.Passed {
			count++
		}
	}
	return count
}

func calcTotalDuration(results model.TestSuiteResult) int64 {
	var total int64
	for _, result := range results.Results {
		total += result.Duration
	}
	return total
}

func printResults(runnerName string, results model.TestSuiteResult) {
	fmt.Println()
	fmt.Println("==================================================")
	fmt.Printf("      %s\n", runnerName)
	fmt.Println("--------------------------------------------------")
	fmt.Println()
	fmt.Printf("     Total Tests: %d\n", len(results.Results))
	fmt.Printf("    Total Passed: %d\n", countPassed(results))
	fmt.Printf("  Total Duration: %d µs\n", calcTotalDuration(results))
	fmt.Println()
	fmt.Println("  Detailed Results:")
	fmt.Println()
	for _, result := range results.Results {
		status := "FAILED"
		if result.Passed {
			status = "PASSED"
		}
		fmt.Printf("   Test: %s - %s (Duration: %d µs)\n", result.Test.Name, status, result.Duration)
	}
	fmt.Println()
}

func loadTestsFromFile(filePath string) (model.TestCollection, error) {
	var testCollection model.TestCollection

	fileBytes, err := os.ReadFile(filePath)
	if err != nil {
		return testCollection, err
	}

	err = json.Unmarshal(fileBytes, &testCollection)
	if err != nil {
		return testCollection, err
	}

	return testCollection, nil
}

func main() {

	var runners = []model.TestRunner{
		&xiatechs.XiatechsRunner{},
		&blues.BluesRunner{},
	}

	var testResults = make(map[string]model.TestSuiteResult)

	fmt.Println()
	if len(os.Args) < 2 {
		log.Fatal("Usage: program <json-file>")
	}

	input, err := loadTestsFromFile(os.Args[1])
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Tests loaded, ready to evaluate.")
	fmt.Println()

	for _, runner := range runners {
		fmt.Printf("Running tests with %-25s ", runner.Name())

		results, err := runner.RunTests(input.Tests)
		if err != nil {
			fmt.Println("ERROR")
			fmt.Printf("Error running tests with %s: %v\n", runner.Name(), err)
			continue
		}
		fmt.Println("DONE")

		testResults[runner.Name()] = results
	}

	for runnerName, results := range testResults {
		printResults(runnerName, results)
	}
}
