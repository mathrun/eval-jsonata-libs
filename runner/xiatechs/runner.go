package xiatechs

import (
	"time"

	model "cybus.io/jsonata/eval/model"

	"github.com/xiatechs/jsonata-go"
)

type XiatechsRunner struct{}

func (r *XiatechsRunner) RunTests(collection model.TestCollection) (model.TestCollectionResult, error) {
	var results []model.TestResult

	for _, test := range collection.Tests {
		timestamp_start := time.Now()
		e := jsonata.MustCompile(test.Expression)
		res, err := e.Eval(test.JSON)
		timestamp_elapsed := time.Since(timestamp_start).Microseconds()
		passed := err == nil && res == test.Expected
		results = append(results, model.TestResult{
			TestName: test.Name,
			Passed:   passed,
			Duration: timestamp_elapsed,
		})
	}

	return model.TestCollectionResult{File: collection.FileName, Results: results}, nil
}

func (r *XiatechsRunner) Name() string {
	return "Xiatechs JSONata Runner"
}
