package authio

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/sonyamoonglade/authio/cookies"
	"github.com/sonyamoonglade/authio/gcmcrypt"
	"github.com/sonyamoonglade/authio/hash"
	"github.com/sonyamoonglade/authio/session"
	"github.com/sonyamoonglade/authio/store"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
)

func TestCanBuildAuth(t *testing.T) {

	lg, _ := zap.NewProduction()
	logger := lg.Sugar()

	auth := NewAuthBuilder().
		AddCookieSetting(cookies.DefaultSetting).
		UseLogger(logger).
		UseStore(store.NewInMemoryStore(&store.Config{
			EntryTTL:         time.Hour * 24,
			OverflowStrategy: store.LRU,
			ParseFunc:        store.ToInt64,
		}, &store.InMemoryConfig{
			MaxItems: int64(100),
		})).
		Build()

	require.NotNil(t, auth)
	require.NotNil(t, auth.settings)
	require.NotNil(t, auth.logger)
	require.NotNil(t, auth.store)
	require.NotNil(t, auth.settings[cookies.DefaultLabel])

	require.Nil(t, auth.settings["jo-mama"])
}

func TestCanBuildWithDefaultLogger(t *testing.T) {

	auth := NewAuthBuilder().
		UseLogger(nil). // <---
		UseStore(store.NewInMemoryStore(&store.Config{
			EntryTTL:         time.Hour * 24,
			OverflowStrategy: store.LRU,
			ParseFunc:        store.ToInt64,
		}, &store.InMemoryConfig{
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

	auth := newDefaultAuth(store.ToInt64)

	var mockUserID int64 = 542 //random user_id

	cookieSetting := newSignedSessionSetting()
	auth.settings["signed-cookie-label"] = cookieSetting //save cookie setting config by label

	mockID := uuid.NewString()

	//Now encrypt mockID as if it was done by register/login endpoint and set to cookie
	encryptedID, err := gcmcrypt.Encrypt(cookieSetting.Secret, mockID)
	require.NoError(t, err)

	//Create request with encrypted cookie
	r := newRequestWithCustomCookie(encryptedID, cookieSetting)

	err = auth.store.Save(mockID, session.FromInt64(mockUserID)) //immitate that session has been already created beforehand...
	require.NoError(t, err)

	handler := func(w http.ResponseWriter, r *http.Request) {

		ctx := r.Context()

		//Specify int64 as mockUserID
		userID, ok := ValueFromContext[int64](ctx)
		require.True(t, ok)
		require.Equal(t, mockUserID, userID)

	}

	m := auth.HTTPGetSessionWithLabel(handler, "signed-cookie-label") // init middleware with settings by label

	w := httptest.NewRecorder()

	m.ServeHTTP(w, r)

}

func newDefaultAuth(pf store.ParseFromStoreFunc) *Auth {
	return NewAuthBuilder().
		UseLogger(nil).
		UseStore(store.NewInMemoryStore(&store.Config{
			EntryTTL:         time.Hour * 24,
			OverflowStrategy: store.LRU,
			ParseFunc:        pf,
		}, &store.InMemoryConfig{
			MaxItems: int64(100),
		})).
		Build()
}

func newRequestWithDefaultCookie(cookieValue string) *http.Request {
	req := httptest.NewRequest(http.MethodGet, "http://cool-url.com", nil)
	req.AddCookie(&http.Cookie{
		Name:     cookies.DefaultSetting.Name,
		Value:    hash.SHA1(cookieValue), //important!! With DefaultSetting cookie value is not signed but hashed
		Path:     "",
		Domain:   "",
		Expires:  time.Now().Add(cookies.DefaultExpiresAt),
		Secure:   false,
		HttpOnly: true,
		SameSite: cookies.DefaultSetting.SameSite,
	})
	return req
}
func newRequestWithCustomCookie(cookieValue string, setting *cookies.Setting) *http.Request {
	req := httptest.NewRequest(http.MethodGet, "http://cool-url.com", nil)
	req.AddCookie(&http.Cookie{
		Name:     setting.Name,
		Value:    cookieValue,
		Path:     setting.Path,
		Domain:   setting.Domain,
		Expires:  time.Now().Add(cookies.DefaultExpiresAt),
		Secure:   setting.Secure,
		HttpOnly: setting.HttpOnly,
		SameSite: setting.SameSite,
	})

	return req
}

func newSignedSessionSetting() *cookies.Setting {
	return &cookies.Setting{
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
