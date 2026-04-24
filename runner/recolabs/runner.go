package recolabs

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/recolabs/gnata"

	. "cycus.io/jsonata/eval/model"
)

type RecolabsRunner struct {
	rawDatasets     map[string]json.RawMessage
	decodedDatasets map[string]interface{}
	gnataFuncs      map[string]gnata.CustomFunc
	// evalCustom captures the gnata env in a closure to avoid naming its internal type.
	evalCustom func(compiled *gnata.Expression, data any) (any, error)
}

func (r *RecolabsRunner) RegisterCustomFunction(name string, fn CustomFunc) error {
	if r.gnataFuncs == nil {
		r.gnataFuncs = make(map[string]gnata.CustomFunc)
	}
	r.gnataFuncs[name] = func(args []any, _ any) (any, error) {
		iargs := make([]interface{}, len(args))
		copy(iargs, args)
		return fn(iargs)
	}
	env := gnata.NewCustomEnv(r.gnataFuncs)
	r.evalCustom = func(compiled *gnata.Expression, data any) (any, error) {
		return compiled.EvalWithCustomFuncs(context.Background(), data, env)
	}
	return nil
}

func (r *RecolabsRunner) Name() string {
	return "recolabs"
}

func (r *RecolabsRunner) getData(tc *TestCase) (interface{}, error) {
	if tc.Dataset != "" {
		if data, ok := r.decodedDatasets[tc.Dataset]; ok {
			return data, nil
		}
		raw, ok := r.rawDatasets[tc.Dataset]
		if !ok {
			return nil, fmt.Errorf("dataset %q not found", tc.Dataset)
		}
		data, err := gnata.DecodeJSON(raw)
		if err != nil {
			return nil, err
		}
		r.decodedDatasets[tc.Dataset] = data
		return data, nil
	}
	if len(tc.RawData) > 0 && string(tc.RawData) != "null" {
		return gnata.DecodeJSON(tc.RawData)
	}
	return nil, nil
}

func (r *RecolabsRunner) Eval(expr string, data interface{}, bindings map[string]interface{}) (interface{}, error) {
	compiled, err := gnata.Compile(expr)
	if err != nil {
		return nil, err
	}
	if len(bindings) > 0 {
		return compiled.EvalWithVars(context.Background(), data, bindings)
	}
	if r.evalCustom != nil {
		return r.evalCustom(compiled, data)
	}
	return compiled.Eval(context.Background(), data)
}

func (r *RecolabsRunner) RunTestCase(testCase *TestCase) *TestCaseResult {
	data, err := r.getData(testCase)
	if err != nil {
		return &TestCaseResult{TestCase: testCase, Error: err}
	}

	start := time.Now()
	result, evalErr := r.Eval(testCase.Expr, data, testCase.Bindings)
	duration := time.Since(start).Microseconds()

	tcr := &TestCaseResult{
		TestCase: testCase,
		Duration: duration,
	}

	if testCase.Error != "" || testCase.ErrorObject != nil {
		if evalErr != nil {
			if testCase.Error != "" {
				tcr.Passed = strings.Contains(evalErr.Error(), testCase.Error)
			} else {
				tcr.Passed = true
			}
		} else {
			tcr.Error = fmt.Errorf("expected error but expression succeeded with: %v", result)
		}
		return tcr
	}

	if evalErr != nil {
		tcr.Error = evalErr
		return tcr
	}

	if testCase.Undefined {
		tcr.Passed = result == nil
		return tcr
	}

	normalized := gnata.NormalizeValue(result)
	if testCase.Unordered {
		tcr.Passed = unorderedEqual(normalized, testCase.Result)
	} else {
		tcr.Passed = gnata.DeepEqual(normalized, testCase.Result)
	}
	return tcr
}

func (r *RecolabsRunner) RunTests(testData *TestData) (*RunnerResult, error) {
	r.rawDatasets = testData.Datasets
	r.decodedDatasets = make(map[string]interface{})

	runRes := &RunnerResult{
		Runner:   r,
		TestData: testData,
		Results:  make([]TestCaseResult, 0, len(testData.TestCases)),
	}
	for i := range testData.TestCases {
		result := r.RunTestCase(&testData.TestCases[i])
		runRes.Results = append(runRes.Results, *result)
	}
	return runRes, nil
}

func unorderedEqual(a, b any) bool {
	aSlice, aOk := a.([]any)
	bSlice, bOk := b.([]any)
	if !aOk || !bOk {
		return gnata.DeepEqual(a, b)
	}
	if len(aSlice) != len(bSlice) {
		return false
	}
	used := make([]bool, len(bSlice))
	for _, av := range aSlice {
		found := false
		for j, bv := range bSlice {
			if !used[j] && gnata.DeepEqual(av, bv) {
				used[j] = true
				found = true
				break
			}
		}
		if !found {
			return false
		}
	}
	return true
}
