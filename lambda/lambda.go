package main

import (
	"net/http"
	"sync"
	"time"

	"github.com/aws/aws-lambda-go/lambda"
)

type Monitor struct {
	URL string
}

type Result struct {
	StatusCode int
	Duration   float64
	Message    string
}

func main() {
	lambda.Start(check)
}

func check(checks []*Monitor) ([]*Result, error) {
	results := make([]*Result, len(checks))
	wg := sync.WaitGroup{}
	for i, c := range checks {
		wg.Add(1)
		go func(i int, c *Monitor) {
			defer wg.Done()
			result := &Result{}
			start := time.Now()
			resp, err := http.Get(c.URL)
			if err != nil {
				result.StatusCode = -1
				result.Message = err.Error()
			} else {
				result.StatusCode = resp.StatusCode
			}
			result.Duration = float64(time.Now().Sub(start)) / float64(time.Millisecond)
			results[i] = result
		}(i, c)
	}
	wg.Wait()
	return results, nil
}
