package http_client

import (
	"context"
	"fmt"
	"net/http"
	"time"
)

func DoRequest(ctx context.Context, url string, timeout time.Duration) (*http.Response, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, fmt.Errorf("new request: %v", err)
	}

	cl := http.Client{
		Timeout: timeout,
	}

	return cl.Do(req)
}
