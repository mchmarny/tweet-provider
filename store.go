package main

import (
	"context"
	"errors"
	"fmt"
	"log"

	"cloud.google.com/go/firestore"
	en "github.com/mchmarny/gcputil/env"
)

const (
	eventCollectionDefaultName = "twitter-query-state"
	stateCollID                = "twitter-query-state-id"
)

var (
	fsClient            *firestore.Client
	stateColl           *firestore.CollectionRef
	errNilDocRef        = errors.New("firestore: nil DocumentRef")
	eventCollectionName = en.MustGetEnvVar("COLL_NAME", eventCollectionDefaultName)
)

type storeState struct {
	LastID int64 `json:"last_id" firestore:"last_id"`
}

func initStore(ctx context.Context) {

	c, err := firestore.NewClient(ctx, projectID)
	if err != nil {
		log.Fatalf("Error while creating Firestore client: %v", err)
	}
	fsClient = c
	stateColl = c.Collection(eventCollectionName)
}

func getState(ctx context.Context) (state *storeState, err error) {

	d, err := stateColl.Doc(stateCollID).Get(ctx)
	if err != nil {
		return nil, err
	}

	var s storeState
	if err := d.DataTo(&s); err != nil {
		return nil, fmt.Errorf("Stored data not user: %v", err)
	}

	return &s, nil
}

func saveState(ctx context.Context, state *storeState) error {
	_, err := stateColl.Doc(stateCollID).Set(ctx, state)
	if err != nil {
		return fmt.Errorf("Error on save: %v", err)
	}
	return nil
}
