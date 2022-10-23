package hash

import (
	"crypto/sha1"
	"encoding/hex"
)

func SHA1(ID string) string {
	sha := sha1.New()
	sha.Write([]byte(ID))
	return hex.EncodeToString(sha.Sum(nil))
}
func Compare(hashed string, actual string) bool {
	return hashed == SHA1(actual)
}
