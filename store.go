package main

import (
	"context"
	"errors"
	"fmt"
	"hash/fnv"
	"log"

	"cloud.google.com/go/firestore"
	en "github.com/mchmarny/gcputil/env"
)

const (
	eventCollectionDefaultName = "twitter-query-state"
	recordIDPrefix             = "id-"
)

var (
	fsClient            *firestore.Client
	stateColl           *firestore.CollectionRef
	errNilDocRef        = errors.New("firestore: nil DocumentRef")
	eventCollectionName = en.MustGetEnvVar("COLL_NAME", eventCollectionDefaultName)
)

type storeState struct {
	ID     string `json:"id" firestore:"id"`
	LastID int64  `json:"last_id" firestore:"last_id"`
	Query  string `json:"query" firestore:"query"`
}

func initStore(ctx context.Context) {

	c, err := firestore.NewClient(ctx, projectID)
	if err != nil {
		log.Fatalf("Error while creating Firestore client: %v", err)
	}
	fsClient = c
	stateColl = c.Collection(eventCollectionName)
}

func getState(ctx context.Context, query string) (state *storeState, err error) {

	if query == "" {
		return nil, errors.New("Nil query")
	}

	id := getQueryID(query)
	d, err := stateColl.Doc(id).Get(ctx)
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
	if state == nil {
		return errors.New("Nil state")
	}
	state.ID = getQueryID(state.Query)
	_, err := stateColl.Doc(state.ID).Set(ctx, state)
	if err != nil {
		return fmt.Errorf("Error on save: %v", err)
	}
	return nil
}

func getQueryID(query string) string {
	h := fnv.New32a()
	h.Write([]byte(query))
	return fmt.Sprintf("%s%d", recordIDPrefix, h.Sum32())
}
