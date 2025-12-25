package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

type USDBRL struct {
	Code       string `json:"code"`
	Codein     string `json:"codein"`
	Name       string `json:"name"`
	High       string `json:"high"`
	Low        string `json:"low"`
	VarBid     string `json:"varBid"`
	PctChange  string `json:"pctChange"`
	Bid        string `json:"bid"`
	Ask        string `json:"ask"`
	Timestamp  string `json:"timestamp"`
	CreateDate string `json:"create_date"`
}

type DolarResponse struct {
	USDBRL USDBRL `json:"USDBRL"`
}

type Cotacao struct {
	ID uint `gorm:"primarykey"`
	USDBRL
	CreatedAt time.Time
	UpdatedAt time.Time
}

func getDolarPrice(ctx context.Context) (*DolarResponse, error) {
	ctx, cancel := context.WithTimeout(ctx, 200*time.Millisecond)
	defer cancel()

	url := "https://economia.awesomeapi.com.br/json/last/USD-BRL"
	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return nil, err
	}

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	data, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	var dolarResponse DolarResponse
	err = json.Unmarshal(data, &dolarResponse)
	if err != nil {
		return nil, err
	}

	return &dolarResponse, nil
}

func writeInternalError(w http.ResponseWriter) {
	w.WriteHeader(http.StatusInternalServerError)
	json.NewEncoder(w).Encode(
		map[string]string{"error": "erro interno"},
	)
}

func storeCotacao(ctx context.Context, db *gorm.DB, dolarResponse *DolarResponse) (Cotacao, error) {
	ctx, cancel := context.WithTimeout(ctx, 10*time.Millisecond)
	defer cancel()

	cotacao := Cotacao{
		USDBRL: dolarResponse.USDBRL,
	}
	err := db.WithContext(ctx).Create(&cotacao).Error

	return cotacao, err
}

func cotacaoController(w http.ResponseWriter, r *http.Request, db *gorm.DB) {
	ctx := r.Context()
	defer fmt.Println("Request has been ended")

	if ctx.Err() != nil {
		fmt.Println("Request cancelled")
		return
	}

	w.Header().Set("Content-Type", "application/json")

	dolarResponse, err := getDolarPrice(ctx)
	if err != nil {
		fmt.Println("Falha ao obter os dados da API", err)
		writeInternalError(w)
		return
	}

	cotacao, err := storeCotacao(ctx, db, dolarResponse)
	if err != nil {
		fmt.Println("Falha ao criar a cotação no banco de dados", err)
		writeInternalError(w)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(
		map[string]string{"bid": cotacao.Bid},
	)
}

func main() {
	dsn := "cotacoes.db"
	db, err := gorm.Open(sqlite.Open(dsn), &gorm.Config{})
	if err != nil {
		fmt.Println("Falha ao tentar se conectar com o sqlite", err)
		return
	}

	err = db.AutoMigrate(&Cotacao{})
	if err != nil {
		fmt.Println("Falha ao carregar o sqlite", err)
		return
	}

	http.HandleFunc("/cotacao", func(w http.ResponseWriter, r *http.Request) {
		cotacaoController(w, r, db)
	})
	http.ListenAndServe(":8080", nil)
}
