// server.go
package main

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	_ "github.com/mattn/go-sqlite3"
)

// Response from the external API
type ExchangeAPIResponse struct {
	USDBRL struct {
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
	} `json:"USDBRL"`
}

type ExchangeRate struct {
	Bid       string
	Timestamp time.Time
}

const (
	apiURL        = "https://economia.awesomeapi.com.br/json/last/USD-BRL"
	dbTimeout     = 10 * time.Millisecond
	apiTimeout    = 200 * time.Millisecond
	serverAddress = ":8080"
	dbFileName    = "exchange_rates.db"
)

func main() {
	// Initialize the database
	err := initializeDatabase()
	if err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}

	http.HandleFunc("/cotacao", getExchangeRateHandler)
	log.Println("Server started on", serverAddress)
	log.Fatal(http.ListenAndServe(serverAddress, nil))
}

func initializeDatabase() error {
	db, err := sql.Open("sqlite3", dbFileName)
	if err != nil {
		return err
	}
	defer db.Close()

	// Create the rates table if it doesn't exist
	_, err = db.Exec(`CREATE TABLE IF NOT EXISTS rates (
		bid TEXT,
		timestamp DATETIME
	)`)
	if err != nil {
		return fmt.Errorf("failed to create table: %v", err)
	}

	log.Println("Database initialized and rates table created (if not already exists)")
	return nil
}

func getExchangeRateHandler(w http.ResponseWriter, r *http.Request) {

	log.Printf("start getExchangeRateHandler")
	ctx, cancel := context.WithTimeout(r.Context(), apiTimeout)
	defer cancel()

	exchangeRate, err := fetchExchangeRate(ctx)
	if err != nil {
		http.Error(w, "Failed to fetch exchange rate", http.StatusInternalServerError)
		return
	}

	dbCtx, dbCancel := context.WithTimeout(context.Background(), dbTimeout)
	defer dbCancel()

	err = logExchangeRate(dbCtx, exchangeRate)
	if err != nil {
		log.Printf("Failed to Save in SQLite3")
		//http.Error(w, "Failed to log exchange rate", http.StatusInternalServerError)
		return
	}

	log.Printf("server exchangeRate: %v", exchangeRate.Bid)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{
		"bid": exchangeRate.Bid,
	})
}

func fetchExchangeRate(ctx context.Context) (*ExchangeRate, error) {
	req, err := http.NewRequestWithContext(ctx, "GET", apiURL, nil)
	if err != nil {
		return nil, err
	}

	log.Printf("start fetchExchangeRate")
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var apiResp ExchangeAPIResponse
	if err := json.NewDecoder(resp.Body).Decode(&apiResp); err != nil {
		return nil, err
	}

	log.Printf("server apiResp: %v", apiResp)
	return &ExchangeRate{
		Bid:       apiResp.USDBRL.Bid,
		Timestamp: time.Now(),
	}, nil
}

func logExchangeRate(ctx context.Context, rate *ExchangeRate) error {
	db, err := sql.Open("sqlite3", dbFileName)
	if err != nil {
		return err
	}
	defer db.Close()

	_, err = db.ExecContext(ctx, `INSERT INTO rates (bid, timestamp) VALUES (?, ?)`, rate.Bid, rate.Timestamp)
	return err
}
