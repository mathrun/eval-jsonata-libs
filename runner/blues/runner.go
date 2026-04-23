package blues

import (
	"encoding/json"
	"fmt"
	"reflect"
	"time"

	jsonata "github.com/blues/jsonata-go"

	. "cycus.io/jsonata/eval/model"
)

type BluesRunner struct {
	rawDatasets     map[string]json.RawMessage
	decodedDatasets map[string]interface{}
}

func (r *BluesRunner) Name() string {
	return "blues"
}

func (r *BluesRunner) getData(tc *TestCase) (interface{}, error) {
	if tc.Dataset != "" {
		if data, ok := r.decodedDatasets[tc.Dataset]; ok {
			return data, nil
		}
		raw, ok := r.rawDatasets[tc.Dataset]
		if !ok {
			return nil, fmt.Errorf("dataset %q not found", tc.Dataset)
		}
		var dest interface{}
		if err := json.Unmarshal(raw, &dest); err != nil {
			return nil, err
		}
		r.decodedDatasets[tc.Dataset] = dest
		return dest, nil
	}
	if len(tc.RawData) > 0 && string(tc.RawData) != "null" {
		var dest interface{}
		if err := json.Unmarshal(tc.RawData, &dest); err != nil {
			return nil, err
		}
		return dest, nil
	}
	return nil, nil
}

func (r *BluesRunner) Eval(expr string, data interface{}, bindings map[string]interface{}) (result interface{}, err error) {
	defer func() {
		if rec := recover(); rec != nil {
			err = fmt.Errorf("panic: %v", rec)
		}
	}()

	compiled, err := jsonata.Compile(expr)
	if err != nil {
		return nil, err
	}
	if len(bindings) > 0 {
		if err := compiled.RegisterVars(bindings); err != nil {
			return nil, err
		}
	}
	result, err = compiled.Eval(data)
	if err == jsonata.ErrUndefined {
		return nil, nil
	}
	return result, err
}

func (r *BluesRunner) RunTestCase(testCase *TestCase) *TestCaseResult {
	// blues has no tail-call optimisation; recursive expressions overflow the Go stack
	// which is a fatal crash that cannot be recovered, so skip that group entirely.
	if testCase.Category == "tail-recursion" {
		return &TestCaseResult{TestCase: testCase, Passed: false, Error: fmt.Errorf("skipped: tail-recursion not supported by blues")}
	}

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

	// Blues uses different error codes than the JSONata spec, so only check that an error occurred.
	if testCase.Error != "" || testCase.ErrorObject != nil {
		if evalErr != nil {
			tcr.Passed = true
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

	if testCase.Unordered {
		tcr.Passed = unorderedEqual(result, testCase.Result)
	} else {
		tcr.Passed = deepEqual(result, testCase.Result)
	}
	return tcr
}

func (r *BluesRunner) RunTests(testData *TestData) (*RunnerResult, error) {
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

// deepEqual normalises both values through JSON to ensure consistent types
// (e.g. int vs float64) before comparing.
func deepEqual(a, b interface{}) bool {
	if a == nil && b == nil {
		return true
	}
	aNorm, err := roundTripJSON(a)
	if err != nil {
		return reflect.DeepEqual(a, b)
	}
	bNorm, err := roundTripJSON(b)
	if err != nil {
		return reflect.DeepEqual(a, b)
	}
	return reflect.DeepEqual(aNorm, bNorm)
}

func roundTripJSON(v interface{}) (interface{}, error) {
	b, err := json.Marshal(v)
	if err != nil {
		return nil, err
	}
	var out interface{}
	if err := json.Unmarshal(b, &out); err != nil {
		return nil, err
	}
	return out, nil
}

func unorderedEqual(a, b interface{}) bool {
	aSlice, aOk := a.([]interface{})
	bSlice, bOk := b.([]interface{})
	if !aOk || !bOk {
		return deepEqual(a, b)
	}
	if len(aSlice) != len(bSlice) {
		return false
	}
	used := make([]bool, len(bSlice))
	for _, av := range aSlice {
		found := false
		for j, bv := range bSlice {
			if !used[j] && deepEqual(av, bv) {
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
