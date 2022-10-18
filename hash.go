package authio

import "crypto/sha1"

func SHA1(ID string) string {
	s := sha1.New()
	return string(s.Sum([]byte(ID)))
}
