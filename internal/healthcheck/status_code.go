package healthcheck

import (
	"context"
	"errors"
	"net/http"
)

const (
	CheckStatusCodeError = "status code != 200"
)

type (
	StatusCodeCheck struct {
	}
)

func (c *StatusCodeCheck) Execute(ctx context.Context, response *http.Response) error {
	if response.StatusCode != http.StatusOK {
		return errors.New(CheckStatusCodeError)
	}

	return nil
}
