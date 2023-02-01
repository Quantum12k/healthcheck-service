package healthcheck

import (
	"context"
	"io"
	"net/http"
	"strings"
	"testing"
)

func TestTextCheck_Execute(t *testing.T) {
	type args struct {
		response *http.Response
	}

	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "text_check_OK",
			args: args{
				response: &http.Response{
					Body: io.NopCloser(strings.NewReader("response text: ok")),
				},
			},
			wantErr: false,
		},
		{
			name: "text_check_ERROR",
			args: args{
				response: &http.Response{
					Body: io.NopCloser(strings.NewReader("response text: failed")),
				},
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			testCtx := context.Background()

			c := &TextCheck{}

			if err := c.Execute(testCtx, tt.args.response); (err != nil) != tt.wantErr {
				t.Errorf("Execute() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
