package cookie

import (
	"net/http"
	"time"
)

var DefaultSetting = &Setting{
	HttpOnly:  true,
	Secure:    false,
	SameSite:  http.SameSiteDefaultMode,
	ExpiresAt: time.Now().Add(time.Hour * 1),
	Path:      "",
	Domain:    "",
}

var DefaultLabel = "default"

type Setting struct {
	HttpOnly  bool
	Secure    bool
	SameSite  http.SameSite
	Path      string
	Domain    string
	ExpiresAt time.Time
}
