package weather

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestWeatherAPI_CelsiusFor(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if got := r.URL.Query().Get("q"); got != "São Paulo" {
				t.Errorf("query q = %q, want %q", got, "São Paulo")
			}
			_, _ = w.Write([]byte(`{"current":{"temp_c":28.5}}`))
		}))
		defer server.Close()

		api := NewWeatherAPI("test-key", server.Client())
		api.BaseURL = server.URL

		got, err := api.CelsiusFor(context.Background(), "São Paulo")
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if got != 28.5 {
			t.Errorf("temp = %v, want 28.5", got)
		}
	})

	t.Run("non-200 status", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusBadRequest)
		}))
		defer server.Close()

		api := NewWeatherAPI("test-key", server.Client())
		api.BaseURL = server.URL

		if _, err := api.CelsiusFor(context.Background(), "Nowhere"); err == nil {
			t.Error("expected error, got nil")
		}
	})
}
