package authio

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestCanCreateAuthSessionFromCtx(t *testing.T) {

	initialSession := New(NewValueFromInt64(200))
	ctx := context.WithValue(context.Background(), CtxKey, initialSession)

	authSession, ok := getFromCtx(ctx)
	require.True(t, ok)

	require.Equal(t, initialSession.ID, authSession.ID)
	require.Equal(t, initialSession.Raw(), authSession.Raw())
}

func TestCanCastFromCtx(t *testing.T) {

	var userID int64 = 200

	sv := NewValueFromInt64(userID)
	authSession := New(sv)

	ctx := context.WithValue(context.Background(), CtxKey, authSession)

	//Specify int64 type as userID
	userIDFromSession, ok := ValueFromContext[int64](ctx)
	require.True(t, ok)
	require.Equal(t, userID, userIDFromSession)
}
