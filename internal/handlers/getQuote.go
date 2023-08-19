package handlers

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
)

type Quote struct {
	C  float32 `json:"c"`
	H  float32 `json:"h"`
	L  float32 `json:"l"`
	O  float32 `json:"o"`
	PC float32 `json:"pc"`
	DP float32 `json:"dp"`
	D  float32 `json:"d"`
	T  int64   `json:"t"`
}

type QuoteResponse struct {
	Status int32 `json:"status"`
	Data   Quote `json:"data"`
}

// Returns a quote for the given symbol
// Expects query parameter 'symbol'
//
// Ex. .../stock/quote?symbol=TSLA
func GetQuote(rw http.ResponseWriter, r *http.Request) {
	symbol := r.URL.Query().Get("symbol")

	key := os.Getenv("FINNHUB_KEY")

	if key == "" {
		log.Fatal("FINNHUB_KEY is undefined")
	}

	// insert api key and symbol to external api url
	url := fmt.Sprintf("https://finnhub.io/api/v1/quote?symbol=%s&token=%s", symbol, key)

	// make GET request to external api
	res, err := http.Get(url)
	if err != nil {
		http.Error(rw, "Error fetching ticker", http.StatusInternalServerError)
		return
	}
	defer res.Body.Close()

	var quoteResponse QuoteResponse

	err = json.NewDecoder(res.Body).Decode(&quoteResponse.Data)

	if err != nil {
		http.Error(rw, "Error fetching ticker", http.StatusInternalServerError)
		return
	}

	// return json
	quoteResponse.Status = 200
	rw.Header().Add("Content-Type", "application/json")
	err = json.NewEncoder(rw).Encode(quoteResponse)
	if err != nil {
		fmt.Println(err)
	}
}
