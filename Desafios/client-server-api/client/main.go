package main

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"os"
	"time"
)

func saveFile(data []byte) error {
	f, err := os.Create("./cotacao.txt")
	if err != nil {
		panic(err)
	}

	_, err = f.Write(data)
	if err != nil {
		return err
	}

	f.Close()

	return nil
}

func main() {
	ctx, cancel := context.WithTimeout(context.Background(), time.Millisecond*300)
	defer cancel()

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, "http://localhost:8080/cotacao", nil)
	if err != nil {
		panic(err)
	}

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		panic(err)
	}
	defer res.Body.Close()

	data, err := io.ReadAll(res.Body)
	if err != nil {
		panic(err)
	}

	fmt.Println(string(data))

	if res.StatusCode != http.StatusOK && res.StatusCode != http.StatusCreated {
		panic("Request failed")
	}

	err = saveFile(data)
	if err != nil {
		panic(err)
	}

	fmt.Printf("Arquivo salvo com sucesso.\n")
}
