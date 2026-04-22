package model

type Test struct {
	Name        string      `json:"name"`
	Description string      `json:"description"`
	JSON        interface{} `json:"json"`
	Expression  string      `json:"expression"`
	Expected    interface{} `json:"expected"`
}

type TestCollection struct {
	Tests []Test `json:"tests"`
}

type TestResult struct {
	Test     Test  `json:"test"`
	Duration int64 `json:"duration"`
	Passed   bool  `json:"passed"`
	Error    error `json:"error,omitempty"`
}

type TestSuiteResult struct {
	Results []TestResult `json:"results"`
}

type TestRunner interface {
	RunTests(tests []Test) (TestSuiteResult, error)
	Name() string
}
