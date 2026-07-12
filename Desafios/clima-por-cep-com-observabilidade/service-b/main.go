// Serviço B (Orquestração): recebe um CEP válido, descobre a cidade via ViaCEP,
// consulta a temperatura no provedor de clima configurado e devolve em Celsius,
// Fahrenheit e Kelvin.
package main

import (
	"context"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/Whallas/goexpert/Desafios/clima-por-cep-com-observabilidade/internal/tracing"
	"github.com/Whallas/goexpert/Desafios/clima-por-cep-com-observabilidade/internal/weather"

	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/trace"
)

type weatherResponse struct {
	City  string  `json:"city"`
	TempC float64 `json:"temp_C"`
	TempF float64 `json:"temp_F"`
	TempK float64 `json:"temp_K"`
}

var httpClient = &http.Client{
	Transport: otelhttp.NewTransport(http.DefaultTransport),
	Timeout:   15 * time.Second,
}

func main() {
	ctx := context.Background()

	shutdown, err := tracing.Init(ctx, "service-b", getenv("OTEL_EXPORTER_OTLP_ENDPOINT", "otel-collector:4318"))
	if err != nil {
		log.Fatalf("falha ao inicializar tracing: %v", err)
	}
	defer shutdown(context.Background())

	providerName := getenv("WEATHER_PROVIDER", weather.DefaultProvider)
	provider, err := weather.New(providerName, os.Getenv("WEATHER_API_KEY"), httpClient)
	if err != nil {
		log.Fatalf("falha ao configurar provedor de clima: %v", err)
	}
	log.Printf("provedor de clima: %s", providerName)

	tracer := otel.Tracer("service-b")

	mux := http.NewServeMux()
	mux.Handle("/", handler(tracer, provider))

	addr := ":" + getenv("PORT", "8081")
	log.Printf("service-b ouvindo em %s", addr)
	if err := http.ListenAndServe(addr, otelhttp.NewHandler(mux, "service-b")); err != nil {
		log.Fatal(err)
	}
}

func handler(tracer trace.Tracer, provider weather.Provider) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		cep := strings.Trim(r.URL.Path, "/")
		if !isValidCep(cep) {
			writeText(w, http.StatusUnprocessableEntity, "invalid zipcode")
			return
		}

		ctx := r.Context()

		city, found, err := fetchCity(ctx, tracer, cep)
		if err != nil {
			writeText(w, http.StatusInternalServerError, "error looking up zipcode")
			return
		}
		if !found {
			writeText(w, http.StatusNotFound, "can not find zipcode")
			return
		}

		tempC, err := fetchTemperature(ctx, tracer, provider, city)
		if err != nil {
			writeText(w, http.StatusInternalServerError, "error looking up temperature")
			return
		}

		resp := weatherResponse{
			City:  city,
			TempC: tempC,
			TempF: tempC*1.8 + 32,
			TempK: tempC + 273,
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(resp)
	}
}

// fetchCity consulta o ViaCEP. Span manual "busca-cep" mede o tempo da chamada.
func fetchCity(ctx context.Context, tracer trace.Tracer, cep string) (string, bool, error) {
	ctx, span := tracer.Start(ctx, "busca-cep")
	defer span.End()

	endpoint := "https://viacep.com.br/ws/" + cep + "/json/"
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, endpoint, nil)
	if err != nil {
		return "", false, err
	}
	resp, err := httpClient.Do(req)
	if err != nil {
		return "", false, err
	}
	defer resp.Body.Close()

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", false, err
	}

	var result struct {
		Localidade string          `json:"localidade"`
		Erro       json.RawMessage `json:"erro"`
	}
	if err := json.Unmarshal(data, &result); err != nil {
		return "", false, err
	}

	// ViaCEP devolve {"erro": true} (ou ausência de localidade) quando não acha.
	if len(result.Erro) > 0 || result.Localidade == "" {
		return "", false, nil
	}
	return result.Localidade, true, nil
}

// fetchTemperature delega ao provedor de clima configurado. Span manual
// "busca-temperatura" mede o tempo total da consulta de temperatura.
func fetchTemperature(ctx context.Context, tracer trace.Tracer, provider weather.Provider, city string) (float64, error) {
	ctx, span := tracer.Start(ctx, "busca-temperatura")
	defer span.End()
	return provider.CelsiusFor(ctx, city)
}

func isValidCep(cep string) bool {
	if len(cep) != 8 {
		return false
	}
	for _, c := range cep {
		if c < '0' || c > '9' {
			return false
		}
	}
	return true
}

func writeText(w http.ResponseWriter, status int, msg string) {
	w.WriteHeader(status)
	w.Write([]byte(msg))
}

func getenv(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}
