package response

import (
	"errors"
	"net/http"

	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/render"
)

type (
	ErrResponse struct {
		Err            error `json:"-"`
		HTTPStatusCode int   `json:"-"`

		StatusText string      `json:"errorTxt"`
		ErrorText  string      `json:"error,omitempty"`
		RequestID  interface{} `json:"request_id,omitempty"`
	}

	Response struct {
		Data           interface{} `json:"data,omitempty"`
		RequestID      interface{} `json:"request_id,omitempty"`
		HTTPStatusCode int         `json:"-"`
	}
)

func (e *ErrResponse) Render(w http.ResponseWriter, r *http.Request) error {
	render.Status(r, e.HTTPStatusCode)
	return nil
}

func (mr *Response) Render(w http.ResponseWriter, r *http.Request) error {
	render.Status(r, mr.HTTPStatusCode)
	return nil
}

func RenderResp(data interface{}, statusCode int) render.Renderer {
	return &Response{
		Data:           data,
		HTTPStatusCode: statusCode,
		RequestID:      middleware.NextRequestID(),
	}
}

func ErrInternal(errInt interface{}) render.Renderer {
	err := getError(errInt)

	return &ErrResponse{
		Err:            err,
		HTTPStatusCode: http.StatusInternalServerError,
		StatusText:     "Internal Server Error",
		ErrorText:      err.Error(),
		RequestID:      middleware.NextRequestID(),
	}
}

func getError(errInt interface{}) error {
	var errStr string

	err, ok := errInt.(error)
	if !ok {
		errStr, ok = errInt.(string)
		if ok {
			err = errors.New(errStr)
		} else {
			err = errors.New("undefined error type")
		}
	}

	return err
}
