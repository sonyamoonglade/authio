package store

import (
	"errors"
	"time"

	"github.com/sonyamoonglade/authio/internal/session"
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
	UseConfig(config *Config)
	UseParseFunc(pf session.ParseFromStoreFunc)
}

type Config struct {
	MaxItems         int
	EntryTTL         time.Duration
	OverflowStrategy OverflowStrategy
}
