package authio

import (
	"errors"
	"strconv"
	"time"
)

var (
	ErrNoEntry = errors.New("entry does not exist")
)

type OverflowStrategy int

const (
	LRU OverflowStrategy = iota
)

type Store interface {
	Save(session *AuthSession) error
	Delete(ID string) error
	Get(ID string) (SessionValue, error)
	GetSessionIDByValue(v SessionValue) (string, error)
}

type Config struct {
	EntryTTL         time.Duration
	OverflowStrategy OverflowStrategy
	ParseFunc        ParseFromStoreFunc
}

type ParseFromStoreFunc func(v string) (SessionValue, error)

func ToInt64(v string) (SessionValue, error) {
	parsed, err := strconv.ParseInt(v, 10, 64)
	if err != nil {
		return nil, err
	}

	return NewValueFromInt64(parsed), nil
}

func ToString(v string) (SessionValue, error) {
	return NewValueFromString(v), nil
}
