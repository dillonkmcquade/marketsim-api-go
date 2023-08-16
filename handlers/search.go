package handlers

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/jackc/pgx/v5/pgxpool"
)

type Ticker struct {
	Symbol      string `json:"symbol"`
	Description string `json:"description"`
}

type SearchResponse struct {
	Status int      `json:"status"`
	Data   []Ticker `json:"data"`
}

func Search(rw http.ResponseWriter, req *http.Request, pool *pgxpool.Pool) {
	name := req.URL.Query().Get("name")
	if name == "" {
		http.Error(rw, "No search query given", http.StatusBadRequest)
		return
	}
	ctx := req.Context()
	rows, err := pool.Query(ctx, `
        SELECT *
        FROM tickers
        WHERE LOWER(description) LIKE $1 || '%'
        LIMIT 10
    `, name)
	defer rows.Close()
	if err != nil {
		http.Error(rw, "Error retrieving resource from database", http.StatusInternalServerError)
		fmt.Println(err)
		return
	}
	var results []Ticker
	for rows.Next() {
		var ticker Ticker
		if err = rows.Scan(&ticker.Symbol, &ticker.Description); err != nil {
			log.Println(err)
		}
		results = append(results, ticker)
	}
	if err = rows.Err(); err != nil {
		log.Println(err)
	}
	if len(results) == 0 {
		http.Error(rw, "No results", http.StatusBadRequest)
		return
	}

	encoder := json.NewEncoder(rw)

	err = encoder.Encode(SearchResponse{Data: results, Status: 200})
	if err != nil {
		log.Println(err)
	}
}
