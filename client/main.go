package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"time"
)

func main() {
	ctx, cancel := context.WithTimeout(context.Background(), 300*time.Millisecond)
	defer cancel()

	req, err := http.NewRequestWithContext(ctx, "GET", "http://localhost:8080/cotacao", nil)
	if err != nil {
		log.Fatal(err)
	}

	resp, error := http.DefaultClient.Do(req)
	if error != nil {
		log.Fatal(error)
	}
	defer resp.Body.Close()

	body, error := io.ReadAll(resp.Body)
	if error != nil {
		log.Fatal(error)
	}

	gravarCotacao(string(body))
}

func gravarCotacao(cotacao string) {
	file, err := os.Create("cotacao.txt")
	if err != nil {
		fmt.Fprintf(os.Stderr, "Erro ao criar arquivo: %v\n", err)
	}
	defer file.Close()
	_, err = file.WriteString(fmt.Sprintf("Dólar: %v", cotacao))
	if err != nil {
		fmt.Fprintf(os.Stderr, "Erro ao gravar cotacao: %v\n", err)
	}
	fmt.Println("Cotacao gravada com sucesso!")
}
