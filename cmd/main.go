package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"

	"cycus.io/jsonata/eval/loader"
	"cycus.io/jsonata/eval/model"
	"cycus.io/jsonata/eval/report"
	"cycus.io/jsonata/eval/runner/blues"
	"cycus.io/jsonata/eval/runner/recolabs"
	"cycus.io/jsonata/eval/runner/xiatechs"
)

func main() {
	var group string
	var verbose bool

	flag.BoolVar(&verbose, "verbose", false, "verbose output")
	flag.StringVar(&group, "group", "", "restrict to one or more test groups")
	flag.Parse()

	if flag.NArg() != 1 {
		fmt.Fprintln(os.Stderr, "Syntax: jsonata-test [options] <directory>")
		os.Exit(1)
	}

	root := flag.Arg(0)
	testdir := filepath.Join(root, "groups")
	datadir := filepath.Join(root, "datasets")

	err := run(testdir, datadir, group, verbose)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error while running: %s\n", err)
		os.Exit(2)
	}

	fmt.Fprintln(os.Stdout, "OK")
}

func run(testdir string, datadir string, filter string, verbose bool) error {
	testData, err := loader.LoadTestData(testdir, datadir, filter)
	if err != nil {
		return err
	}
	fmt.Printf("Loaded %d test cases\n", len(testData.TestCases))

	runners := []model.Runner{
		&recolabs.RecolabsRunner{},
		&blues.BluesRunner{},
		&xiatechs.XiaTechsRunner{},
	}

	allResults := make(map[string]*model.RunnerResult, len(runners))

	for _, r := range runners {
		fmt.Printf("Running tests with %s runner\n", r.Name())
		res, err := r.RunTests(&testData)
		if err != nil {
			return fmt.Errorf("running tests with %s: %s", r.Name(), err)
		}
		allResults[r.Name()] = res

		if verbose {
			for _, result := range res.Results {
				if !result.Passed {
					if result.Error != nil {
						fmt.Printf("  FAIL %s: %s\n", result.TestCase.ID, result.Error)
					} else {
						fmt.Printf("  FAIL %s\n", result.TestCase.ID)
					}
				}
			}
		}
		fmt.Printf("  %s: %d passed, %d failed, total duration: %dms\n",
			r.Name(), res.Passed(), res.Failed(), res.TotalDuration()/1000)
	}

	path := report.ReportPath()
	if err := report.WriteReport(path, allResults); err != nil {
		return fmt.Errorf("writing report: %w", err)
	}
	fmt.Printf("Results written to %s\n", path)

	return nil
}
