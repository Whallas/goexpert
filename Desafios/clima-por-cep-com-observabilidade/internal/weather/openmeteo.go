package weather

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
)

// OpenMeteo is a Provider backed by https://open-meteo.com.
//
// Open-Meteo needs coordinates rather than a city name, so each lookup first
// resolves the location through the Open-Meteo geocoding API and then queries
// the forecast endpoint. No API key is required.
type OpenMeteo struct {
	GeocodeURL  string
	ForecastURL string
	Client      *http.Client
}

// NewOpenMeteo builds an OpenMeteo provider. A nil HTTP client falls back to
// http.DefaultClient.
func NewOpenMeteo(client *http.Client) *OpenMeteo {
	if client == nil {
		client = http.DefaultClient
	}
	return &OpenMeteo{
		GeocodeURL:  "https://geocoding-api.open-meteo.com/v1/search",
		ForecastURL: "https://api.open-meteo.com/v1/forecast",
		Client:      client,
	}
}

type geocodeResponse struct {
	Results []struct {
		Latitude  float64 `json:"latitude"`
		Longitude float64 `json:"longitude"`
	} `json:"results"`
}

type forecastResponse struct {
	Current struct {
		Temperature float64 `json:"temperature_2m"`
	} `json:"current"`
}

// CelsiusFor resolves the location to coordinates and returns its current
// temperature in Celsius.
func (o *OpenMeteo) CelsiusFor(ctx context.Context, location string) (float64, error) {
	lat, lon, err := o.geocode(ctx, location)
	if err != nil {
		return 0, err
	}
	return o.forecast(ctx, lat, lon)
}

func (o *OpenMeteo) geocode(ctx context.Context, location string) (lat, lon float64, err error) {
	endpoint := fmt.Sprintf("%s?name=%s&count=1&language=pt&format=json",
		o.GeocodeURL, url.QueryEscape(location))

	var data geocodeResponse
	if err := o.get(ctx, endpoint, &data); err != nil {
		return 0, 0, err
	}
	if len(data.Results) == 0 {
		return 0, 0, fmt.Errorf("open-meteo: no coordinates found for %q", location)
	}
	return data.Results[0].Latitude, data.Results[0].Longitude, nil
}

func (o *OpenMeteo) forecast(ctx context.Context, lat, lon float64) (float64, error) {
	endpoint := fmt.Sprintf("%s?latitude=%f&longitude=%f&current=temperature_2m",
		o.ForecastURL, lat, lon)

	var data forecastResponse
	if err := o.get(ctx, endpoint, &data); err != nil {
		return 0, err
	}
	return data.Current.Temperature, nil
}

func (o *OpenMeteo) get(ctx context.Context, endpoint string, out any) error {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, endpoint, nil)
	if err != nil {
		return err
	}

	resp, err := o.Client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("open-meteo returned status %d", resp.StatusCode)
	}

	return json.NewDecoder(resp.Body).Decode(out)
}
