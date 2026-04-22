package xiatechs

import (
	"time"

	model "cybus.io/jsonata/eval/model"

	jsonata "github.com/blues/jsonata-go"
)

type BluesRunner struct{}

func (r *BluesRunner) RunTests(tests []model.Test) (model.TestSuiteResult, error) {
	var results []model.TestResult

	for _, test := range tests {
		timestamp_start := time.Now()
		e := jsonata.MustCompile(test.Expression)
		res, err := e.Eval(test.JSON)
		timestamp_elapsed := time.Since(timestamp_start).Microseconds()
		passed := err == nil && res == test.Expected
		results = append(results, model.TestResult{
			Test:     test,
			Passed:   passed,
			Duration: timestamp_elapsed,
		})
	}

	return model.TestSuiteResult{Results: results}, nil
}

func (r *BluesRunner) Name() string {
	return "Blues JSONata Runner"
}
