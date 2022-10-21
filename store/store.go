package store

import (
	"errors"
	"strconv"
	"time"

	"github.com/sonyamoonglade/authio/session"
)

var (
	ErrEntryDoesNotExist = errors.New("entry does not exist")
)

type OverflowStrategy int

const (
	LRU OverflowStrategy = iota
	LFU
	RANDOM
)

type Store interface {
	Save(ID string, v session.SessionValue) error
	Delete(ID string) error
	Get(ID string) (session.SessionValue, error)
}

type Config struct {
	EntryTTL         time.Duration
	OverflowStrategy OverflowStrategy
	ParseFunc        ParseFromStoreFunc
}

type ParseFromStoreFunc func(v string) (session.SessionValue, error)

func ToInt64(v string) (session.SessionValue, error) {
	parsed, err := strconv.ParseInt(v, 10, 64)
	if err != nil {
		return nil, err
	}

	return session.FromInt64(parsed), nil
}
