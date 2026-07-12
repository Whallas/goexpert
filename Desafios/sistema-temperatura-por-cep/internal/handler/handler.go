package handler

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/Whallas/goexpert/Desafios/sistema-temperatura-por-cep/internal/cep"
	"github.com/Whallas/goexpert/Desafios/sistema-temperatura-por-cep/internal/temperature"
	"github.com/Whallas/goexpert/Desafios/sistema-temperatura-por-cep/internal/weather"
)

// Weather handles GET /weather/{cep}, returning the current temperature for the
// city associated with the CEP in Celsius, Fahrenheit and Kelvin.
type Weather struct {
	Finder   cep.Finder
	Provider weather.Provider
}

// NewWeather wires the CEP finder and weather provider into a handler.
func NewWeather(finder cep.Finder, provider weather.Provider) *Weather {
	return &Weather{Finder: finder, Provider: provider}
}

func (h *Weather) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	code := r.PathValue("cep")

	city, err := h.Finder.City(r.Context(), code)
	switch {
	case errors.Is(err, cep.ErrInvalid):
		writeError(w, http.StatusUnprocessableEntity, cep.ErrInvalid.Error())
		return
	case errors.Is(err, cep.ErrNotFound):
		writeError(w, http.StatusNotFound, cep.ErrNotFound.Error())
		return
	case err != nil:
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}

	celsius, err := h.Provider.CelsiusFor(r.Context(), city)
	if err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}

	writeJSON(w, http.StatusOK, temperature.FromCelsius(celsius))
}

func writeJSON(w http.ResponseWriter, status int, body any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(body)
}

func writeError(w http.ResponseWriter, status int, message string) {
	writeJSON(w, status, map[string]string{"message": message})
}
