package cep

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"regexp"
)

// ErrInvalid is returned when a CEP does not match the expected 8-digit format.
var ErrInvalid = errors.New("invalid zipcode")

// ErrNotFound is returned when a well-formed CEP is not present in the database.
var ErrNotFound = errors.New("can not find zipcode")

var cepPattern = regexp.MustCompile(`^[0-9]{8}$`)

// Finder resolves a CEP into the name of its city.
type Finder interface {
	City(ctx context.Context, cep string) (string, error)
}

// ViaCEP is a Finder backed by the public https://viacep.com.br service.
type ViaCEP struct {
	BaseURL string
	Client  *http.Client
}

// NewViaCEP builds a ViaCEP client using the given HTTP client. A nil client
// falls back to http.DefaultClient.
func NewViaCEP(client *http.Client) *ViaCEP {
	if client == nil {
		client = http.DefaultClient
	}
	return &ViaCEP{BaseURL: "https://viacep.com.br/ws", Client: client}
}

type viaCEPResponse struct {
	Localidade string `json:"localidade"`
	Erro       any    `json:"erro"`
}

// City returns the city for the given CEP. It returns ErrInvalid for malformed
// input and ErrNotFound when the CEP does not exist.
func (v *ViaCEP) City(ctx context.Context, cep string) (string, error) {
	if !cepPattern.MatchString(cep) {
		return "", ErrInvalid
	}

	url := fmt.Sprintf("%s/%s/json/", v.BaseURL, cep)
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return "", err
	}

	resp, err := v.Client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", ErrNotFound
	}

	var data viaCEPResponse
	if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
		return "", err
	}

	// ViaCEP answers 200 with `{ "erro": true }` for unknown CEPs.
	if data.Erro != nil || data.Localidade == "" {
		return "", ErrNotFound
	}

	return data.Localidade, nil
}
