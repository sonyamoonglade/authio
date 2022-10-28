package authio

import (
	"net/http"
)

type AuthioMiddlewareFactory func(h http.HandlerFunc) http.HandlerFunc

func (au *Auth) AuthRequired(settingLabel string) AuthioMiddlewareFactory {
	return func(h http.HandlerFunc) http.HandlerFunc {
		return au.newAuthRequiredWithSetting(h, settingLabel)
	}
}

//TODO: accept array of labels and check for every setting
func (au *Auth) RedirectAuthed(settingLabel string) AuthioMiddlewareFactory {
	return func(h http.HandlerFunc) http.HandlerFunc {
		return au.newRedirectAuthed(h, settingLabel)
	}
}
