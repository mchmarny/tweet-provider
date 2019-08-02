package main

import (
	"fmt"
	"log"
	"net"
	"net/http"
	"os"

	en "github.com/mchmarny/gcputil/env"
	pj "github.com/mchmarny/gcputil/project"
)

var (
	logger = log.New(os.Stdout, "", 0)

	projectID      = pj.GetIDOrFail()
	topic          = en.MustGetEnvVar("TOPIC", "search-tweets")
	port           = en.MustGetEnvVar("PORT", "8080")
	consumerKey    = en.MustGetEnvVar("T_CONSUMER_KEY", "")
	consumerSecret = en.MustGetEnvVar("T_CONSUMER_SECRET", "")
	accessToken    = en.MustGetEnvVar("T_ACCESS_TOKEN", "")
	accessSecret   = en.MustGetEnvVar("T_ACCESS_SECRET", "")
)

func main() {

	logger.Printf("Project: %s, Topic: %s, Port: %s", projectID, topic, port)

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, "OK")
	})

	http.HandleFunc("/query", queryHandler)

	hostPost := net.JoinHostPort("0.0.0.0", port)
	if err := http.ListenAndServe(hostPost, nil); err != nil {
		logger.Fatal(err)
	}
}
