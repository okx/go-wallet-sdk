package waves

import (
	"github.com/okx/go-wallet-sdk/coins/waves/crypto"
	"github.com/okx/go-wallet-sdk/coins/waves/types"
	"github.com/stretchr/testify/assert"
	"reflect"
	"testing"
)

func Test_newAddressFromPublicKeyHash(t *testing.T) {
	addr, _ := NewAddressFromString("3P22DpDfBvLr9E7WEfC5sGnWCyQ2M9wods3")
	type args struct {
		scheme     byte
		pubKeyHash []byte
	}
	tests := []struct {
		name    string
		args    args
		want    types.WavesAddress
		wantErr bool
	}{
		{
			name: "test1",
			args: args{
				scheme: MainNetScheme,
				pubKeyHash: []byte{0x01, 0x02, 0x03, 0x04, 0x05, 0x06, 0x07, 0x08, 0x09,
					0x0a, 0x0b, 0x0c, 0x0d, 0x0e, 0x0f, 0x10, 0x11, 0x12, 0x13, 0x14, 0x15, 0x16, 0x17,
					0x18, 0x19, 0x1a, 0x1b, 0x1c, 0x1d, 0x1e, 0x1f, 0x20},
			},
			want:    addr,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := types.NewAddressFromPublicKeyHash(tt.args.scheme, tt.args.pubKeyHash)
			if (err != nil) != tt.wantErr {
				t.Errorf("newAddressFromPublicKeyHash() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("newAddressFromPublicKeyHash() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGenerateKeyPair(t *testing.T) {
	privKey, pubKey, err := GenerateKeyPair()
	if err != nil {
		t.Errorf("GenerateKeyPair() error = %v", err)
		return
	}
	assert.Equal(t, 44, len(privKey))
	assert.Equal(t, 44, len(pubKey))
}

func TestPrivateKeyPublicKey(t *testing.T) {
	type args struct {
		privateKey string
	}
	tests := []struct {
		name    string
		args    args
		want    string
		wantErr bool
	}{
		{
			name: "test1",
			args: args{
				privateKey: "5NBbF9dHDfuJw2WC8m3Am5kJwKMXbLmN2eh4Cmqsgo5w",
			},
			want:    "tMUA9XRwPTiUXCTmEvU6kFkqTFKxSpaAFvQwyAT29GR",
			wantErr: false,
		},
		{
			name: "test2",
			args: args{
				privateKey: "6QhEoSnJ12QDgeEAt3HYkPDBiYe15BArgSKWrV3DUctG",
			},
			want:    "GRcXDTsfpJZU6qUPkhjBX7dY1yKJ5mV2JJyWHWW1mUYK",
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			privateKey, err := crypto.NewSecretKeyFromBase58(tt.args.privateKey)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetAddress() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			publicKey := crypto.GeneratePublicKey(privateKey)
			if publicKey.String() != tt.want {
				t.Errorf("GetAddress() got = %v, want %v", publicKey.String(), tt.want)
			}
		})
	}
}

func TestGetAddress(t *testing.T) {
	pubKeyHash1, _ := crypto.NewPublicKeyFromBase58("2wySdbAsXi1bfAfMBKC1NcyyJemUWLM4R5ECwXJiADUx")
	pubKeyHash2, _ := crypto.NewPublicKeyFromBase58("tMUA9XRwPTiUXCTmEvU6kFkqTFKxSpaAFvQwyAT29GR")
	pubKeyHash3, _ := crypto.NewPublicKeyFromBase58("GRcXDTsfpJZU6qUPkhjBX7dY1yKJ5mV2JJyWHWW1mUYK")

	type args struct {
		scheme     byte
		pubKeyHash []byte
	}
	tests := []struct {
		name    string
		args    args
		want    string
		wantErr bool
	}{
		{
			name: "test1",
			args: args{
				scheme:     MainNetScheme,
				pubKeyHash: pubKeyHash1.Bytes(),
			},
			want:    "3P4ZvCY6W5WCMfPEdKCnWxqWWA8qdEDuJp3",
			wantErr: false,
		},
		{
			name: "test2",
			args: args{
				scheme:     TestNetScheme,
				pubKeyHash: pubKeyHash1.Bytes(),
			},
			want:    "3MrZ7FDCdwxojD5pNEwnZWTh9Gd4o9hAd3Y",
			wantErr: false,
		},
		{
			name: "test3",
			args: args{
				scheme:     TestNetScheme,
				pubKeyHash: pubKeyHash2.Bytes(),
			},
			want:    "3Mq7eCKTgNAoEag4eQVHZYGZKRNYKmodEpM",
			wantErr: false,
		},
		{
			name: "test4",
			args: args{
				scheme:     TestNetScheme,
				pubKeyHash: pubKeyHash3.Bytes(),
			},
			want:    "3NAorunHiZ5aJNQuyhZ3XBZy9Msc8pedYfA",
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := GetAddress(tt.args.scheme, tt.args.pubKeyHash)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetAddress() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("GetAddress() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestValidAddress(t *testing.T) {
	type args struct {
		addr string
	}
	tests := []struct {
		name    string
		args    args
		want    bool
		wantErr bool
	}{
		{
			name: "test1",
			args: args{
				addr: "3P4ZvCY6W5WCMfPEdKCnWxqWWA8qdEDuJp3",
			},
			want:    true,
			wantErr: false,
		},
		{
			name: "test2",
			args: args{
				addr: "3MrZ7FDCdwxojD5pNEwnZWTh9Gd4o9hAd3Y",
			},
			want:    true,
			wantErr: false,
		},
		{
			name: "test3",
			args: args{
				addr: "3Mq7eCKTgNAoEag4eQVHZYGZKRNYKmodEpM",
			},
			want:    true,
			wantErr: false,
		},
		{
			name: "test4",
			args: args{
				addr: "3NB2pUqjoavApZeAmdsVYS84hyRGXZpeytA",
			},
			want:    true,
			wantErr: false,
		},
		{
			name: "test5",
			args: args{
				addr: "edwards25519",
			},
			want:    false,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ValidAddress(tt.args.addr)
			if (err != nil) != tt.wantErr {
				t.Errorf("ValidAddress() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("ValidAddress() got = %v, want %v", got, tt.want)
			}
		})
	}
}
