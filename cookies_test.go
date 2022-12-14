package authio

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/sonyamoonglade/authio/internal/gcmcrypt"
	"github.com/stretchr/testify/require"
)

func TestShouldWriteCookie(t *testing.T) {

	setting := DefaultSetting

	mockID := uuid.NewString()

	w := httptest.NewRecorder()

	err := write(w, setting, mockID)
	require.NoError(t, err)

	res := w.Result()
	cookies := res.Cookies()

	var sessionCookieExists bool
	for _, c := range cookies {
		if c.Name == DefaultSetting.Name {

			sessionCookieExists = true

			//With default setting cookie is not signed but hashed,
			//Thus should compare two hashes
			require.Equal(t, mockID, c.Value) // see cookies.write

			require.Equal(t, DefaultSetting.HttpOnly, c.HttpOnly)
			require.Equal(t, DefaultSetting.Secure, c.Secure)
			require.Equal(t, DefaultSetting.SameSite, c.SameSite)

		}
	}

	require.True(t, sessionCookieExists)
}

func TestShouldNotWriteCookieErrToLong(t *testing.T) {

	setting := DefaultSetting

	var mockID string
	//prepare 4097 length cookie value (4096 is limit)
	for i := 0; i < 4097; i++ {
		mockID += "a"
	}

	w := httptest.NewRecorder()

	err := write(w, setting, mockID)
	require.Error(t, err)
	require.Equal(t, ErrCookieTooLong, err)

	res := w.Result()
	cookies := res.Cookies()

	var sessionCookieExists bool
	for _, c := range cookies {
		if c.Name == DefaultSetting.Name {
			sessionCookieExists = true
		}
	}

	require.False(t, sessionCookieExists)
}

func TestShouldWriteSignedCookie(t *testing.T) {

	setting := &Setting{
		Label:    "my-label",
		Signed:   true,
		Secret:   gcmcrypt.KeyFromString("my-ultra-mega-cool-secret-gcc"),
		Name:     DefaultSetting.Name,
		Path:     "",
		Domain:   "",
		Expires:  DefaultSetting.Expires,
		Secure:   false,
		HttpOnly: true,
		SameSite: DefaultSetting.SameSite,
	}

	mockID := uuid.NewString()

	w := httptest.NewRecorder()

	err := writeSigned(w, setting, mockID)
	require.NoError(t, err)

	res := w.Result()
	cookies := res.Cookies()

	for _, c := range cookies {
		if c.Name == DefaultSetting.Name {
			require.NotEqual(t, mockID, c.Value) // not equal (see cookies.writeSigned impl)
		}
	}
}

func TestShouldGetCookie(t *testing.T) {

	mockID := uuid.NewString()
	req := newRequestWithDefaultCookie(mockID)

	cookieValue, err := get(req, DefaultSetting.Name, false)
	require.NoError(t, err)
	require.Equal(t, mockID, cookieValue)
}

func TestShouldNotGetCookieByInvalidKey(t *testing.T) {

	mockID := uuid.NewString()
	req := newRequestWithDefaultCookie(mockID)

	//See newRequestWithDefaultCookie (cookie.Name)
	cookieValue, err := get(req, "some-invalid-key", false)
	require.Error(t, err)
	require.Equal(t, "", cookieValue)
	require.Equal(t, http.ErrNoCookie, err)
}

func TestShouldGetAndUnsignCookieValue(t *testing.T) {

	mockID := uuid.NewString()

	setting := &Setting{
		Label:    "my-custom-cookie-setting",
		Signed:   true,
		HttpOnly: true,
		Secret:   gcmcrypt.KeyFromString("ajadfkjsadfkasdkfjkasfjdskfj23984124789247"),
		Secure:   false,
		SameSite: DefaultSetting.SameSite,
		Expires:  DefaultSetting.Expires,
		Path:     "",
		Domain:   "",
		Name:     "SESSION_ID_SIGNED",
	}

	//emulate that cookie has been already set and signed with such key beforehand...
	encryptedID, err := gcmcrypt.Encrypt(setting.Secret, mockID)
	require.NoError(t, err)

	req := newRequestWithCustomCookie(encryptedID, setting)
	decryptedID, err := getSigned(req, setting.Name, setting.Secret)
	require.NoError(t, err)

	require.NotEqual(t, encryptedID, decryptedID)
	require.Equal(t, mockID, decryptedID)

}

func TestShouldDeleteCookie(t *testing.T) {

	mockID := uuid.NewString()
	w := httptest.NewRecorder()
	write(w, DefaultSetting, mockID)

	deleteCookie(w, DefaultSetting)

	res := w.Result()
	cookies := res.Cookies()

	var sessionCookieCount int
	for _, c := range cookies {
		if c.Name == DefaultSetting.Name {
			sessionCookieCount += 1
		}
		fmt.Printf("%v\n", c)
	}

	//1st - original session cookie
	//2nd - same cookie but with zero value and -1 MaxAge (browsers will remove 1st cookie)
	require.Equal(t, 2, sessionCookieCount)
}

func newRequestWithDefaultCookie(cookieValue string) *http.Request {
	req := httptest.NewRequest(http.MethodGet, "http://cool-url.com", nil)
	req.AddCookie(&http.Cookie{
		Name:     DefaultSetting.Name,
		Value:    cookieValue, //important!! With DefaultSetting cookie value is not signed but hashed
		Path:     "",
		Domain:   "",
		Expires:  time.Now().Add(DefaultSetting.Expires),
		Secure:   false,
		HttpOnly: true,
		SameSite: DefaultSetting.SameSite,
	})
	return req
}

func newRequestWithCustomCookie(cookieValue string, setting *Setting) *http.Request {
	req := httptest.NewRequest(http.MethodGet, "http://cool-url.com", nil)
	req.AddCookie(&http.Cookie{
		Name:     setting.Name,
		Value:    cookieValue,
		Path:     setting.Path,
		Domain:   setting.Domain,
		Expires:  time.Now().Add(DefaultSetting.Expires),
		Secure:   setting.Secure,
		HttpOnly: setting.HttpOnly,
		SameSite: setting.SameSite,
	})

	return req
}
