// Serviço A (Input): recebe a requisição do usuário via POST, valida o CEP e,
// se válido, encaminha para o Serviço B propagando o contexto de tracing.
package main

import (
	"context"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/Whallas/goexpert/Desafios/clima-por-cep-com-observabilidade/internal/tracing"

	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/trace"
)

type cepRequest struct {
	Cep string `json:"cep"`
}

func main() {
	ctx := context.Background()

	shutdown, err := tracing.Init(ctx, "service-a", getenv("OTEL_EXPORTER_OTLP_ENDPOINT", "otel-collector:4318"))
	if err != nil {
		log.Fatalf("falha ao inicializar tracing: %v", err)
	}
	defer shutdown(context.Background())

	client := &http.Client{
		Transport: otelhttp.NewTransport(http.DefaultTransport),
		Timeout:   15 * time.Second,
	}
	serviceBURL := getenv("SERVICE_B_URL", "http://service-b:8081")
	tracer := otel.Tracer("service-a")

	mux := http.NewServeMux()
	mux.Handle("/", handler(client, serviceBURL, tracer))

	addr := ":" + getenv("PORT", "8080")
	log.Printf("service-a ouvindo em %s", addr)
	if err := http.ListenAndServe(addr, otelhttp.NewHandler(mux, "service-a")); err != nil {
		log.Fatal(err)
	}
}

func handler(client *http.Client, serviceBURL string, tracer trace.Tracer) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
			return
		}

		var body cepRequest
		if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
			invalidZipcode(w)
			return
		}
		if !isValidCep(body.Cep) {
			invalidZipcode(w)
			return
		}

		ctx, span := tracer.Start(r.Context(), "encaminha-servico-b")
		defer span.End()

		req, err := http.NewRequestWithContext(ctx, http.MethodGet, serviceBURL+"/"+body.Cep, nil)
		if err != nil {
			http.Error(w, "internal error", http.StatusInternalServerError)
			return
		}

		resp, err := client.Do(req)
		if err != nil {
			http.Error(w, "error contacting service B", http.StatusInternalServerError)
			return
		}
		defer resp.Body.Close()

		data, _ := io.ReadAll(resp.Body)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(resp.StatusCode)
		w.Write(data)
	}
}

// isValidCep garante que o CEP seja uma string com exatamente 8 dígitos.
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

func invalidZipcode(w http.ResponseWriter) {
	w.WriteHeader(http.StatusUnprocessableEntity)
	w.Write([]byte("invalid zipcode"))
}

func getenv(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}
