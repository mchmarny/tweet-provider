package main

import (
	"encoding/json"
	"time"

	"github.com/dghubble/go-twitter/twitter"
	"github.com/dghubble/oauth1"
)

// Provider provides content
type Provider interface {
	Provide(query string, ch chan<- *[]byte)
}

// BatchProvider provides
type BatchProvider struct {
	Client    *twitter.Client
	SinceID   int64
	Frequency int
}

// NewBatchProvider returns BatchProvider instance pointer
func getProvider(frequency int, consumerKey, consumerSecret, accessToken, accessSecret string) (p Provider) {
	config := oauth1.NewConfig(consumerKey, consumerSecret)
	token := oauth1.NewToken(accessToken, accessSecret)
	httpClient := config.Client(oauth1.NoContext, token)
	client := twitter.NewClient(httpClient)
	provider := &BatchProvider{
		Client:    client,
		SinceID:   0,
		Frequency: frequency,
	}

	return provider
}

// Provide initiates the Tweeter stream subscription and pumps all messages into
// the passed in channel
func (p *BatchProvider) Provide(query string, ch chan<- *[]byte) {

	logger.Printf("Starting providing for: %s every %d sec\n",
		query, p.Frequency)

	p.getTweets(query, ch)

	if p.Frequency > 0 {
		tick := time.Tick(time.Duration(p.Frequency) * time.Second)
		for {
			select {
			case <-tick:
				p.getTweets(query, ch)
			}
		}
	}
}

// getTweets initiates the Tweeter stream subscription and pumps all messages into
// the passed in channel
func (p *BatchProvider) getTweets(query string, ch chan<- *[]byte) {

	logger.Printf("Quering for tweets since ID:%d\n", p.SinceID)

	search, resp, err := p.Client.Search.Tweets(&twitter.SearchTweetParams{
		Query:           query,
		Lang:            "en",
		Count:           100,
		SinceID:         p.SinceID,
		IncludeEntities: twitter.Bool(true),
	})

	if err != nil {
		logger.Fatalf("Error executing search %s - %v", resp.Status, err)
	}

	for _, t := range search.Statuses {

		debug.Printf("Id:%s On:%s By:@%s", t.IDStr, t.CreatedAt, t.User.ScreenName)

		msg, err := json.Marshal(t)
		if err != nil {
			logger.Printf("Error while marshaling object: %v", err)
			return
		}

		// tweets come in newest first order so just make sure we capture the highest number
		// and start from there the next time
		if t.ID > p.SinceID {
			p.SinceID = t.ID
		}

		ch <- &msg

	}

}
