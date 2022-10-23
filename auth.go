package authio

import (
	"errors"

	"github.com/sonyamoonglade/authio/cookies"
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
