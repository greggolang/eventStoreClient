package eventeur

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

type TestType struct {
	InternalData string
	ExternalID   string
}

func TestCustomStore(t *testing.T) {
	item := TestType{
		InternalData: "Test",
		ExternalID:   "1234",
	}

	store := NewStore[TestType]()

	// Trying to update store with an index the doesn't exist.
	// The Update function should return `false`.
	assert.Equal(t, false, store.Update(1, &item))

	// Inserting item into the store and Getting it.
	id := store.Put(&item)
	assert.Equal(t, uint64(1), id)

	itemFromStore, ok := store.Get(id)
	assert.Equal(t, &item, itemFromStore)
	assert.True(t, ok)
}
