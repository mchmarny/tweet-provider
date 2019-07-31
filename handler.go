package main

import (
	"encoding/json"
	"net/http"
	"time"
)

// requestHandler handles the HTTP request
func requestHandler(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Content-type", "application/json")

	var sr serviceRequest
	if err := json.NewDecoder(r.Body).Decode(&sr); err != nil {
		logger.Printf("Error decoding message: %v", err)
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode("Invalid Content")
		return
	}

	logger.Printf("Query: %v", sr)
	startTime := time.Now()
	lastID := work(sr.Query)

	msg := &serviceResponse{
		Query:    sr.Query,
		LastID:   lastID,
		Duration: time.Since(startTime),
	}

	logger.Printf("Result: %v", *msg)
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(msg)

	return
}

type serviceRequest struct {
	Query string `json:"query"`
}

type serviceResponse struct {
	Query    string `json:"query"`
	LastID   int64  `json:"lastID"`
	Duration time.Duration
}
