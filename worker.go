package main

import (
	"context"
	"encoding/json"

	"github.com/mchmarny/gcputil/metric"

	"github.com/dghubble/go-twitter/twitter"
	"github.com/dghubble/oauth1"
)

func work(query string) int64 {

	ctx := context.Background()

	q, err := newQueue(ctx)
	if err != nil {
		logger.Fatalf("Error creating pubsub client: %v", err)
	}

	initStore(ctx)

	savedState, err := getState(ctx, query)
	if err != nil {
		logger.Fatalf("Error getting state for %s: %v", query, err)
	}

	config := oauth1.NewConfig(consumerKey, consumerSecret)
	token := oauth1.NewToken(accessToken, accessSecret)

	httpClient := config.Client(oauth1.NoContext, token)
	tc := twitter.NewClient(httpClient)

	logger.Printf("Starting search for %s since ID: %d\n", query, savedState.LastID)
	search, resp, err := tc.Search.Tweets(&twitter.SearchTweetParams{
		Query:           query,
		Lang:            "en",
		Count:           100,
		SinceID:         savedState.LastID,
		IncludeEntities: twitter.Bool(true),
	})

	if err != nil {
		logger.Fatalf("Error executing search %s - %v", resp.Status, err)
	}

	foundTweets := len(search.Statuses)
	errorTweets := 0
	logger.Printf("Processing tweets: %d\n", foundTweets)
	for _, t := range search.Statuses {

		data, err := json.Marshal(t)
		if err != nil {
			errorTweets++
			logger.Printf("Error while marshaling tweet: %v", err)
			continue
		}

		// tweets come in newest first order so just make sure we capture the highest number
		// and start from there the next time
		if t.ID > savedState.LastID {
			if err = q.push(ctx, data); err != nil {
				errorTweets++
				logger.Printf("Error while publishing tweet: %v", err)
				continue
			}
			savedState.LastID = t.ID
		}
	}

	saveState(ctx, savedState)
	publishMetrics(ctx, foundTweets, errorTweets)

	return savedState.LastID
}

func publishMetrics(ctx context.Context, total, errs int) {

	c, err := metric.NewClient(ctx)
	if err != nil {
		logger.Fatalf("Error creating metric client: %v", err)
	}

	if err := c.Publish(ctx, "total", "search-tweets", total); err != nil {
		logger.Printf("Error logging metrics: %v", err)
	}

	if err := c.Publish(ctx, "errors", "search-tweets", errs); err != nil {
		logger.Printf("Error logging metrics: %v", err)
	}
}
