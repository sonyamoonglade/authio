package store

import (
	"sync"

	"github.com/sonyamoonglade/authio/session"
)

//In memory implementation of store.Store
type InMemoryStore struct {
	mu               *sync.RWMutex
	data             map[string]string
	maxItems         int64
	currItems        int64
	overflowStrategy OverflowStrategy
	parseFunc        ParseFromStoreFunc
}

type InMemoryConfig struct {
	MaxItems int64
}

func NewInMemoryStore(cfg *Config, inMemoryCfg *InMemoryConfig) *InMemoryStore {
	return &InMemoryStore{
		mu:               new(sync.RWMutex),
		data:             make(map[string]string),
		maxItems:         inMemoryCfg.MaxItems,
		currItems:        0,
		parseFunc:        cfg.ParseFunc,
		overflowStrategy: cfg.OverflowStrategy,
	}
}

func (i *InMemoryStore) Save(au *session.AuthSession) error {
	if i.currItems == i.maxItems {
		panic("LRU!!")
	}

	i.mu.Lock()
	i.data[au.ID] = au.Value.String()
	i.mu.Unlock()

	i.currItems += 1
	return nil
}

func (i *InMemoryStore) Delete(ID string) error {
	i.mu.RLock()
	_, ok := i.data[ID]
	if !ok {
		return ErrNoEntry
	}
	i.mu.RUnlock()

	i.mu.Lock()
	delete(i.data, ID)
	i.mu.Unlock()

	return nil
}

func (i *InMemoryStore) Get(ID string) (session.SessionValue, error) {
	i.mu.RLock()
	stringValue, ok := i.data[ID]
	i.mu.RUnlock()

	if !ok {
		return nil, ErrNoEntry
	}

	return i.parseFunc(stringValue)
}
