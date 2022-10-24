package hash

import (
	"crypto/sha1"
	"encoding/hex"
)

func SHA1(v string) string {
	s := sha1.New()
	s.Write([]byte(v))
	return hex.EncodeToString(s.Sum(nil))
}
