package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"time"
)

type BrasilApiAddress struct {
	Cep          string `json:"cep"`
	State        string `json:"state"`
	City         string `json:"city"`
	Neighborhood string `json:"neighborhood"`
	Street       string `json:"street"`
	Service      string `json:"service"`
}

type ViaCEPAddress struct {
	Cep         string `json:"cep"`
	Logradouro  string `json:"logradouro"`
	Complemento string `json:"complemento"`
	Unidade     string `json:"unidade"`
	Bairro      string `json:"bairro"`
	Localidade  string `json:"localidade"`
	Uf          string `json:"uf"`
	Estado      string `json:"estado"`
	Regiao      string `json:"regiao"`
	Ibge        string `json:"ibge"`
	Gia         string `json:"gia"`
	Ddd         string `json:"ddd"`
	Siafi       string `json:"siafi"`
}

func getBrasilApiAddress(ctx context.Context, cep string) (*BrasilApiAddress, error) {
	url := "https://brasilapi.com.br/api/cep/v1/" + cep
	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return nil, err
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	res, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var data BrasilApiAddress
	err = json.Unmarshal(res, &data)
	if err != nil {
		return nil, err
	}

	return &data, nil
}

func getViaCEPAddress(ctx context.Context, cep string) (*ViaCEPAddress, error) {
	url := "https://viacep.com.br/ws/" + cep + "/json/"
	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return nil, err
	}
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	res, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var data ViaCEPAddress
	err = json.Unmarshal(res, &data)
	if err != nil {
		return nil, err
	}

	return &data, nil
}

func main() {
	if len(os.Args) < 2 {
		println("Cep is required. Operation Aborted.")
		return
	}
	cep := os.Args[1]

	brasilApiChan := make(chan *BrasilApiAddress, 1)
	viaCEPChan := make(chan *ViaCEPAddress, 1)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*1)
	defer cancel()

	go func() {
		address, err := getBrasilApiAddress(ctx, cep)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Erro ao buscar CEP na BrasilApi: %v\n", err)
			close(brasilApiChan)
			return
		}
		brasilApiChan <- address
	}()

	go func() {
		address, err := getViaCEPAddress(ctx, cep)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Erro ao buscar CEP na ViaCEP: %v\n", err)
			close(viaCEPChan)
			return
		}
		viaCEPChan <- address
	}()

	select {
	case vcAddress := <-viaCEPChan:
		println("ViaCep got first:")
		json.NewEncoder(os.Stdout).Encode(vcAddress)

	case baAddress := <-brasilApiChan:
		println("BrasilApi got first:")
		json.NewEncoder(os.Stdout).Encode(baAddress)

	case <-ctx.Done():
		println("timeout")
	}
}
