package session

import (
	"github.com/google/uuid"
)

type AuthSession struct {
	ID       string
	SignedID string
	Value    SessionValue
}

func NewAuthSession(v SessionValue) *AuthSession {
	UUID := uuid.NewString()
	return &AuthSession{
		ID:       UUID,
		SignedID: UUID,
		Value:    v,
	}
}
