package main

import (
	"fmt"

	"cloud.google.com/go/pubsub"
	"golang.org/x/net/context"
)

// Publisher represents generic publisher interface
type Publisher interface {
	Publish(content *[]byte) (id string, err error)
}

// PubSubPublisher represents basic implementation of the publisher interface
type PubSubPublisher struct {
	Context *context.Context
	Topic   *pubsub.Topic
	Client  *pubsub.Client
}

// Publish publishes
func (p *PubSubPublisher) Publish(content *[]byte) (id string, err error) {

	msg := &pubsub.Message{Data: *content}
	result := p.Topic.Publish(*p.Context, msg)
	serverID, err := result.Get(*p.Context)
	if err != nil {
		return serverID,
			fmt.Errorf("Error while publishing message: %v:%s", err, serverID)
	}

	return serverID, nil
}

func getPublisher(ctx *context.Context, projectID, topicName string) (p Publisher, err error) {

	client, err := pubsub.NewClient(*ctx, projectID)
	if err != nil {
		logger.Fatalf("Failed to create client: %v", err)
		return nil, err
	}

	publisher := &PubSubPublisher{
		Context: ctx,
		Client:  client,
		Topic:   client.Topic(topicName),
	}

	return publisher, nil

}
