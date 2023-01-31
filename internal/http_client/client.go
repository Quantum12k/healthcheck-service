package http_client

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"time"
)

func DoRequest(ctx context.Context, method string, url string, body io.Reader, timeout time.Duration) (*http.Response, error) {
	req, err := http.NewRequestWithContext(ctx, method, url, body)
	if err != nil {
		return nil, fmt.Errorf("new request: %v", err)
	}

	cl := http.Client{
		Timeout: timeout,
	}

	return cl.Do(req)
}
