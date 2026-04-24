package model

import "encoding/json"

// CustomFunc is the common signature for user-defined functions that can be
// registered with any runner via RegisterCustomFunction.
type CustomFunc func(args []interface{}) (interface{}, error)

type Runner interface {
	RunTests(testData *TestData) (*RunnerResult, error)
	RunTestCase(testCase *TestCase) *TestCaseResult
	Eval(expr string, data interface{}, bindings map[string]interface{}) (interface{}, error)
	RegisterCustomFunction(name string, fn CustomFunc) error
	Name() string
}

type TestData struct {
	TestDir   string
	DataDir   string
	Filter    string
	TestCases []TestCase
	Datasets  map[string]json.RawMessage // raw bytes of each dataset file, keyed by name (without .json)
}

type TestCaseError struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

type TestCase struct {
	Expr        string
	ExprFile    string         `json:"expr-file"`
	Category    string
	RawData     json.RawMessage `json:"data"`
	Dataset     string
	Description string
	TimeLimit   int
	Depth       int
	Bindings    map[string]interface{}
	Result      interface{}
	Undefined   bool           `json:"undefinedResult"`
	Error       string         `json:"code"`
	ErrorObject *TestCaseError `json:"error"`
	Token       string
	Unordered   bool
	ID          string
}

type TestCaseResult struct {
	TestCase *TestCase
	Passed   bool
	Error    error
	Duration int64
}

type RunnerResult struct {
	Runner   Runner
	TestData *TestData
	Results  []TestCaseResult
}

func (r *RunnerResult) Passed() int {
	n := 0
	for _, res := range r.Results {
		if res.Passed {
			n++
		}
	}
	return n
}

func (r *RunnerResult) Failed() int {
	return len(r.Results) - r.Passed()
}

func (r *RunnerResult) TotalDuration() int64 {
	var total int64
	for _, res := range r.Results {
		total += res.Duration
	}
	return total
}
