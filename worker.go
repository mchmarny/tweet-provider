package main

import (
	"context"
	"encoding/json"

	"github.com/dghubble/go-twitter/twitter"
	"github.com/dghubble/oauth1"

	en "github.com/mchmarny/gcputil/env"
	pj "github.com/mchmarny/gcputil/project"
)

var (
	project        = pj.GetIDOrFail()
	topic          = en.MustGetEnvVar("TOPIC", "tweets")
	consumerKey    = en.MustGetEnvVar("T_CONSUMER_KEY", "")
	consumerSecret = en.MustGetEnvVar("T_CONSUMER_SECRET", "")
	accessToken    = en.MustGetEnvVar("T_ACCESS_TOKEN", "")
	accessSecret   = en.MustGetEnvVar("T_ACCESS_SECRET", "")
)

func work(query string) int64 {

	ctx := context.Background()

	// queue
	q, err := newQueue(ctx, project, topic)
	if err != nil {
		logger.Fatalf("Error creating pubsub client: %v", err)
	}

	// last ID
	var sinceID int64

	config := oauth1.NewConfig(consumerKey, consumerSecret)
	token := oauth1.NewToken(accessToken, accessSecret)

	httpClient := config.Client(oauth1.NoContext, token)
	tc := twitter.NewClient(httpClient)

	logger.Printf("Starting search for %s since ID: %d\n", query, sinceID)
	search, resp, err := tc.Search.Tweets(&twitter.SearchTweetParams{
		Query:           query,
		Lang:            "en",
		Count:           100,
		SinceID:         sinceID,
		IncludeEntities: twitter.Bool(true),
	})

	if err != nil {
		logger.Fatalf("Error executing search %s - %v", resp.Status, err)
	}

	logger.Printf("Processing tweets: %d\n", len(search.Statuses))
	for _, t := range search.Statuses {

		data, err := json.Marshal(t)
		if err != nil {
			logger.Printf("Error while marshaling tweet: %v", err)
			continue
		}

		// tweets come in newest first order so just make sure we capture the highest number
		// and start from there the next time
		if t.ID > sinceID {
			// publish
			if err = q.push(ctx, data); err != nil {
				logger.Printf("Error while publishing tweet: %v", err)
				continue
			}
			sinceID = t.ID
		}
	}

	logger.Printf("Done. Last ID: %d\n", sinceID)
	return sinceID
}
