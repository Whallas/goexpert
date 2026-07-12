package handler

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/Whallas/goexpert/Desafios/sistema-temperatura-por-cep/internal/cep"
	"github.com/Whallas/goexpert/Desafios/sistema-temperatura-por-cep/internal/temperature"
)

type stubFinder struct {
	city string
	err  error
}

func (s stubFinder) City(context.Context, string) (string, error) { return s.city, s.err }

type stubProvider struct {
	celsius float64
	err     error
}

func (s stubProvider) CelsiusFor(context.Context, string) (float64, error) {
	return s.celsius, s.err
}

func TestWeather_ServeHTTP(t *testing.T) {
	tests := []struct {
		name       string
		finder     stubFinder
		provider   stubProvider
		wantStatus int
		wantBody   string
	}{
		{
			name:       "success",
			finder:     stubFinder{city: "São Paulo"},
			provider:   stubProvider{celsius: 28.5},
			wantStatus: http.StatusOK,
		},
		{
			name:       "invalid zipcode",
			finder:     stubFinder{err: cep.ErrInvalid},
			wantStatus: http.StatusUnprocessableEntity,
			wantBody:   "invalid zipcode",
		},
		{
			name:       "zipcode not found",
			finder:     stubFinder{err: cep.ErrNotFound},
			wantStatus: http.StatusNotFound,
			wantBody:   "can not find zipcode",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mux := http.NewServeMux()
			mux.Handle("GET /weather/{cep}", NewWeather(tt.finder, tt.provider))

			req := httptest.NewRequest(http.MethodGet, "/weather/01001000", nil)
			rec := httptest.NewRecorder()
			mux.ServeHTTP(rec, req)

			if rec.Code != tt.wantStatus {
				t.Fatalf("status = %d, want %d", rec.Code, tt.wantStatus)
			}

			if tt.wantStatus == http.StatusOK {
				var got temperature.Temperatures
				if err := json.Unmarshal(rec.Body.Bytes(), &got); err != nil {
					t.Fatalf("decode body: %v", err)
				}
				want := temperature.FromCelsius(28.5)
				if got != want {
					t.Errorf("body = %+v, want %+v", got, want)
				}
				return
			}

			var got map[string]string
			if err := json.Unmarshal(rec.Body.Bytes(), &got); err != nil {
				t.Fatalf("decode body: %v", err)
			}
			if got["message"] != tt.wantBody {
				t.Errorf("message = %q, want %q", got["message"], tt.wantBody)
			}
		})
	}
}
