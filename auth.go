package authio

import (
	"fmt"

	"github.com/sonyamoonglade/authio/internal/cookie"
	"github.com/sonyamoonglade/authio/internal/logger"
	"github.com/sonyamoonglade/authio/internal/store"
)

type Auth struct {
	logger   logger.Logger
	store    store.Store
	settings map[string]*cookie.Setting
}

func NewAuth(logger logger.Logger, store store.Store) *Auth {
	au := &Auth{
		logger:   logger,
		store:    store,
		settings: make(map[string]*cookie.Setting),
	}

	au.settings[cookie.DefaultLabel] = cookie.DefaultSetting

	return au
}

func (au *Auth) EatBuilder(b *SettingBuilder) error {
	for label, setting := range b.settings {
		_, ok := au.settings[label]
		if ok {
			return fmt.Errorf("label %s already exist\n", label)
		}
		au.settings[label] = setting
	}

	return nil
}
