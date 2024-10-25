package main

import (
	"context"
	"database/sql"
	"encoding/json"
	"io/ioutil"
	"log"
	_ "modernc.org/sqlite"
	"net/http"
	"time"
)

type Response struct {
	USDBRL USDBRL `json:"USDBRL"`
}

type USDBRL struct {
	Bid string `json:"bid"`
}

func main() {
	http.HandleFunc("/cotacao", quoteHandler)
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func quoteHandler(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), 200*time.Millisecond)
	defer cancel()

	req, err := http.NewRequestWithContext(ctx, "GET", "https://economia.awesomeapi.com.br/json/last/USD-BRL", nil)
	if err != nil {
		http.Error(w, "Request failed", http.StatusInternalServerError)
		log.Println("Error creating request:", err)
		return
	}

	resp, err := http.DefaultClient.Do(req)

	if err != nil {
		http.Error(w, "Failed to fetch cotacao", http.StatusInternalServerError)
		log.Println("Error fetching cotacao:", err)
		return
	}

	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)

	if err != nil {
		log.Fatal(err)
	}
	var response Response

	err = json.Unmarshal(body, &response)

	if err != nil {
		log.Fatal(err)
	}

	bid := response.USDBRL.Bid

	_, err = w.Write([]byte(bid))

	if err != nil {
		return
	}

	go insertQuote(bid)
}

func insertQuote(bid string) {
	db, err := sql.Open("sqlite", "cotacoes.db")

	if err != nil {
		log.Println("Failed to open database:", err)
		return
	}

	defer db.Close()

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Millisecond)

	defer cancel()

	_, err = db.ExecContext(ctx, "CREATE TABLE IF NOT EXISTS cotacao (id INTEGER PRIMARY KEY, bid TEXT, timestamp DATETIME DEFAULT CURRENT_TIMESTAMP)")

	if err != nil {
		log.Println("Failed to create table:", err)
		return
	}

	_, err = db.ExecContext(ctx, "INSERT INTO cotacao (bid) VALUES (?)", bid)

	if err != nil {
		log.Println("Failed to insert cotacao:", err)
	}
}
