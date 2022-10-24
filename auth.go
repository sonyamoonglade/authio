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
}

//todo: check nils
func newAuth(logger Logger,
	store store.Store,
	settings map[string]*cookies.Setting) *Auth {

	if store == nil {
		panic("store is nil")
	}

	au := &Auth{
		logger:   logger,
		store:    store,
		settings: settings,
	}

	if logger == nil {
		au.logger = NewDefaultLogger(ErrorLevel)
	}

	au.settings[cookies.DefaultLabel] = cookies.DefaultSetting

	return au
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

func (au *Auth) HTTPGetSessionWithLabel(h http.HandlerFunc, settingLabel string) http.HandlerFunc {

	setting, ok := au.settings[settingLabel]
	if !ok {
		panic(fmt.Sprintf("setting at %s does not exist\n", settingLabel))
	}
	return func(w http.ResponseWriter, r *http.Request) {

		unsignedID, err := cookies.Get(r, setting.Name, setting.Secret)
		if err != nil {
			//authErrors.WithError(err)
			http.Error(w, "test", http.StatusUnauthorized)
			return
		}

		sessionValue, err := au.store.Get(unsignedID)
		if err != nil {
			http.Error(w, "test", http.StatusForbidden)
			return
		}
		//Can set explicit true because if it was not valid then
		//cookies.GetSigned would've returned an error
		session := session.NewFromCookie(unsignedID, sessionValue)

		h.ServeHTTP(w, r.WithContext(au.enrichContext(r.Context(), session)))
		return
	}
}

func (au *Auth) enrichContext(requestCtx context.Context, s *session.AuthSession) context.Context {
	return context.WithValue(requestCtx, session.CtxKey, s)
}
