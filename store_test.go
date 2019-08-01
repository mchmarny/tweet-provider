package main

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestStateStore(t *testing.T) {

	if testing.Short() {
		t.Skip("Skipping TestStateStore")
	}

	ctx := context.Background()
	initStore(ctx)

	state := &storeState{
		LastID: int64(1),
	}

	// save
	err := saveState(ctx, state)
	assert.Nil(t, err)

	// get
	savedState, err := getState(ctx)
	assert.Nil(t, err)
	assert.Equalf(t, savedState, state,
		"Retreaved ID doesn't equal saved (%s != %s)",
		savedState, state)

}
