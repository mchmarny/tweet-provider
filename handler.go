package main

import (
	"encoding/json"
	"net/http"
	"time"
)

func queryHandler(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Content-type", "application/json")

	var c serviceRequest
	if err := json.NewDecoder(r.Body).Decode(&c); err != nil {
		logger.Printf("Error decoding message: %v", err)
		writeResp(w, http.StatusBadRequest, "Invalid Content")
		return
	}
	logger.Printf("Query: %v", c)

	startTime := time.Now()
	lastID := work(c.Query)

	msg := &serviceResponse{
		Query:    c.Query,
		LastID:   lastID,
		Duration: time.Since(startTime),
	}

	logger.Printf("Result: %v", *msg)
	writeResp(w, http.StatusOK, msg)

	return
}

func writeResp(w http.ResponseWriter, status int, msg interface{}) {
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(msg)
}

type serviceRequest struct {
	Query string `json:"query"`
}

type serviceResponse struct {
	Query    string `json:"query"`
	LastID   int64  `json:"lastID"`
	Duration time.Duration
}
