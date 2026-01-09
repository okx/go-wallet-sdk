package cardano

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewAddress(t *testing.T) {
	mnemonic := "north bulb crunch need badge orient tissue web east scan invite energy canal solar eight"
	path := "m/1852'/1815'/0'/0/0"

	prvKeyHex, err := DerivePrvKey(mnemonic, path)
	assert.Nil(t, err)

	address, err := NewAddressFromPrvKey(prvKeyHex)
	assert.Nil(t, err)
	assert.Equal(t, "addr1qyxdgpwcqsrfsfv7gs3el47ym205hxaxnnpvs550czrjr8gr7z40zns2zm4kdd5jgxhawpstcgnyt4zdwzn4e9g6qmksvhsufu", address)

	pubKeyHex, err := PubKeyFromPrvKey(prvKeyHex)
	assert.Nil(t, err)
	assert.Equal(t, "28390ae3e51963cc935acd5b522a1dbbb6e06119ae3a9d93143acd30bc8b2a535de8f3bbab3225021a403adbb7a13e9973f356bdd154730fbc3ed16745210204", pubKeyHex)

	address, err = NewAddressFromPubKey(pubKeyHex)
	assert.Nil(t, err)
	assert.Equal(t, "addr1qyxdgpwcqsrfsfv7gs3el47ym205hxaxnnpvs550czrjr8gr7z40zns2zm4kdd5jgxhawpstcgnyt4zdwzn4e9g6qmksvhsufu", address)
}

func TestValidateAddress(t *testing.T) {
	isValid := ValidateAddress("addr1qyxdgpwcqsrfsfv7gs3el47ym205hxaxnnpvs550czrjr8gr7z40zns2zm4kdd5jgxhawpstcgnyt4zdwzn4e9g6qmksvhsufu")
	assert.True(t, isValid)

	isValid = ValidateAddress("addr1pyxdgpwcqsrfsfv7gs3el47ym205hxaxnnpvs550czrjr8gr7z40zns2zm4kdd5jgxhawpstcgnyt4zdwzn4e9g6qmksvhsufu")
	assert.False(t, isValid)

	isValid = ValidateAddress("addr1v92fkk3qu3y68cu5ka38qhmyhx3xhxgpxqp6907m5guevlqs8vk7u")
	assert.True(t, isValid)

	isValid = ValidateAddress("addr1s92fkk3qu3y68cu5ka38qhmyhx3xhxgpxqp6907m5guevlqs8vk7u")
	assert.False(t, isValid)
}
