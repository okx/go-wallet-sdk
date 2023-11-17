package eos

import (
	"github.com/eoscanada/eos-go/ecc"
	"github.com/stretchr/testify/require"
	"reflect"
	"testing"
)

func TestSigDigest(t *testing.T) {
	type args struct {
		chainID         []byte
		payload         []byte
		contextFreeData []byte
	}
	tests := []struct {
		name string
		args args
		want []byte
	}{
		{
			name: "test1",
			args: args{
				chainID:         []byte{0x01, 0x02, 0x03, 0x04, 0x05, 0x06, 0x07, 0x08, 0x09, 0x0a, 0x0b, 0x0c, 0x0d, 0x0e, 0x0f, 0x10, 0x11, 0x12, 0x13, 0x14, 0x15, 0x16, 0x17, 0x18, 0x19, 0x1a, 0x1b, 0x1c, 0x1d, 0x1e, 0x1f, 0x20},
				payload:         []byte{0x01, 0x02, 0x03, 0x04, 0x05, 0x06, 0x07, 0x08, 0x09, 0x0a, 0x0b, 0x0c, 0x0d, 0x0e, 0x0f, 0x10, 0x11, 0x12, 0x13, 0x14, 0x15, 0x16, 0x17, 0x18, 0x19, 0x1a, 0x1b, 0x1c, 0x1d, 0x1e, 0x1f, 0x20},
				contextFreeData: []byte{0x01, 0x02, 0x03, 0x04, 0x05, 0x06, 0x07, 0x08, 0x09, 0x0a, 0x0b, 0x0c, 0x0d, 0x0e, 0x0f, 0x10, 0x11, 0x12, 0x13, 0x14, 0x15, 0x16, 0x17, 0x18, 0x19, 0x1a, 0x1b, 0x1c, 0x1d, 0x1e, 0x1f, 0x20},
			},
			want: []byte{228, 228, 43, 61, 42, 228, 34, 3, 94, 76, 198, 84, 29, 104, 9, 110,
				184, 33, 188, 247, 180, 222, 142, 242, 14, 223, 81, 242, 78, 71, 231, 242},
		},
		{
			name: "test2 chainId is nil",
			args: args{
				chainID:         nil,
				payload:         []byte{0x01, 0x02, 0x03, 0x04, 0x05, 0x06, 0x07, 0x08, 0x09, 0x0a, 0x0b, 0x0c, 0x0d, 0x0e, 0x0f, 0x10, 0x11, 0x12, 0x13, 0x14, 0x15, 0x16, 0x17, 0x18, 0x19, 0x1a, 0x1b, 0x1c, 0x1d, 0x1e, 0x1f, 0x20},
				contextFreeData: []byte{0x01, 0x02, 0x03, 0x04, 0x05, 0x06, 0x07, 0x08, 0x09, 0x0a, 0x0b, 0x0c, 0x0d, 0x0e, 0x0f, 0x10, 0x11, 0x12, 0x13, 0x14, 0x15, 0x16, 0x17, 0x18, 0x19, 0x1a, 0x1b, 0x1c, 0x1d, 0x1e, 0x1f, 0x20},
			},
			want: []byte{148, 200, 233, 60, 109, 227, 5, 44, 133, 152, 255, 6, 249,
				55, 190, 69, 27, 212, 243, 86, 160, 70, 72, 191, 159, 36, 117, 155, 39, 250, 215, 77},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := SigDigest(tt.args.chainID, tt.args.payload, tt.args.contextFreeData); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("SigDigest() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestSigner_Add(t *testing.T) {
	type fields struct {
		Keys []*ecc.PrivateKey
	}
	type args struct {
		wifKey string
	}
	p, _ := GenerateKeyPair()
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			name: "test1",
			args: args{
				wifKey: p,
			},
			fields: fields{
				Keys: []*ecc.PrivateKey{},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			b := &Signer{
				Keys: tt.fields.Keys,
			}
			if err := b.Add(tt.args.wifKey); (err != nil) != tt.wantErr {
				t.Errorf("Add() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestNewSigner(t *testing.T) {
	p, _ := GenerateKeyPair()
	wifs := []string{
		p,
	}
	var keys []*ecc.PrivateKey
	for _, wif := range wifs {
		key, err := ecc.NewPrivateKey(wif)
		require.NoError(t, err)
		keys = append(keys, key)
	}
	var addWifs []string
	addWifs = []string{
		p,
	}
	signer := NewSigner(keys)
	if len(signer.Keys) != len(keys) {
		t.Errorf("NewSigner() error")
	}
	for _, wif := range addWifs {
		err := signer.Add(wif)
		require.NoError(t, err)
	}
}

func TestErrorTxSign(t *testing.T) {
	signer, _ := NewSignerFromWIFs([]string{""})
	_, err := signer.Sign(nil, nil)
	require.Error(t, err)
}
