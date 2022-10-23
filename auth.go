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
	ErrSettingExist = errors.New("setting already exists")
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

func (au *Auth) HTTPGetSessionWithLabel(h http.HandlerFunc, settingLabel string) http.HandlerFunc {

	setting, ok := au.settings[settingLabel]
	if !ok {
		panic(fmt.Sprintf("setting at %s does not exist\n", settingLabel))
	}
	return func(w http.ResponseWriter, r *http.Request) {

		unsignedID, err := cookies.Get(r, setting.Name, setting.Secret)
		if err != nil {
			//authErrors.WithError(err)
			panic(err)
		}

		sessionValue, err := au.store.Get(unsignedID)
		if err != nil {
			//handle err...
			panic(err)
		}
		//Can set explicit true because if it was not valid then
		//cookies.GetSigned would've returned an error
		session := session.NewFromCookie(unsignedID, sessionValue)

		h.ServeHTTP(w, r.WithContext(au.enrichContext(r, session)))
		return
	}
}

func (au *Auth) enrichContext(r *http.Request, s *session.AuthSession) context.Context {
	return context.WithValue(r.Context(), session.CtxKey, s)
}
