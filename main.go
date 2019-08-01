package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	en "github.com/mchmarny/gcputil/env"
	pj "github.com/mchmarny/gcputil/project"
)

var (
	logger = log.New(os.Stdout, "[TP] ", 0)

	projectID      = pj.GetIDOrFail()
	topic          = en.MustGetEnvVar("TOPIC", "tweets")
	consumerKey    = en.MustGetEnvVar("T_CONSUMER_KEY", "")
	consumerSecret = en.MustGetEnvVar("T_CONSUMER_SECRET", "")
	accessToken    = en.MustGetEnvVar("T_ACCESS_TOKEN", "")
	accessSecret   = en.MustGetEnvVar("T_ACCESS_SECRET", "")
)

func main() {

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, "OK")
	})

	http.HandleFunc("/query", queryHandler)

	port := fmt.Sprintf(":%s", en.MustGetEnvVar("PORT", "8080"))
	if err := http.ListenAndServe(port, nil); err != nil {
		logger.Fatal(err)
	}
}
