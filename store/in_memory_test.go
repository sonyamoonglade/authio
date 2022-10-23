package store

import (
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/sonyamoonglade/authio/session"
	"github.com/stretchr/testify/require"
)

func TestShouldInitInMemoryStore(t *testing.T) {

	inMemory := newStore()
	require.NotNil(t, inMemory)
	require.Equal(t, LRU, inMemory.overflowStrategy)
	require.Equal(t, int64(100), inMemory.maxItems)
	require.Zero(t, inMemory.currItems)

}

func TestGet(t *testing.T) {

	store := newStore()

	ID := uuid.NewString()
	store.data[ID] = "random-value"

	sv, err := store.Get(ID)
	require.NoError(t, err)
	require.Equal(t, "random-value", sv.String())

}

func TestGetNoEntry(t *testing.T) {
	store := newStore()

	ID := uuid.NewString()

	sv, err := store.Get(ID)
	require.Error(t, err)
	require.Nil(t, sv)
	require.Equal(t, ErrNoEntry, err)
}

func TestSave(t *testing.T) {
	store := newStore()
	ID := uuid.NewString()

	err := store.Save(ID, session.FromString("random-value"))
	require.NoError(t, err)

	sv, err := store.Get(ID)
	require.NoError(t, err)

	require.Equal(t, "random-value", sv.String())

}

func TestDelete(t *testing.T) {

	store := newStore()
	ID := uuid.NewString()

	err := store.Save(ID, session.FromString("random-value"))
	require.NoError(t, err)

	err = store.Delete(ID)
	require.NoError(t, err)

	sv, err := store.Get(ID)
	require.Error(t, err)
	require.Equal(t, ErrNoEntry, err)

	require.Nil(t, sv)
}

func newStore() *InMemoryStore {
	return NewInMemoryStore(&Config{
		EntryTTL:         time.Hour * 1,
		OverflowStrategy: LRU,
		ParseFunc:        ToString,
	}, &InMemoryConfig{
		MaxItems: int64(100),
	})

}
