package store

import "github.com/sonyamoonglade/authio/session"

//In memory implementation of store.Store
type InMemoryStore struct {
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
		data:             make(map[string]string),
		maxItems:         inMemoryCfg.MaxItems,
		currItems:        0,
		parseFunc:        cfg.ParseFunc,
		overflowStrategy: cfg.OverflowStrategy,
	}
}

func (i *InMemoryStore) Save(ID string, v session.SessionValue) error {
	if i.currItems == i.maxItems {
		panic("LRU!!")
	}

	i.data[ID] = v.String()
	i.currItems += 1
	return nil
}

func (i *InMemoryStore) Delete(ID string) error {
	_, ok := i.data[ID]
	if !ok {
		return ErrNoEntry
	}

	delete(i.data, ID)

	return nil
}

func (i *InMemoryStore) Get(ID string) (session.SessionValue, error) {
	stringValue, ok := i.data[ID]
	if !ok {
		return nil, ErrNoEntry
	}

	return i.parseFunc(stringValue)
}
