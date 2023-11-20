package crypto

import (
	"crypto/ecdsa"
	"github.com/decred/dcrd/dcrec/secp256k1/v4"
	"math/big"
	"reflect"
	"testing"
)

func TestSecp256k1Key_PubKey(t *testing.T) {
	type fields struct {
		PrivateKey *ecdsa.PrivateKey
	}
	tests := []struct {
		name   string
		fields fields
		want   []byte
	}{
		{
			name: "test",
			fields: fields{
				PrivateKey: &ecdsa.PrivateKey{
					PublicKey: ecdsa.PublicKey{
						Curve: secp256k1.S256(),
						X:     big.NewInt(1),
						Y:     big.NewInt(2),
					},
					D: big.NewInt(3),
				},
			},
			want: []byte{2, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			k := &Secp256k1Key{
				PrivateKey: tt.fields.PrivateKey,
			}
			if got := k.PubKey(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("PubKey() = %v, want %v", got, tt.want)
			}
		})
	}
}
