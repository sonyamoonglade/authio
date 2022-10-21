package authio

import (
	"errors"

	"github.com/sonyamoonglade/authio/cookies"
	"github.com/sonyamoonglade/authio/store"
)

var (
	ErrSettingExist = errors.New("setting already exists")
)

type AuthBuilder struct {
	logger   Logger
	store    store.Store
	settings map[string]*cookies.Setting
	pf       store.ParseFromStoreFunc
}

func NewAuthBuilder() *AuthBuilder {
	return &AuthBuilder{
		logger:   nil,
		store:    nil,
		settings: make(map[string]*cookies.Setting),
		pf:       nil,
	}
}

func (b *AuthBuilder) AddCookieSetting(setting *cookies.Setting) *AuthBuilder {
	b.settings[setting.Label] = setting
	return b
}

func (b *AuthBuilder) UseStore(store store.Store) *AuthBuilder {
	b.store = store
	return b
}

func (b *AuthBuilder) UseLogger(logger Logger) *AuthBuilder {
	b.logger = logger
	return b
}

func (b *AuthBuilder) Build() *Auth {

	return newAuth(b.logger, b.store, b.settings)
}

type Auth struct {
	logger   Logger
	store    store.Store
	settings map[string]*cookies.Setting
}

//todo: check nils
func newAuth(logger Logger,
	store store.Store,
	settings map[string]*cookies.Setting) *Auth {

	au := &Auth{
		logger:   logger,
		store:    store,
		settings: settings,
	}

	au.settings[cookies.DefaultLabel] = cookies.DefaultSetting

	return au
}
