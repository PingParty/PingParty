package main

import (
	"net/http"
	"sync"

	"github.com/aws/aws-lambda-go/lambda"
)

type Monitor struct {
	URL string
}

type Result struct {
	StatusCode int
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
			resp, err := http.Get(c.URL)
			if err != nil {
				result.StatusCode = -1
			} else {
				result.StatusCode = resp.StatusCode
			}
			results[i] = result
		}(i, c)
	}
	wg.Wait()
	return results, nil
}
