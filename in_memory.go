package authio

import (
	"sync"
)

//In memory implementation of store.Store
type InMemoryStore struct {
	mu               *sync.RWMutex
	data             map[string]string
	reversedData     map[string]string
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
		reversedData:     make(map[string]string),
		maxItems:         inMemoryCfg.MaxItems,
		currItems:        0,
		parseFunc:        cfg.ParseFunc,
		overflowStrategy: cfg.OverflowStrategy,
	}
}

func (i *InMemoryStore) Save(au *AuthSession) error {
	if i.currItems == i.maxItems {
		panic("LRU!!")
	}

	v := au.Value.String()

	i.mu.Lock()
	i.data[au.ID] = v
	i.reversedData[v] = au.ID
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
	delete(i.reversedData, i.data[ID])
	delete(i.data, ID)
	i.mu.Unlock()

	return nil
}

func (i *InMemoryStore) Get(ID string) (SessionValue, error) {
	i.mu.RLock()
	stringValue, ok := i.data[ID]
	i.mu.RUnlock()

	if !ok {
		return nil, ErrNoEntry
	}

	return i.parseFunc(stringValue)
}

//todo: test
func (i *InMemoryStore) GetSessionIDByValue(v SessionValue) (string, error) {
	i.mu.RLock()
	sessionID, ok := i.reversedData[v.String()]
	i.mu.RUnlock()

	if !ok {
		return "", ErrNoEntry
	}

	return sessionID, nil
}
