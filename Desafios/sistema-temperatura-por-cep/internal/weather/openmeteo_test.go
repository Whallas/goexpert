package weather

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestOpenMeteo_CelsiusFor(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		geo := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if got := r.URL.Query().Get("name"); got != "São Paulo" {
				t.Errorf("name = %q, want %q", got, "São Paulo")
			}
			_, _ = w.Write([]byte(`{"results":[{"latitude":-23.55,"longitude":-46.63}]}`))
		}))
		defer geo.Close()

		forecast := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			_, _ = w.Write([]byte(`{"current":{"temperature_2m":28.5}}`))
		}))
		defer forecast.Close()

		om := NewOpenMeteo(http.DefaultClient)
		om.GeocodeURL = geo.URL
		om.ForecastURL = forecast.URL

		got, err := om.CelsiusFor(context.Background(), "São Paulo")
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if got != 28.5 {
			t.Errorf("temp = %v, want 28.5", got)
		}
	})

	t.Run("location not geocoded", func(t *testing.T) {
		geo := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			_, _ = w.Write([]byte(`{"results":[]}`))
		}))
		defer geo.Close()

		om := NewOpenMeteo(http.DefaultClient)
		om.GeocodeURL = geo.URL
		om.ForecastURL = "http://unused"

		if _, err := om.CelsiusFor(context.Background(), "Nowhere"); err == nil {
			t.Error("expected error, got nil")
		}
	})
}
