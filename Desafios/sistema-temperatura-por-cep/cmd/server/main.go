package main

import (
	"log"
	"net/http"
	"os"
	"time"

	"github.com/Whallas/goexpert/Desafios/sistema-temperatura-por-cep/internal/cep"
	"github.com/Whallas/goexpert/Desafios/sistema-temperatura-por-cep/internal/handler"
	"github.com/Whallas/goexpert/Desafios/sistema-temperatura-por-cep/internal/weather"
)

func main() {
	// WEATHER_PROVIDER selects the weather backend: "openmeteo" (default) or
	// "weatherapi". WEATHER_API_KEY is only required by "weatherapi".
	providerName := os.Getenv("WEATHER_PROVIDER")
	if providerName == "" {
		providerName = weather.DefaultProvider
	}
	apiKey := os.Getenv("WEATHER_API_KEY")

	// Cloud Run injects the listening port via PORT.
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	httpClient := &http.Client{Timeout: 10 * time.Second}
	finder := cep.NewViaCEP(httpClient)
	provider, err := weather.New(providerName, apiKey, httpClient)
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("weather provider: %s", providerName)

	mux := http.NewServeMux()
	mux.Handle("GET /weather/{cep}", handler.NewWeather(finder, provider))
	mux.HandleFunc("GET /health", func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte("ok"))
	})

	server := &http.Server{
		Addr:         ":" + port,
		Handler:      mux,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
	}

	log.Printf("listening on :%s", port)
	if err := server.ListenAndServe(); err != nil {
		log.Fatal(err)
	}
}
