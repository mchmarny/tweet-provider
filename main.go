package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	en "github.com/mchmarny/gcputil/env"
)

var (
	logger  = log.New(os.Stdout, "[TP] ", 0)
)

func main() {
	http.HandleFunc("/", requestHandler)
	port := fmt.Sprintf(":%s", en.MustGetEnvVar("PORT", "8080"))
	if err := http.ListenAndServe(port, nil); err != nil {
		logger.Fatal(err)
	}
}
