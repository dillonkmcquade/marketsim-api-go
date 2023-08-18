package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
)

type Candle struct{}

func GetCandle(rw http.ResponseWriter, r *http.Request) {
	symbol := r.URL.Query().Get("symbol")
	from := r.URL.Query().Get("from")
	key := os.Getenv("FINNHUB_KEY")
	if symbol == "" || from == "" || key == "" {
		http.Error(rw, "Missing query parameters", http.StatusBadRequest)
		return
	}

	url := fmt.Sprintf("https://finnhub.io/api/v1/stock/candle?symbol=%s&resolution=1&from=%s&to=1679649780&token=%s", symbol, from, key)

	res, err := http.Get(url)
	if err != nil {
		fmt.Println(err)
		http.Error(rw, "Error fetching candlestick data", http.StatusInternalServerError)
		return
	}

	var result Candle
	err = json.NewDecoder(res.Body).Decode(&result)
	if err != nil {
		fmt.Println(err)
		http.Error(rw, "Error processing candlestick data", http.StatusInternalServerError)
		return
	}

	rw.WriteHeader(200)
	err = json.NewEncoder(rw).Encode(result)
	if err != nil {
		fmt.Println(err)
	}
}
