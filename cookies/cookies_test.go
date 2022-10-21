package cookies

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/google/uuid"
	"github.com/sonyamoonglade/authio/session"
	"github.com/stretchr/testify/require"
)

func TestShouldWriteCookie(t *testing.T) {

	setting := DefaultSetting

	randomID := int64(527) //e.g. userID
	sv := session.FromInt64(randomID)

	session := session.NewAuthSession(sv)

	w := httptest.NewRecorder()
	write(w, setting, session)

	res := w.Result()
	cookies := res.Cookies()

	var sessionCookieExists bool
	for _, c := range cookies {
		if c.Name == DefaultSetting.Name {
			sessionCookieExists = true
			require.Equal(t, DefaultSetting.HttpOnly, c.HttpOnly)
			require.Equal(t, DefaultSetting.Secure, c.Secure)
			require.Equal(t, DefaultSetting.SameSite, c.SameSite)
			require.Equal(t, session.SignedID, c.Value) // see cookies.write
		}
	}

	require.True(t, sessionCookieExists)
}

func TestShouldWriteSignedCookie(t *testing.T) {

	setting := &Setting{
		Label:    "my-label",
		Signed:   true,
		Secret:   "my-ultra-mega-cool-secret-gcc",
		Name:     DefaultName,
		Path:     "",
		Domain:   "",
		Expires:  DefaultExpiresAt,
		Secure:   false,
		HttpOnly: true,
		SameSite: DefaultSetting.SameSite,
	}

	randomID := int64(527) //e.g. userID
	sv := session.FromInt64(randomID)

	session := session.NewAuthSession(sv)

	w := httptest.NewRecorder()
	err := writeSigned(w, setting, session)
	require.NoError(t, err)

	res := w.Result()
	cookies := res.Cookies()

	for _, c := range cookies {
		if c.Name == DefaultSetting.Name {
			require.NotEqual(t, session.ID, c.Value) // not equal (see cookies.write impl)
			require.Equal(t, session.SignedID, c.Value)
		}
	}

}

func TestShouldWriteAndGetCookie(t *testing.T) {

	mockValue := uuid.NewString()
	req := httptest.NewRequest(http.MethodGet, "http://cool-url.com", nil)
	req.AddCookie(&http.Cookie{
		Name:     DefaultSetting.Name,
		Value:    mockValue,
		Path:     "",
		Domain:   "",
		Expires:  DefaultExpiresAt,
		Secure:   false,
		HttpOnly: true,
		SameSite: DefaultSetting.SameSite,
	})

	cookieValue, err := get(req, DefaultSetting.Name)
	require.NoError(t, err)
	require.Equal(t, cookieValue, mockValue)
}

func TestShouldWriteAndNotGetCookie(t *testing.T) {

	mockValue := uuid.NewString()
	req := httptest.NewRequest(http.MethodGet, "http://cool-url.com", nil)
	req.AddCookie(&http.Cookie{
		Name:     DefaultSetting.Name,
		Value:    mockValue,
		Path:     "",
		Domain:   "",
		Expires:  DefaultExpiresAt,
		Secure:   false,
		HttpOnly: true,
		SameSite: DefaultSetting.SameSite,
	})

	cookieValue, err := get(req, "some-invalid-key")
	require.Error(t, err)
	require.Equal(t, "", cookieValue)
	require.Equal(t, http.ErrNoCookie, err)
}
