package authio

import (
	"context"
	"errors"
	"fmt"
	"net/http"

	"github.com/sonyamoonglade/authio/cookies"
	"github.com/sonyamoonglade/authio/session"
	"github.com/sonyamoonglade/authio/store"
)

var (
	ErrNoSetting = errors.New("setting does not exist")
)

type Auth struct {
	logger   Logger
	store    store.Store
	settings map[string]*cookies.Setting
	Config   AuthioConfig
}

func (au *Auth) SaveSession(w http.ResponseWriter, settingLabel string, sv session.SessionValue) error {

	setting, ok := au.settings[settingLabel]
	if !ok {
		return ErrNoSetting
	}

	authSession := session.New(sv)
	err := au.store.Save(authSession)
	if err != nil {
		return fmt.Errorf("could not save to store: %w", err)
	}

	err = cookies.Write(w, setting, authSession.ID)
	if err != nil {
		return fmt.Errorf("could not write a cookie: %w", err)
	}

	return nil
}

//todo: test
func (au *Auth) InvalidateSession(w http.ResponseWriter, settingLabel string, sessionID string) error {

	setting, ok := au.settings[settingLabel]
	if !ok {
		return ErrNoSetting
	}

	err := au.store.Delete(sessionID)
	if err != nil {
		return fmt.Errorf("coult not delete a cookie: %w", err)
	}

	cookies.Delete(w, setting)

	return nil
}

// probably move to sep. file
func (au *Auth) newRedirectAuthed(h http.HandlerFunc, settingLabel string) http.HandlerFunc {

	setting, ok := au.settings[settingLabel]
	if !ok {
		panic(fmt.Sprintf("setting at %s does not exist\n", settingLabel))
	}

	return func(w http.ResponseWriter, r *http.Request) {
		_, err := cookies.Get(r, setting.Name, setting.Secret)
		// Either cookie does not exist or sessionID is some random garbage.
		if err != nil {
			au.logger.Warnf("could not get cookies: %s\n", setting.Name)
			h.ServeHTTP(w, r)
			return
		}
		au.logger.Debugf("err: %v", err)

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

		unsignedID, err := cookies.Get(r, setting.Name, setting.Secret)
		if err != nil {
			//authErrors.WithError(err)
			//redirect
			http.Error(w, "test", http.StatusUnauthorized)
			return
		}

		//lookup ttl
		sessionValue, err := au.store.Get(unsignedID)
		if err != nil {
			http.Error(w, "test", http.StatusForbidden)
			//redirect
			return
		}
		session := session.NewFromCookie(unsignedID, sessionValue)

		h.ServeHTTP(w, r.WithContext(au.addToContext(r.Context(), session)))
		return
	}
}

func (au *Auth) addToContext(requestCtx context.Context, s *session.AuthSession) context.Context {
	return context.WithValue(requestCtx, session.CtxKey, s)
}

func newAuth(authioConfig *AuthioConfig,
	logger Logger,
	store store.Store,
	settings cookies.CookieSettings) *Auth {

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

	au.settings[cookies.DefaultLabel] = cookies.DefaultSetting

	return au
}
