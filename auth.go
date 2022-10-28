package authio

import (
	"context"
	"errors"
	"fmt"
	"net/http"
)

var (
	ErrNoSetting = errors.New("setting does not exist")
)

type Auth struct {
	logger   Logger
	store    Store
	settings map[string]*Setting
	Config   AuthioConfig
}

func (au *Auth) SaveSession(w http.ResponseWriter, settingLabel string, sv SessionValue) error {

	setting, ok := au.settings[settingLabel]
	if !ok {
		return ErrNoSetting
	}

	authSession := New(sv)
	err := au.store.Save(authSession)
	if err != nil {
		return fmt.Errorf("could not save to store: %w", err)
	}

	err = Write(w, setting, authSession.ID)
	if err != nil {
		return fmt.Errorf("could not write a cookie: %w", err)
	}

	return nil
}

//todo: test
func (au *Auth) InvalidateSessionByID(w http.ResponseWriter, settingLabel string, sessionID string) error {

	setting, ok := au.settings[settingLabel]
	if !ok {
		return ErrNoSetting
	}

	err := au.store.Delete(sessionID)
	if err != nil {
		return fmt.Errorf("coult not delete a cookie: %w", err)
	}

	Delete(w, setting)

	return nil
}

//todo: test
// InvalidateSessionByValue does not actually delete a cookie
// but invalidates it (deletes) in store, so no longer
// accesses would success and user asosiated with
// invalidated session will get no access (rejected by a middleware)
func (au *Auth) InvalidateSessionByValue(sv SessionValue) error {

	sessionID, err := au.store.GetSessionIDByValue(sv)
	if err != nil {
		return err
	}
	au.logger.Debugf("sessionID: %s", sessionID)

	err = au.store.Delete(sessionID)
	if err != nil {
		return err
	}

	return nil
}

// probably move to sep. file
func (au *Auth) newRedirectAuthed(h http.HandlerFunc, settingLabel string) http.HandlerFunc {

	setting, ok := au.settings[settingLabel]
	if !ok {
		panic(fmt.Sprintf("setting at %s does not exist\n", settingLabel))
	}

	return func(w http.ResponseWriter, r *http.Request) {
		_, err := Get(r, setting.Name, setting.Secret)
		// Either cookie does not exist or sessionID is some random garbage.
		if err != nil {
			au.logger.Warnf(err.Error())
			h.ServeHTTP(w, r)
			return
		}

		// Session has passed de-encryption so redirect user
		http.Redirect(w, r, au.Config.Paths.OnAuthNotRequired, http.StatusTemporaryRedirect)
		au.logger.Debugf("redirected to: %s", au.Config.Paths.OnAuthNotRequired)
		return
	}
}

func (au *Auth) newAuthRequiredWithSetting(h http.HandlerFunc, settingLabel string) http.HandlerFunc {

	setting, ok := au.settings[settingLabel]
	if !ok {
		panic(fmt.Sprintf("setting at %s does not exist\n", settingLabel))
	}

	return func(w http.ResponseWriter, r *http.Request) {

		unsignedID, err := Get(r, setting.Name, setting.Secret)
		if err != nil {
			//authErrors.WithError(err)
			//redirect
			if err == http.ErrNoCookie {
				http.Error(w, fmt.Sprintf("Missing a %s cookie", setting.Name), http.StatusUnauthorized)
				//redirect...
				return
			}
			http.Error(w, "Internal error", http.StatusInternalServerError)
			return
		}

		//todo: lookup ttl
		sessionValue, err := au.store.Get(unsignedID)
		if err != nil {
			http.Error(w, "Invalid auth token", http.StatusForbidden)
			//redirect
			return
		}
		session := NewFromCookie(unsignedID, sessionValue)

		h.ServeHTTP(w, r.WithContext(au.addToContext(r.Context(), session)))
		return
	}
}

func (au *Auth) addToContext(requestCtx context.Context, s *AuthSession) context.Context {
	return context.WithValue(requestCtx, CtxKey, s)
}

func newAuth(authioConfig *AuthioConfig,
	logger Logger,
	store Store,
	settings CookieSettings) *Auth {

	if store == nil {
		panic("store is nil")
	}

	if authioConfig == nil {
		authioConfig = &AuthioConfig{}
		authioConfig.Defaults()
	}

	au := &Auth{
		logger:   logger,
		store:    store,
		settings: settings,
		Config:   *authioConfig,
	}

	if logger == nil {
		au.logger = NewDefaultLogger(DebugLevel)
	}

	au.settings[DefaultLabel] = DefaultSetting

	return au
}
