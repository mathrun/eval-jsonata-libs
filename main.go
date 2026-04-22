package main

import (
	"fmt"
	"log"
	"os"
	"time"

	. "cybus.io/jsonata/eval/model"
	blues "cybus.io/jsonata/eval/runner/blues"
	xiatechs "cybus.io/jsonata/eval/runner/xiatechs"

	"github.com/goccy/go-json"
)

func countAbsolutPassed(results []TestCollectionResult) int {
	var count int
	for _, collectionResult := range results {
		for _, result := range collectionResult.Results {
			if result.Passed {
				count++
			}
		}
	}
	return count
}

func countPassed(results TestCollectionResult) int {
	var count int
	for _, result := range results.Results {
		if result.Passed {
			count++
		}
	}
	return count
}

func calcAbsoluteDuration(results []TestCollectionResult) int64 {
	var total int64
	for _, collectionResult := range results {
		for _, result := range collectionResult.Results {
			total += result.Duration
		}
	}
	return total
}

func calcTotalDuration(results TestCollectionResult) int64 {
	var total int64
	for _, result := range results.Results {
		total += result.Duration
	}
	return total
}

func countTotalTests(results []TestCollectionResult) int {
	var count int
	for _, collectionResult := range results {
		count += len(collectionResult.Results)
	}
	return count
}

func printSummary(summary TestSummary) {
	fmt.Println()
	fmt.Println("==================================================")
	fmt.Printf("                  Test Summary\n")
	fmt.Println("==================================================")
	fmt.Println()

	for runnerName, collectionResults := range summary.Results {
		fmt.Printf("Runner: %s\n", runnerName)
		fmt.Printf("     Total Tests: %d\n", countTotalTests(collectionResults))
		fmt.Printf("    Total Passed: %d\n", countAbsolutPassed(collectionResults))
		fmt.Printf("  Total Duration: %d µs\n", calcAbsoluteDuration(collectionResults))
		fmt.Println()
	}
}

func printResults(runnerName string, results TestCollectionResult) {
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
		fmt.Printf("   Test: %s - %s (Duration: %d µs)\n", result.TestName, status, result.Duration)
	}
	fmt.Println()
}

func loadTestsFromFile(filePath string) (TestCollection, error) {
	var testCollection TestCollection

	fileBytes, err := os.ReadFile(filePath)
	if err != nil {
		return testCollection, err
	}

	err = json.Unmarshal(fileBytes, &testCollection)
	if err != nil {
		return testCollection, err
	}
	testCollection.FileName = filePath

	return testCollection, nil
}

func saveResultsToFile(results TestSummary) (string, error) {
	folder := "./results"
	if _, err := os.Stat(folder); os.IsNotExist(err) {
		err = os.Mkdir(folder, 0755)
		if err != nil {
			return "", err
		}
	}
	timestamp := time.Now().Format("20060102_150405")
	fileName := fmt.Sprintf("%s/%s.json", folder, timestamp)

	resultsBytes, err := json.MarshalIndent(results, "", "  ")
	if err != nil {
		return "", err
	}

	err = os.WriteFile(fileName, resultsBytes, 0644)
	if err != nil {
		return "", err
	}

	return fileName, nil
}

func main() {

	var runners = []TestRunner{
		&xiatechs.XiatechsRunner{},
		&blues.BluesRunner{},
	}

	testSummary := TestSummary{
		Results: make(map[string][]TestCollectionResult),
	}

	// for all files in folder ./testdata with suffix .json, load the tests and run them with all runners
	files, err := os.ReadDir("./testdata")
	if err != nil {
		log.Fatalf("Error reading testdata directory: %v\n", err)
	}

	for _, file := range files {

		testCollection, err := loadTestsFromFile(fmt.Sprintf("./testdata/%s", file.Name()))
		if err != nil {
			log.Printf("Error loading tests from file %s: %v\n", file.Name(), err)
			continue
		}

		for _, runner := range runners {
			fmt.Printf("Testing %s tests with %-25s ", file.Name(), runner.Name())

			results, err := runner.RunTests(testCollection)
			if err != nil {
				fmt.Println("ERROR")
				fmt.Printf("Error running tests with %s: %v\n", runner.Name(), err)
				continue
			}
			testSummary.Results[runner.Name()] = append(testSummary.Results[runner.Name()], results)
			fmt.Println("DONE")

		}
	}

	printSummary(testSummary)

	fileName, err := saveResultsToFile(testSummary)
	if err != nil {
		log.Printf("Error saving results to file: %v\n", err)
	} else {
		fmt.Printf("Results saved to file: %s\n", fileName)
	}

}
