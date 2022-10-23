package cookies

import (
	"errors"
	"net/http"
	"time"

	"github.com/sonyamoonglade/authio/gcmcrypt"
)

var (
	DefaultLabel     = "default"
	DefaultName      = "SESSION_ID"
	DefaultExpiresAt = time.Hour * 1
)

var DefaultSetting = &Setting{
	HttpOnly: true,
	Secure:   false,
	SameSite: http.SameSiteNoneMode,
	Expires:  DefaultExpiresAt,
	Path:     "",
	Domain:   "",
	Name:     DefaultName,
	Label:    DefaultLabel,
}

var (
	ErrCookieTooLong = errors.New("can not set cookie with 4096 or more lenght")
)

type Setting struct {
	Label    string
	Name     string
	Path     string
	Domain   string
	Secret   [16]byte
	Signed   bool
	HttpOnly bool
	Secure   bool
	SameSite http.SameSite
	Expires  time.Duration
}

func Write(w http.ResponseWriter, setting *Setting, unsignedValue string) error {
	return writeSigned(w, setting, unsignedValue)
}

func Get(r *http.Request, name string, key [16]byte) (string, error) {
	return getSigned(r, name, key)
}

func write(w http.ResponseWriter, setting *Setting, cookieValue string) error {

	if len(cookieValue) > 2<<11 { //4096
		return ErrCookieTooLong
	}

	http.SetCookie(w, &http.Cookie{
		Name:     setting.Name,
		Value:    cookieValue,
		Path:     setting.Path,
		Domain:   setting.Domain,
		Expires:  time.Now().Add(setting.Expires),
		Secure:   setting.Secure,
		HttpOnly: setting.HttpOnly,
		SameSite: setting.SameSite,
	})

	return nil
}

func writeSigned(w http.ResponseWriter, setting *Setting, unsignedValue string) error {

	signedValue, err := gcmcrypt.Encrypt(setting.Secret, unsignedValue)
	if err != nil {
		return err
	}

	write(w, setting, signedValue)

	return nil
}

func get(r *http.Request, name string, signed bool) (string, error) {
	cookie, err := r.Cookie(name)
	if err != nil {
		return "", err
	}

	return cookie.Value, nil
}

func getSigned(r *http.Request, name string, key [16]byte) (string, error) {
	//Explicitly set signed=true(3rd arg). See method name
	signedValue, err := get(r, name, true)
	if err != nil {
		return "", err
	}

	unsignedValue, err := gcmcrypt.Decrypt(key, signedValue)
	if err != nil {
		return "", err
	}

	return unsignedValue, nil
}
