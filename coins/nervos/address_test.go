package nervos

import (
	"github.com/okx/go-wallet-sdk/coins/nervos/types"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestConvertScriptToBech32mFullAddress(t *testing.T) {
	type args struct {
		mode   Mode
		script *types.Script
	}
	tests := []struct {
		name    string
		args    args
		want    string
		wantErr bool
	}{
		{
			name: "testnet",
			args: args{
				mode: Testnet,
				script: &types.Script{
					// CodeHash: "0x3419a1c09eb2567f6552ee7a8ecffd64155cffe0f1796e6e61ec088d740c1356",
					HashType: "data",
				}},
			// want: "ckt1qyqrqzqxq2q9qwq3q4q5q6q7q8q9q0qaqdqfqgqhqjqkqlqmqnqoqpqq",
			want: "ckt1qqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqgaqanf",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ConvertScriptToBech32mFullAddress(tt.args.mode, tt.args.script)
			if (err != nil) != tt.wantErr {
				t.Errorf("ConvertScriptToBech32mFullAddress() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("ConvertScriptToBech32mFullAddress() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGenerateTestnetAddress(t *testing.T) {
	got, err := GenerateTestnetAddress()
	require.NoError(t, err)
	t.Log(got)
}

func TestValidateAddress(t *testing.T) {
	type args struct {
		address string
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			name: "testnet",
			args: args{
				address: "ckt1qyqrqzqxq2q9qwq3q4q5q6q7q8q9q0qaqdqfqgqhqjqkqlqmqnqoqpqq",
			},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := ValidateAddress(tt.args.address); got != tt.want {
				t.Errorf("ValidateAddress() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGenerateAddressByPrivateKey(t *testing.T) {
	type args struct {
		mode       string
		privateKey string
	}
	tests := []struct {
		name    string
		args    args
		want    string
		wantErr bool
	}{
		{
			name: "testnet",
			args: args{
				mode:       "ckt",
				privateKey: "0171ecab8a308cd26fef99efb7ea02fa17ec9c210d8e9f6e32543694a6623ece",
			},
			want:    "ckt1qzda0cr08m85hc8jlnfp3zer7xulejywt49kt2rr0vthywaa50xwsq08pk6ldw7944vqmvulq555739qnlpap8sglxd9q",
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := GenerateAddressByPrivateKey(tt.args.mode, tt.args.privateKey)
			if (err != nil) != tt.wantErr {
				t.Errorf("GenerateAddressByPrivateKey() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("GenerateAddressByPrivateKey() got = %v, want %v", got, tt.want)
			}
		})
	}
}
