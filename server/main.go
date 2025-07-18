package main

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"time"

	"github.com/google/uuid"
	_ "github.com/mattn/go-sqlite3"
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

type DollarResponse struct {
	USDBRL USDBRL `json:"USDBRL"`
}

const file string = "cotacao.db"

func main() {
	http.HandleFunc("/cotacao", HandlerDollar)

	fmt.Println("Server is running on port 8080")

	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		log.Fatal(err)
	}
}

func HandlerDollar(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/cotacao" {
		w.WriteHeader(http.StatusNotFound)
		return
	}
	dollar, error := GetDollarQuotation()
	if error != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	error = InsertCotacao(dollar)
	if error != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(error.Error())
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(dollar.Bid)
}

func GetDollarQuotation() (*USDBRL, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 200*time.Millisecond)
	defer cancel()

	req, err := http.NewRequestWithContext(ctx, "GET", "https://economia.awesomeapi.com.br/json/last/USD-BRL", nil)
	if err != nil {
		return nil, err
	}

	resp, error := http.DefaultClient.Do(req)
	if error != nil {
		return nil, error
	}
	defer resp.Body.Close()

	body, error := io.ReadAll(resp.Body)
	if error != nil {
		return nil, error
	}

	var dollarResponse DollarResponse
	error = json.Unmarshal(body, &dollarResponse)
	if error != nil {
		return nil, error
	}
	return &dollarResponse.USDBRL, nil
}

func InsertCotacao(usdbrl *USDBRL) error {
	db, err := sql.Open("sqlite3", file)
	if err != nil {
		return err
	}
	defer db.Close()

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Millisecond)
	defer cancel()

	insert := `
		INSERT INTO cotacoes (uuid, code, high, low, bid, created_at) 
		VALUES (?, ?, ?, ?, ?, ?)
	`
	if _, err := db.ExecContext(
		ctx,
		insert,
		uuid.New().String(),
		usdbrl.Code,
		usdbrl.High,
		usdbrl.Low,
		usdbrl.Bid,
		usdbrl.CreateDate); err != nil {
		return err
	}

	return nil
}
