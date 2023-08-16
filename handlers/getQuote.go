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

func GetQuote(rw http.ResponseWriter, r *http.Request) {
	symbol := r.URL.Query().Get("symbol")

	key := os.Getenv("FINNHUB_KEY")

	if key == "" {
		log.Fatal("FINNHUB_KEY is undefined")
	}

	url := fmt.Sprintf("https://finnhub.io/api/v1/quote?symbol=%s&token=%s", symbol, key)

	res, err := http.Get(url)
	if err != nil {
		http.Error(rw, "Error fetching ticker", http.StatusInternalServerError)
		return
	}

	defer res.Body.Close()

	var data Quote

	err = json.NewDecoder(res.Body).Decode(&data)

	if err != nil {
		http.Error(rw, "Error fetching ticker", http.StatusInternalServerError)
		return
	}

	err = json.NewEncoder(rw).Encode(data)
	if err != nil {
		fmt.Println(err)
	}
}
