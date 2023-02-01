package healthcheck

import (
	"context"
	"net/http"
	"testing"
)

func TestStatusCodeCheck_Execute(t *testing.T) {
	type args struct {
		response *http.Response
	}

	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name:    "status_code_check_OK",
			args:    args{
				response: &http.Response{
					StatusCode:       200,
				},
			},
			wantErr: false,
		},
		{
			name:    "status_code_check_ERROR",
			args:    args{
				response: &http.Response{
					StatusCode:       500,
				},
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			testCtx := context.Background()

			c := &StatusCodeCheck{}

			if err := c.Execute(testCtx, tt.args.response); (err != nil) != tt.wantErr {
				t.Errorf("Execute() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}