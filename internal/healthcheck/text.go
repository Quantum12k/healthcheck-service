package healthcheck

import (
	"context"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
)

const (
	CheckTextError = "response body doesn't contain 'ok'"
)

type (
	TextCheck struct {
	}
)

func (c *TextCheck) Execute(ctx context.Context, response *http.Response) error {
	respBody, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return fmt.Errorf("read body: %v", err)
	}

	if !strings.Contains(string(respBody), "ok") {
		return errors.New(CheckTextError)
	}

	return nil
}
