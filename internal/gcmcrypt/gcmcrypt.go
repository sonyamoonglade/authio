package gcmcrypt

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"encoding/hex"
	"fmt"
	"io"
)

const (
	nonceLen = 12
	keyLen   = 16
)

//KeyFromString goes char-by-char and generates a []byte key
//If len(str) > keyLen - extra characters will be ommited.
//Otherwise, if len(str) < keyLen - remaining characters will be just ""(empty str)
func KeyFromString(s string) [keyLen]byte {
	var key [keyLen]byte
	for i, ch := range s {
		if i == keyLen {
			return key
		}
		key[i] = byte(ch)
	}
	return key
}

func Encrypt(key [keyLen]byte, value string) (string, error) {

	block, err := aes.NewCipher(key[:])
	if err != nil {
		return "", err
	}

	//randomly generated byte slice by rand.Reader
	nonce := make([]byte, nonceLen)
	_, err = io.ReadFull(rand.Reader, nonce)
	if err != nil {
		return "", err
	}

	msg := []byte(value)

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return "", err
	}

	out := gcm.Seal(nil, nonce, msg, nil)
	magic := append(nonce, out...)

	return hex.EncodeToString(magic), nil
}

func Decrypt(key [16]byte, encryptedValue string) (string, error) {

	buff, err := hex.DecodeString(encryptedValue)
	if err != nil {
		return "", fmt.Errorf("could not decode value: %v", err)
	}

	nonce := buff[:nonceLen]

	block, err := aes.NewCipher(key[:])
	if err != nil {
		return "", err
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return "", err
	}

	out, err := gcm.Open(nil, nonce, buff[nonceLen:], nil)
	if err != nil {
		return "", err
	}

	return string(out), nil
}

func tob64(in string) string {
	src := []byte(in)

	return base64.StdEncoding.EncodeToString(src)
}

func fromb64(in string) ([]byte, error) {
	out, err := base64.StdEncoding.DecodeString(in)
	if err != nil {
		return nil, err
	}
	return out, nil
}