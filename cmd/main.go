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
}

func run(testdir string, datadir string, filter string, verbose bool) error {
	fmt.Println()
	testData, err := loader.LoadTestData(testdir, datadir, filter)
	if err != nil {
		return err
	}
	fmt.Printf("Loaded %d test cases\n\n", len(testData.TestCases))

	runners := []model.Runner{
		&recolabs.RecolabsRunner{},
		&blues.BluesRunner{},
		&xiatechs.XiaTechsRunner{},
	}

	allResults := make(map[string]*model.RunnerResult, len(runners))

	fmt.Printf("Running tests ...  0%%")

	for i, r := range runners {
		res, err := r.RunTests(&testData)
		if err != nil {
			fmt.Printf("\rRunning tests ... ERROR\n")
			return fmt.Errorf("running tests with %s: %s", r.Name(), err)
		}
		allResults[r.Name()] = res
		pct := (i + 1) * 100 / len(runners)
		fmt.Printf("\rRunning tests ... %3d%%", pct)
	}
	fmt.Printf("\rRunning tests ... DONE   \n")
	fmt.Println("")

	if verbose {
		report.PromptDetailedReport(allResults)
	} else {
		report.PromptShortReport(allResults, verbose)
	}

	path := report.ReportPath()
	if err := report.WriteReport(path, allResults); err != nil {
		return fmt.Errorf("writing report: %w", err)
	}

	fmt.Printf("\nResults written to %s\n\n", path)

	return nil
}
