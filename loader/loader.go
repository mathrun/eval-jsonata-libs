package loader

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"cycus.io/jsonata/eval/model"
)

func LoadTestData(testdir string, datadir string, filter string) (model.TestData, error) {
	testData := model.TestData{
		TestDir:   testdir,
		DataDir:   datadir,
		Filter:    filter,
		TestCases: []model.TestCase{},
	}

	err := filepath.Walk(testdir, func(path string, info os.FileInfo, walkFnErr error) error {
		if info.IsDir() {
			if path == testdir {
				return nil
			}
			dirName := filepath.Base(path)
			if filter != "" && !strings.Contains(dirName, filter) {
				return filepath.SkipDir
			}
			return nil
		}

		if filepath.Ext(path) == ".jsonata" {
			return nil
		}

		testCases, err := loadTestCases(path)
		if err != nil {
			return fmt.Errorf("walk %s: %s", path, err)
		}

		category := filepath.Base(filepath.Dir(path))
		for i := range testCases {
			testCases[i].Category = category
		}

		testData.TestCases = append(testData.TestCases, testCases...)
		return nil
	})

	if err != nil {
		return model.TestData{}, err
	}

	testData.Datasets, err = loadDatasets(datadir)
	if err != nil {
		return model.TestData{}, err
	}

	return testData, nil
}

func loadDatasets(dataDir string) (map[string]json.RawMessage, error) {
	entries, err := os.ReadDir(dataDir)
	if err != nil {
		return nil, fmt.Errorf("reading datasets dir: %w", err)
	}
	datasets := make(map[string]json.RawMessage, len(entries))
	for _, entry := range entries {
		if entry.IsDir() || filepath.Ext(entry.Name()) != ".json" {
			continue
		}
		name := strings.TrimSuffix(entry.Name(), ".json")
		b, err := os.ReadFile(filepath.Join(dataDir, entry.Name()))
		if err != nil {
			return nil, fmt.Errorf("reading dataset %s: %w", entry.Name(), err)
		}
		datasets[name] = b
	}
	return datasets, nil
}

func loadTestCases(path string) ([]model.TestCase, error) {
	// Test cases are contained in json files. They consist of either
	// one test case in the file or an array of test cases.
	// Since we don't know which it will be until we load the file,
	// first try to demarshall it a single case, and if there is an
	// error, try again demarshalling it into an array of test cases
	var tc model.TestCase
	err := readJSONFile(path, &tc)
	if err != nil {
		var tcs []model.TestCase
		if err := readJSONFile(path, &tcs); err != nil {
			return nil, err
		}
		for i := range tcs {
			if err := normalizeTestCase(&tcs[i], path, i); err != nil {
				return nil, err
			}
		}
		return tcs, nil
	}

	if err := normalizeTestCase(&tc, path, 0); err != nil {
		return nil, err
	}
	return []model.TestCase{tc}, nil
}

func normalizeTestCase(tc *model.TestCase, path string, index int) error {
	if tc.ExprFile != "" {
		expr, err := loadExprFile(path, tc.ExprFile)
		if err != nil {
			return err
		}
		tc.Expr = expr
	}
	tc.ID = genTestCaseID(path, index)
	return nil
}

func loadExprFile(testPath string, exprFileName string) (string, error) {
	dir := filepath.Dir(testPath)
	content, err := os.ReadFile(filepath.Join(dir, exprFileName))
	if err != nil {
		return "", err
	}
	return string(content), nil
}

func genTestCaseID(path string, index int) string {
	base := strings.ReplaceAll(path, string(filepath.Separator), "-")
	return fmt.Sprintf("%s-%d", base, index)
}

func readJSONFile(path string, content interface{}) error {
	b, err := os.ReadFile(path)
	if err != nil {
		return fmt.Errorf("ReadFile %s: %s", path, err)
	}
	if err := json.Unmarshal(b, content); err != nil {
		return fmt.Errorf("unmarshal %s: %s", path, err)
	}
	return nil
}
