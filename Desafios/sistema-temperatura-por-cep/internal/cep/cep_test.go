package cep

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestViaCEP_City(t *testing.T) {
	tests := []struct {
		name     string
		cep      string
		status   int
		body     string
		wantCity string
		wantErr  error
	}{
		{
			name:    "invalid format - too short",
			cep:     "1234",
			wantErr: ErrInvalid,
		},
		{
			name:    "invalid format - letters",
			cep:     "abcdefgh",
			wantErr: ErrInvalid,
		},
		{
			name:     "success",
			cep:      "01001000",
			status:   http.StatusOK,
			body:     `{"localidade":"São Paulo"}`,
			wantCity: "São Paulo",
		},
		{
			name:    "not found - erro flag",
			cep:     "99999999",
			status:  http.StatusOK,
			body:    `{"erro":true}`,
			wantErr: ErrNotFound,
		},
		{
			name:    "not found - http status",
			cep:     "00000000",
			status:  http.StatusBadRequest,
			body:    ``,
			wantErr: ErrNotFound,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(tt.status)
				_, _ = w.Write([]byte(tt.body))
			}))
			defer server.Close()

			client := NewViaCEP(server.Client())
			client.BaseURL = server.URL

			city, err := client.City(context.Background(), tt.cep)
			if !errors.Is(err, tt.wantErr) {
				t.Fatalf("error = %v, want %v", err, tt.wantErr)
			}
			if city != tt.wantCity {
				t.Errorf("city = %q, want %q", city, tt.wantCity)
			}
		})
	}
}
