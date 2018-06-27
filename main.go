package main

import (
	"flag"
	"log"
	"os"
	"os/signal"
	"syscall"

	"golang.org/x/net/context"
)

const (
	defaultQueryFrequency = 5 //sec
	defaultPubSubTopic    = "tweets"
)

var (
	logger = log.New(os.Stdout, "[app] ", log.Lshortfile|log.Ldate|log.Ltime)
	debug  = log.New(os.Stdout, "[tweet] ", 0)
)

func main() {

	// GCP
	projectID := flag.String("project", os.Getenv("GCLOUD_PROJECT"), "GCP Project ID")
	topicName := flag.String("topic", defaultPubSubTopic, "Google PubSub topic name [tweets]")

	// twitter
	consumerKey := flag.String("consumerKey", os.Getenv("T_CONSUMER_KEY"), "Twitter consumer key")
	consumerSecret := flag.String("consumerSecret", os.Getenv("T_CONSUMER_SECRET"), "Twitter consumer secret")
	accessToken := flag.String("accessToken", os.Getenv("T_ACCESS_TOKEN"), "Twitter access token")
	accessSecret := flag.String("accessSecret", os.Getenv("T_ACCESS_SECRET"), "Twitter access secret")

	// app
	queryFreqInMin := flag.Int("queryFreqInMin", defaultQueryFrequency, "Frequency of Twitter queries")
	searchQuery := flag.String("query", os.Getenv("T_QUERY"), "Search Query")

	flag.Parse()

	// VALIDATION
	if *projectID == "" {
		logger.Fatalf("Missing required app configs: project:%v", projectID)
	}
	if *consumerKey == "" || *consumerSecret == "" || *accessToken == "" || *accessSecret == "" {
		logger.Fatalf("Missing required twitter configs: consumerKey:%v, consumerSecret:%v accessToken:%v, accessSecret:%v",
			consumerKey, consumerSecret, accessToken, accessSecret)
	}
	if *searchQuery == "" {
		logger.Fatal("Missing required search query")
	}

	// init context
	appContext, ctxCancel := context.WithCancel(context.Background())

	// initialize dependencies
	publisher, err := getPublisher(&appContext, *projectID, *topicName)
	if err != nil {
		logger.Fatalf("Error while creating publisher: %v", err)
	}

	go func() {
		// Wait for SIGINT and SIGTERM (HIT CTRL-C)
		ch := make(chan os.Signal)
		signal.Notify(ch, syscall.SIGINT, syscall.SIGTERM)
		logger.Println(<-ch)
		ctxCancel()
		os.Exit(0)
	}()

	messages := make(chan *[]byte, 1)

	// configure provider
	provider := getProvider(*queryFreqInMin, *consumerKey, *consumerSecret, *accessToken, *accessSecret)
	go func() {
		provider.Provide(*searchQuery, messages)
	}()

	for {
		select {
		case <-appContext.Done():
			break
		case m := <-messages:
			publisher.Publish(m)
		}
	}

}
