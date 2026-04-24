package report

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"cycus.io/jsonata/eval/model"
)

const resultsDir = "results"

func ReportPath() string {
	ts := time.Now().Format("20060102_150405")
	return filepath.Join(resultsDir, fmt.Sprintf("%s_results.json", ts))
}

type testCaseResultJSON struct {
	ID       string `json:"id"`
	Passed   bool   `json:"passed"`
	Duration int64  `json:"duration_µs"`
	Error    string `json:"error,omitempty"`
}

type runnerSummaryJSON struct {
	Total    int                  `json:"total"`
	Passed   int                  `json:"passed"`
	Duration int64                `json:"duration_µs"`
	Results  []testCaseResultJSON `json:"results"`
}

func PromptShortReport(results map[string]*model.RunnerResult, verbose bool) {
	fmt.Printf("Summary:\n\n")
	for name, res := range results {
		totalDuration := res.TotalDuration()
		if totalDuration > 1000 {
			fmt.Printf("%10s: %4d passed, %4d failed, total duration: %d ms\n",
				name, res.Passed(), res.Failed(), totalDuration/1000)
		} else {
			fmt.Printf("%10s: %4d passed, %4d failed, total duration: %d µs\n",
				name, res.Passed(), res.Failed(), totalDuration)
		}
	}
}

func PromptDetailedReport(results map[string]*model.RunnerResult) {
	fmt.Printf("\nDetailed Results:\n\n")
	for name, res := range results {

		totalDuration := res.TotalDuration()
		if totalDuration > 1000 {
			fmt.Printf("\n%10s: %4d passed, %4d failed, total duration: %d ms\n",
				name, res.Passed(), res.Failed(), totalDuration/1000)
		} else {
			fmt.Printf("\n%10s: %4d passed, %4d failed, total duration: %d µs\n",
				name, res.Passed(), res.Failed(), totalDuration)
		}

		for _, r := range res.Results {
			if !r.Passed {
				fmt.Printf("    - %s: FAILED : %v)\n", r.TestCase.ID, r.Error)
			}
		}
	}
}

func WriteReport(path string, results map[string]*model.RunnerResult) error {
	out := make(map[string]runnerSummaryJSON, len(results))

	for name, res := range results {
		entries := make([]testCaseResultJSON, len(res.Results))
		for i, r := range res.Results {
			entry := testCaseResultJSON{
				ID:       r.TestCase.ID,
				Passed:   r.Passed,
				Duration: r.Duration,
			}
			if r.Error != nil {
				entry.Error = r.Error.Error()
			}
			entries[i] = entry
		}
		out[name] = runnerSummaryJSON{
			Total:    len(res.Results),
			Passed:   res.Passed(),
			Duration: res.TotalDuration(),
			Results:  entries,
		}
	}

	b, err := json.MarshalIndent(out, "", "  ")
	if err != nil {
		return err
	}
	if err := os.MkdirAll(filepath.Dir(path), 0755); err != nil {
		return err
	}
	return os.WriteFile(path, b, 0644)
}
