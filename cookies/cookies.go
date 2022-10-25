package cookies

import (
	"errors"
	"net/http"
	"time"

	"github.com/sonyamoonglade/authio/gcmcrypt"
)

var (
	DefaultLabel = "default"
)

var DefaultSetting = &Setting{
	HttpOnly: true,
	Secure:   false,
	SameSite: http.SameSiteNoneMode,
	Expires:  time.Hour * 1,
	Path:     "/",
	Name:     "SESSION_ID",
	Label:    DefaultLabel,
	Secret:   gcmcrypt.KeyFromString("zxZC1b2316ZXC!ZXC!B@#"),
	Signed:   true,
}

var (
	ErrCookieTooLong = errors.New("can not set cookie with 4096 or more length")
)

type CookieSettings map[string]*Setting

type Setting struct {
	Label    string
	Name     string
	Path     string //optional
	Domain   string //optional
	Secret   [16]byte
	Signed   bool //optional (not for now)
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

//The exact setting should be passed that cookie has been written with!
func Delete(w http.ResponseWriter, setting *Setting) {
	deleteCookie(w, setting)
}

//sessionID is what's written in the cookie
func write(w http.ResponseWriter, setting *Setting, sessionID string) error {
	if len(sessionID) > 2<<11 { //4096
		return ErrCookieTooLong
	}
	exp := time.Now().Add(setting.Expires)
	c := &http.Cookie{
		Name:     setting.Name,
		Value:    sessionID,
		Expires:  exp,
		MaxAge:   exp.Second(),
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

func deleteCookie(w http.ResponseWriter, setting *Setting) {
	c := &http.Cookie{
		Name:     setting.Name,
		Value:    "", //zero
		Secure:   setting.Secure,
		HttpOnly: setting.HttpOnly,
		SameSite: setting.SameSite,
		MaxAge:   -1,
	}

	if setting.Path != "" {
		c.Path = setting.Path
	}

	if setting.Domain != "" {
		c.Domain = setting.Domain
	}

	http.SetCookie(w, c)
}
