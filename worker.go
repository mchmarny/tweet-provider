package main

import (
	"context"
	"encoding/json"

	"github.com/dghubble/go-twitter/twitter"
	"github.com/dghubble/oauth1"
)

func work(query string) int64 {

	ctx := context.Background()

	// queue
	q, err := newQueue(ctx, projectID, topic)
	if err != nil {
		logger.Fatalf("Error creating pubsub client: %v", err)
	}

	initStore(ctx)

	savedState, err := getState(ctx, query)
	{
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

	logger.Printf("Processing tweets: %d\n", len(search.Statuses))
	for _, t := range search.Statuses {

		data, err := json.Marshal(t)
		if err != nil {
			logger.Printf("Error while marshaling tweet: %v", err)
			continue
		}

		// tweets come in newest first order so just make sure we capture the highest number
		// and start from there the next time
		if t.ID > savedState.LastID {
			// publish
			if err = q.push(ctx, data); err != nil {
				logger.Printf("Error while publishing tweet: %v", err)
				continue
			}
			savedState.LastID = t.ID
		}
	}

	saveState(ctx, savedState)

	logger.Printf("Done. Last ID: %d\n", savedState.LastID)
	return savedState.LastID
}
