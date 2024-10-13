// client.go
package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"time"
)

const (
	serverURL      = "http://server:8080/cotacao" // Server container hostname in Docker Compose
	clientTimeout  = 300 * time.Millisecond
	outputFileName = "cotacao.txt"
)

func main() {
	ctx, cancel := context.WithTimeout(context.Background(), clientTimeout)
	defer cancel()

	exchangeRate, err := getExchangeRate(ctx)
	if err != nil {
		log.Fatalf("Error fetching exchange rate: %v", err)
	}

	err = saveToFile(exchangeRate)
	if err != nil {
		log.Fatalf("Error saving exchange rate to file: %v", err)
	}

	fmt.Printf("Exchange rate saved to %s: Dollar: %s\n", outputFileName, exchangeRate)
}

func getExchangeRate(ctx context.Context) (string, error) {
	req, err := http.NewRequestWithContext(ctx, "GET", serverURL, nil)
	if err != nil {
		log.Printf("Req err: %v", err)
		return "", err
	}

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Printf("resp err: %v", err)
		return "", err
	}
	defer resp.Body.Close()

	log.Printf("resp.Status: %v", resp.Status)
	log.Printf("resp.Body: %v", resp.Body)

	var data map[string]string
	if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
		log.Printf("NewDecoder err: %v", err)
		return "", err
	}

	return data["bid"], nil
}

func saveToFile(rate string) error {
	fileContent := fmt.Sprintf("Dollar: %s", rate)
	return ioutil.WriteFile(outputFileName, []byte(fileContent), 0644)
}
