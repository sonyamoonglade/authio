package store

import "github.com/sonyamoonglade/authio/session"

//In memory implementation of store.Store
type InMemoryStore struct {
	data             map[string]string
	currItems        int
	maxItems         int
	overflowStrategy OverflowStrategy
}

func NewInMemoryStore(cfg *Config, maxItems int) *InMemoryStore {
	return &InMemoryStore{
		data:             make(map[string]string),
		currItems:        0,
		overflowStrategy: cfg.OverflowStrategy,
		maxItems:         maxItems,
	}
}

func (i *InMemoryStore) Save(ID string, v session.SessionValue) error {
	//todo: setup limits
	i.data[ID] = v.String()
	i.currItems += 1
	return nil
}

func (i *InMemoryStore) Delete(ID string) error {
	_, ok := i.data[ID]
	if !ok {
		return ErrEntryDoesNotExist
	}

	delete(i.data, ID)

	return nil
}

func (i *InMemoryStore) Get(ID string) (session.SessionValue, error) {
	v, ok := i.data[ID]
	if !ok {
		return nil, ErrEntryDoesNotExist
	}
	_ = v
	panic("implement")

}

func (i *InMemoryStore) UseConfig(config *Config) {
	panic("not implemented") // TODO: Implement
}
