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
	Path     string //optional
	Domain   string //optional
	Secret   [16]byte
	Signed   bool //optional (for now not...)
	HttpOnly bool
	Secure   bool
	SameSite http.SameSite
	Expires  time.Duration
}

func Write(w http.ResponseWriter, setting *Setting, sessionID string) error {
	return writeSigned(w, setting, sessionID)
}

func Get(r *http.Request, name string, key [16]byte) (string, error) {
	return getSigned(r, name, key)
}

//sessionID is what's written in the cookie
func write(w http.ResponseWriter, setting *Setting, sessionID string) error {
	if len(sessionID) > 2<<11 { //4096
		return ErrCookieTooLong
	}

	c := &http.Cookie{
		Name:     setting.Name,
		Value:    sessionID,
		Expires:  time.Now().Add(setting.Expires),
		Secure:   setting.Secure,
		HttpOnly: setting.HttpOnly,
		SameSite: setting.SameSite,
	}

	if setting.Path != "" {
		c.Path = setting.Path
	}

	if setting.Domain != "" {
		c.Domain = setting.Domain
	}

	http.SetCookie(w, c)

	return nil
}

func writeSigned(w http.ResponseWriter, setting *Setting, sessionID string) error {

	signedID, err := gcmcrypt.Encrypt(setting.Secret, sessionID)
	if err != nil {
		return err
	}

	write(w, setting, signedID)

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
