package store

import "github.com/sonyamoonglade/authio/internal/session"

//In memory implementation of store.Store

type InMemoryStore struct {
	data             map[string]string
	maxItems         int
	currItems        int
	overflowStrategy OverflowStrategy
}

func (i *InMemoryStore) Save(ID string, v session.SessionValue) error {
	if i.currItems >= i.maxItems {
		panic("implement LRU")
	}

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

}

func (i *InMemoryStore) UseConfig(config *Config) {
	panic("not implemented") // TODO: Implement
}
