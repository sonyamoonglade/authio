package session

import (
	"github.com/google/uuid"
)

type SessionValue interface {
	String() string
}

type ParseFromStoreFunc func(v string) SessionValue

var ParseFunc ParseFromStoreFunc

func SetGlobParseFunc(f ParseFromStoreFunc) {
	ParseFunc = f
}

type AuthSession struct {
	ID    string
	Value SessionValue
}

func NewAuthSession(v SessionValue) *AuthSession {
	return &AuthSession{
		ID:    uuid.New().String(),
		Value: v,
	}
}
