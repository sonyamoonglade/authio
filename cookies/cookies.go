package cookies

import (
	"net/http"
	"time"

	"github.com/sonyamoonglade/authio/gcmcrypt"
	"github.com/sonyamoonglade/authio/session"
)

var (
	DefaultLabel     = "default"
	DefaultName      = "SESSION_ID"
	DefaultExpiresAt = time.Now().Add(time.Hour * 1)
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

type Setting struct {
	Label    string
	Signed   bool
	Secret   string
	HttpOnly bool
	Secure   bool
	SameSite http.SameSite
	Expires  time.Time
	Path     string
	Domain   string
	Name     string
}

func write(w http.ResponseWriter, setting *Setting, session *session.AuthSession) {
	http.SetCookie(w, &http.Cookie{
		Name:     setting.Name,
		Value:    session.SignedID, //pay attention that SignedID is set to http cookie not ID!
		Path:     setting.Path,
		Domain:   setting.Domain,
		Expires:  setting.Expires,
		Secure:   setting.Secure,
		HttpOnly: setting.HttpOnly,
		SameSite: setting.SameSite,
	})
}

func writeSigned(w http.ResponseWriter, setting *Setting, session *session.AuthSession) error {

	//it should encrypt the value
	key := gcmcrypt.KeyFromString(setting.Secret)

	signedID, err := gcmcrypt.Encrypt(key, session.ID)
	if err != nil {
		return err
	}

	session.SignedID = signedID

	write(w, setting, session)

	return nil
}

func get(r *http.Request, name string) (string, error) {
	cookie, err := r.Cookie(name)
	if err != nil {
		return "", err
	}

	return cookie.Value, nil
}
