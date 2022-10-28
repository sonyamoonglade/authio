package authio

import (
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/sonyamoonglade/authio/internal/gcmcrypt"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
)

func TestCanBuildAuth(t *testing.T) {

	lg, _ := zap.NewProduction()
	logger := lg.Sugar()

	auth := NewAuthBuilder().
		AddCookieSetting(DefaultSetting).
		UseLogger(logger).
		UseStore(NewInMemoryStore(&Config{
			EntryTTL:         time.Hour * 24,
			OverflowStrategy: LRU,
			ParseFunc:        ToInt64,
		}, &InMemoryConfig{
			MaxItems: int64(100),
		})).
		Build()

	require.NotNil(t, auth)
	require.NotNil(t, auth.settings)
	require.NotNil(t, auth.logger)
	require.NotNil(t, auth.store)
	require.NotNil(t, auth.settings[DefaultLabel])

	require.Nil(t, auth.settings["jo-mama"])
}

func TestCanBuildWithDefaultLogger(t *testing.T) {

	auth := NewAuthBuilder().
		UseLogger(nil). // <---
		UseStore(NewInMemoryStore(&Config{
			EntryTTL:         time.Hour * 24,
			OverflowStrategy: LRU,
			ParseFunc:        ToInt64,
		}, &InMemoryConfig{
			MaxItems: int64(100),
		})).
		Build()

	require.NotNil(t, auth.logger)

}

func TestCantBuildWithoutStore(t *testing.T) {
	builder := NewAuthBuilder().
		UseLogger(nil) // <---

	require.Panics(t, func() {
		builder.Build()
	})

}

func TestMiddlewareWithLabelShouldFillContext(t *testing.T) {

	auth := newDefaultAuth(ToInt64)

	var mockUserID int64 = 542 //random user_id

	cookieSetting := newSignedSessionSetting()
	auth.settings["signed-cookie-label"] = cookieSetting //save cookie setting config by label

	authSession := New(NewValueFromInt64(mockUserID))
	//Now encrypt mockID as if it was done by register/login endpoint and set to cookie
	encryptedID, err := gcmcrypt.Encrypt(cookieSetting.Secret, authSession.ID)
	require.NoError(t, err)

	//Create request with encrypted cookie
	r := newRequestWithCustomCookie(encryptedID, cookieSetting)

	err = auth.store.Save(authSession) //immitate that session has been already created beforehand...
	require.NoError(t, err)

	handler := func(w http.ResponseWriter, r *http.Request) {
		//Specify int64 as mockUserID
		userID, ok := ValueFromContext[int64](r.Context())
		require.True(t, ok)
		require.Equal(t, mockUserID, userID)

		//use userID like in real handler...
		//...
		//...

		w.WriteHeader(http.StatusOK)
		w.Write([]byte("Hello-world!"))
		return
	}

	m := auth.newAuthRequiredWithSetting(handler, "signed-cookie-label") // init middleware with settings by label

	w := httptest.NewRecorder()

	m.ServeHTTP(w, r)

	b, err := io.ReadAll(w.Result().Body)
	require.NoError(t, err)

	require.Equal(t, w.Result().StatusCode, http.StatusOK)
	require.Equal(t, string(b), "Hello-world!")
}

func newDefaultAuth(pf ParseFromStoreFunc) *Auth {
	return NewAuthBuilder().
		UseLogger(nil).
		UseStore(NewInMemoryStore(&Config{
			EntryTTL:         time.Hour * 24,
			OverflowStrategy: LRU,
			ParseFunc:        pf,
		}, &InMemoryConfig{
			MaxItems: int64(100),
		})).
		Build()
}

func newSignedSessionSetting() *Setting {
	return &Setting{
		Label:    "signed-cookie-label",
		Name:     "signed-name",
		Path:     "",
		Domain:   "",
		Secret:   gcmcrypt.KeyFromString("asdfkjasfkddsfjk123z1234123"),
		Signed:   true,
		HttpOnly: true,
		Secure:   false,
		SameSite: http.SameSiteLaxMode,
		Expires:  time.Hour * 1,
	}
}
