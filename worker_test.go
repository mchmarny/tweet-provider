package main

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestPublishMetrics(t *testing.T) {

	if testing.Short() {
		t.Skip("Skipping TestStateStore")
	}

	ctx := context.Background()
	err := publishMetrics(ctx, 10, 0)
	assert.Nil(t, err)

}
