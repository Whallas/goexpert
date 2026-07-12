package weather

import (
	"fmt"
	"net/http"
)

// Provider names accepted by New.
const (
	ProviderOpenMeteo  = "openmeteo"
	ProviderWeatherAPI = "weatherapi"
)

// DefaultProvider is used when no provider is configured.
const DefaultProvider = ProviderOpenMeteo

// New builds a Provider from a provider name. The apiKey is only required by
// the weatherapi provider; openmeteo ignores it. An empty name selects the
// DefaultProvider.
func New(name, apiKey string, client *http.Client) (Provider, error) {
	switch name {
	case "", ProviderOpenMeteo:
		return NewOpenMeteo(client), nil
	case ProviderWeatherAPI:
		if apiKey == "" {
			return nil, fmt.Errorf("weather: WEATHER_API_KEY is required for the %q provider", ProviderWeatherAPI)
		}
		return NewWeatherAPI(apiKey, client), nil
	default:
		return nil, fmt.Errorf("weather: unknown provider %q", name)
	}
}
