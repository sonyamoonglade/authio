package authio

import (
	"testing"
	"time"

	"github.com/sonyamoonglade/authio/cookies"
	"github.com/sonyamoonglade/authio/store"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
)

//tdd
func TestCanBuildAuth(t *testing.T) {

	lg, _ := zap.NewProduction()
	logger := lg.Sugar()

	auth := NewAuthBuilder().
		AddCookieSetting(cookies.DefaultSetting).
		UseLogger(logger).
		UseStore(store.NewInMemoryStore(&store.Config{
			EntryTTL:         time.Hour * 24,
			OverflowStrategy: store.LRU,
			ParseFunc:        store.ToInt64,
		}, 100)).
		Build()

	require.NotNil(t, auth)
	require.NotNil(t, auth.settings)
	require.NotNil(t, auth.logger)
	require.NotNil(t, auth.store)
	require.NotNil(t, auth.settings[cookies.DefaultLabel])
	require.Nil(t, auth.settings["jo-mama"])
}
