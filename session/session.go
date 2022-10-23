package session

import (
	"github.com/google/uuid"
)

const (
	CtxKey = "session"
)

//todo: implement more
type SessionValueConstraint interface {
	int64 | string
}

type AuthSession struct {
	ID    string
	Value SessionValue
}

func NewFromCookie(ID string, v SessionValue) *AuthSession {
	return &AuthSession{
		ID:    ID,
		Value: v,
	}
}

func New(v SessionValue) *AuthSession {
	return &AuthSession{
		ID:    uuid.NewString(),
		Value: v,
	}
}

func (a *AuthSession) Raw() interface{} {
	return a.Value.Raw()
}
