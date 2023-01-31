package healthcheck

import (
	"context"
	"fmt"
	"net/http"
	"sync"
	"time"

	"github.com/Quantum12k/healthcheck-service/internal/http_client"
)

const (
	MaxActiveHealthChecks = 5

	StatusCode = "status_code"
	Text       = "text"
)

type (
	HealthChecker interface {
		Execute(context.Context, *http.Response) error
	}

	URL struct {
		URL            string   `yaml:"url"`
		Checks         []string `yaml:"checks"`
		MinChecksCount int      `yaml:"min_checks_cnt"`
	}

	CheckResult struct {
		URL    string `db:"url"`
		Result string `db:"result"`
	}

	Config struct {
		TimeoutMS int `yaml:"timeout_ms"`
	}
)

func (r *CheckResult) String() string {
	return fmt.Sprintf("%s %s", r.URL, r.Result)
}

func HandleURLs(ctx context.Context, urls []URL) []CheckResult {
	results := make([]CheckResult, 0, len(urls))
	resultCh := make(chan CheckResult)
	resultWg := sync.WaitGroup{}

	go func() {
		resultWg.Add(1)
		defer resultWg.Done()

		for res := range resultCh {
			results = append(results, res)
		}
	}()

	checksWg := sync.WaitGroup{}
	limiter := make(chan struct{}, MaxActiveHealthChecks)

urlLoop:
	for _, url := range urls {
		select {
		case <-ctx.Done():
			break urlLoop
		default:
			limiter <- struct{}{}
			checksWg.Add(1)

			go func(url URL) {
				defer func() {
					<-limiter
					checksWg.Done()
				}()

				res := CheckResult{
					URL:    url.URL,
					Result: "ok",
				}

				if err := handleURL(ctx, url); err != nil {
					res.Result = err.Error()
				}

				resultCh <- res
			}(url)
		}
	}

	checksWg.Wait()
	close(resultCh)

	resultWg.Wait()

	return results
}

func handleURL(ctx context.Context, url URL) error {
	resp, err := http_client.DoRequest(ctx, http.MethodGet, url.URL, nil, time.Second)
	if err != nil {
		return fmt.Errorf("unable to do request: %v", err)
	}

	errorsAmount := 0
	failText := "fail"

	for _, checkType := range url.Checks {
		var check HealthChecker

		switch checkType {
		case StatusCode:
			check = &StatusCodeCheck{}
		case Text:
			check = &TextCheck{}
		default:
			failText += checkType
			errorsAmount++
			continue
		}

		if err := check.Execute(ctx, resp); err != nil {
			failText += checkType
			errorsAmount++
			continue
		}
	}

	if errorsAmount > (len(url.Checks) - url.MinChecksCount) {
		return fmt.Errorf(failText)
	}

	return nil
}
