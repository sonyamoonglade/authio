package gcmcrypt

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestShouldConvertTob64AndViceVersa(t *testing.T) {

	v := "MY-SUPER-STRING-WITH-__SOME_==CHARS_@@12U3Y2178!@#*&!@*#&!@#&!"

	b64 := tob64(v)

	outb64, err := fromb64(b64)
	require.NoError(t, err)

	require.Equal(t, v, string(outb64))
}

func TestShouldEncryptAndDecryptValue(t *testing.T) {

	randomStr := "asdkfjdaskasjfhy234hkas"
	key := KeyFromString(randomStr)

	value := "MY-πππSUPEππR-π«π½π«π½π«π½SEπ«π½π«π½π«π½CRET-S!@#!&@#^!@#&^TRINGππ" //with emodji and other chars

	encrypted, err := Encrypt(key, value)
	require.NoError(t, err)

	decrypted, err := Decrypt(key, encrypted)
	require.NoError(t, err)

	require.Equal(t, value, decrypted)
}
