package store

import (
	"sync"
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
	au := session.New(session.FromString("random-value"))

	err := store.Save(au)
	require.NoError(t, err)

	sv, err := store.Get(au.ID)
	require.NoError(t, err)

	require.Equal(t, "random-value", sv.String())

}

func TestDelete(t *testing.T) {

	store := newStore()
	au := session.New(session.FromString("random-value"))

	err := store.Save(au)
	require.NoError(t, err)

	err = store.Delete(au.ID)
	require.NoError(t, err)

	sv, err := store.Get(au.ID)
	require.Error(t, err)
	require.Equal(t, ErrNoEntry, err)

	require.Nil(t, sv)
}

func TestConcurrentRW(t *testing.T) {

	sessions := make([]*session.AuthSession, 100, 100)
	for i := 0; i < 100; i++ {
		sessions[i] = session.New(session.FromString(uuid.NewString()))
	}

	wg := new(sync.WaitGroup)
	store := newStore()

	wg.Add(3)
	//writer
	go func() {
		var err error
		for _, s := range sessions {
			err = store.Save(s)
			require.NoError(t, err)
		}

		defer wg.Done()
	}()

	//reader
	go func() {
		//readers
		for i := 0; i < 5; i++ {
			go func() {
				for _, s := range sessions {
					store.Get(s.ID)
				}
			}()
		}

		defer wg.Done()
	}()

	//deleter
	go func() {
		for _, s := range sessions {
			store.Delete(s.ID)
		}

		defer wg.Done()
	}()

	wg.Wait()

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
