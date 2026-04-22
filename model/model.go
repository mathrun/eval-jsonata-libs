package model

type Test struct {
	Name        string      `json:"name"`
	Description string      `json:"description"`
	JSON        interface{} `json:"json"`
	Expression  string      `json:"expression"`
	Expected    interface{} `json:"expected"`
}

type TestCollection struct {
	Tests    []Test `json:"tests"`
	FileName string `json:"-"`
}

type TestResult struct {
	TestName string `json:"name"`
	Duration int64  `json:"duration"`
	Passed   bool   `json:"passed"`
	Error    string `json:"error,omitempty"`
}

type TestCollectionResult struct {
	File    string       `json:"file"`
	Results []TestResult `json:"results"`
}

type TestSummary struct {
	Results map[string][]TestCollectionResult `json:"results"`
}

type TestRunner interface {
	RunTests(collection TestCollection) (TestCollectionResult, error)
	Name() string
}
