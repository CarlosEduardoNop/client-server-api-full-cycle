package main

import (
	"context"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"time"
)

func main() {
	ctx, cancel := context.WithTimeout(context.Background(), 300*time.Millisecond)
	defer cancel()

	req, err := http.NewRequestWithContext(ctx, "GET", "http://localhost:8080/cotacao", nil)

	if err != nil {
		log.Fatalf("Failed to create request: %v", err)
	}

	resp, err := http.DefaultClient.Do(req)

	if err != nil {
		log.Fatalf("Request failed: %v", err)
	}

	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)

	if err != nil {
		log.Fatalf("Failed to read response: %v", err)
	}

	err = ioutil.WriteFile("cotacao.txt", []byte(fmt.Sprintf("Dólar: %s", body)), 0644)

	if err != nil {
		log.Fatalf("Failed to write to file: %v", err)
	}
}
