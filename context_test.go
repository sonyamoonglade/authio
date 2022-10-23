package authio

import (
	"context"
	"testing"

	"github.com/sonyamoonglade/authio/session"
	"github.com/stretchr/testify/require"
)

func TestCanCreateAuthSessionFromCtx(t *testing.T) {

	initialSession := session.New(session.FromInt64(200))
	ctx := context.WithValue(context.Background(), session.CtxKey, initialSession)

	authSession, ok := getFromCtx(ctx)
	require.True(t, ok)
	require.Equal(t, initialSession.ID, authSession.ID)
	require.Equal(t, initialSession.Raw(), authSession.Raw())
}

func TestCanCastFromCtx(t *testing.T) {

	var userID int64 = 200

	sv := session.FromInt64(userID)
	authSession := session.New(sv)

	ctx := context.WithValue(context.Background(), session.CtxKey, authSession)

	//Specify int64 type as userID
	userIDFromSession, ok := ValueFromContext[int64](ctx)
	require.True(t, ok)
	require.Equal(t, userID, userIDFromSession)
}
