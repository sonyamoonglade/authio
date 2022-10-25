package authio

import (
	"context"

	"github.com/sonyamoonglade/authio/session"
)

func getFromCtx(ctx context.Context) (*session.AuthSession, bool) {
	v, ok := ctx.Value(session.CtxKey).(*session.AuthSession)
	return v, ok
}

func ValueFromContext[T session.SessionValueConstraint](ctx context.Context) (T, bool) {
	authSession, ok := getFromCtx(ctx)

	if ok == false {
		return *new(T), false
	}

	casted, ok := (authSession.Raw()).(T)
	if ok == false {
		return *new(T), false
	}

	return casted, true
}

func SessionFromContext(ctx context.Context) *session.AuthSession {
	authSession, _ := getFromCtx(ctx)
	return authSession
}
