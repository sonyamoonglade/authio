package cookies

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/google/uuid"
	"github.com/sonyamoonglade/authio/gcmcrypt"
	"github.com/sonyamoonglade/authio/hash"
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
			require.True(t, hash.Compare(c.Value, mockID)) // see cookies.write

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
		Name:     DefaultName,
		Path:     "",
		Domain:   "",
		Expires:  DefaultExpiresAt,
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

	//With default setting set signed=false (see cookies.DefaultSetting)
	cookieValue, err := get(req, DefaultSetting.Name, false)
	require.NoError(t, err)

	require.True(t, hash.Compare(cookieValue, mockID))
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
		Expires:  DefaultExpiresAt,
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

func newRequestWithDefaultCookie(cookieValue string) *http.Request {
	req := httptest.NewRequest(http.MethodGet, "http://cool-url.com", nil)
	req.AddCookie(&http.Cookie{
		Name:     DefaultSetting.Name,
		Value:    hash.SHA1(cookieValue), //important!! With DefaultSetting cookie value is not signed but hashed
		Path:     "",
		Domain:   "",
		Expires:  DefaultExpiresAt,
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
		Expires:  setting.Expires,
		Secure:   setting.Secure,
		HttpOnly: setting.HttpOnly,
		SameSite: setting.SameSite,
	})

	return req
}
