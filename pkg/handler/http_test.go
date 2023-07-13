package handler

import (
	"indexer/pkg/mock"
	"indexer/pkg/store"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestHTTP_ServeHTTP(t *testing.T) {
	type fields struct {
		repo store.Repository
	}
	type args struct {
		r *http.Request
	}
	type result struct {
		code int
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		result result
	}{
		{
			name: "GET on '/' should be 200 OK",
			fields: fields{
				repo: mock.New(),
			},
			args: args{
				r: httptest.NewRequest(http.MethodGet, "/", nil),
			},
			result: result{
				code: http.StatusOK,
			},
		},
		{
			name: "GET on any other path except '/' should be 404",
			fields: fields{
				repo: mock.New(),
			},
			args: args{
				r: httptest.NewRequest(http.MethodGet, "/something_else", nil),
			},
			result: result{
				code: http.StatusNotFound,
			},
		},
		{
			name: "only GET is allowed, any other Method should be 405",
			fields: fields{
				repo: mock.New(),
			},
			args: args{
				r: httptest.NewRequest(http.MethodGet, "/something_else", nil),
			},
			result: result{
				code: http.StatusNotFound,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			h := &HTTP{
				repo: tt.fields.repo,
			}
			w := httptest.NewRecorder()
			h.ServeHTTP(w, tt.args.r)
			assert.Equal(t, tt.result.code, w.Result().StatusCode, "status code must match expected %d, got %d", tt.result.code, w.Result().StatusCode)
		})
	}
}
