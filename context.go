package authio

import (
	"context"
)

func getFromCtx(ctx context.Context) (*AuthSession, bool) {
	v, ok := ctx.Value(CtxKey).(*AuthSession)
	return v, ok
}

func ValueFromContext[T SessionValueConstraint](ctx context.Context) (T, bool) {
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

func SessionFromContext(ctx context.Context) *AuthSession {
	authSession, _ := getFromCtx(ctx)
	return authSession
}
