package weather

import (
	"net/http"
	"testing"
)

func TestNew(t *testing.T) {
	tests := []struct {
		name     string
		provider string
		apiKey   string
		wantType any
		wantErr  bool
	}{
		{name: "default is openmeteo", provider: "", wantType: &OpenMeteo{}},
		{name: "explicit openmeteo", provider: ProviderOpenMeteo, wantType: &OpenMeteo{}},
		{name: "weatherapi with key", provider: ProviderWeatherAPI, apiKey: "k", wantType: &WeatherAPI{}},
		{name: "weatherapi without key", provider: ProviderWeatherAPI, wantErr: true},
		{name: "unknown provider", provider: "bogus", wantErr: true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := New(tt.provider, tt.apiKey, http.DefaultClient)
			if tt.wantErr {
				if err == nil {
					t.Fatal("expected error, got nil")
				}
				return
			}
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			switch tt.wantType.(type) {
			case *OpenMeteo:
				if _, ok := got.(*OpenMeteo); !ok {
					t.Errorf("got %T, want *OpenMeteo", got)
				}
			case *WeatherAPI:
				if _, ok := got.(*WeatherAPI); !ok {
					t.Errorf("got %T, want *WeatherAPI", got)
				}
			}
		})
	}
}
