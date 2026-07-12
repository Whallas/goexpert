package weather

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
)

// Provider returns the current temperature in Celsius for a given location.
type Provider interface {
	CelsiusFor(ctx context.Context, location string) (float64, error)
}

// WeatherAPI is a Provider backed by https://www.weatherapi.com.
type WeatherAPI struct {
	BaseURL string
	APIKey  string
	Client  *http.Client
}

// NewWeatherAPI builds a WeatherAPI client. A nil HTTP client falls back to
// http.DefaultClient.
func NewWeatherAPI(apiKey string, client *http.Client) *WeatherAPI {
	if client == nil {
		client = http.DefaultClient
	}
	return &WeatherAPI{
		BaseURL: "https://api.weatherapi.com/v1",
		APIKey:  apiKey,
		Client:  client,
	}
}

type weatherResponse struct {
	Current struct {
		TempC float64 `json:"temp_c"`
	} `json:"current"`
}

// CelsiusFor queries the current temperature for the given location name.
func (w *WeatherAPI) CelsiusFor(ctx context.Context, location string) (float64, error) {
	endpoint := fmt.Sprintf("%s/current.json?key=%s&q=%s",
		w.BaseURL, url.QueryEscape(w.APIKey), url.QueryEscape(location))

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, endpoint, nil)
	if err != nil {
		return 0, err
	}

	resp, err := w.Client.Do(req)
	if err != nil {
		return 0, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return 0, fmt.Errorf("weather api returned status %d for %q", resp.StatusCode, location)
	}

	var data weatherResponse
	if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
		return 0, err
	}

	return data.Current.TempC, nil
}
