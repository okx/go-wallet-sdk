package tezos

import (
	"github.com/stretchr/testify/require"
	"testing"
)

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
		{name: "ValidAddress", args: args{addr: "tz1QMREVfenun5z18veejSRjiM44y1QfrfXA"}, want: true, wantErr: false},
		{name: "ValidAddress", args: args{addr: "tz2Athio4Vrc1Ge2GdpFAY7dFceVN5eMy6oj"}, want: true, wantErr: false},
		{name: "ValidAddress", args: args{addr: "tz3YjPMAXtsCQYzVRoYFPg4n32WAvXEh6QQb"}, want: true, wantErr: false},
		{name: "ValidAddress", args: args{addr: "KT1GyeRktoGPEKsWpchWguyy8FAf3aNHkw2T"}, want: true, wantErr: false},
		{name: "ValidAddress", args: args{addr: "error"}, want: false, wantErr: true},
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

func TestGetAddress(t *testing.T) {
	type args struct {
		publicKey string
	}
	tests := []struct {
		name    string
		args    args
		want    string
		wantErr bool
	}{
		{name: "address", args: args{publicKey: "edpkucde3WUTR2s6KgDBwvR7NiezGyHNj1aGz6WrJg6SeZWHNjDA8N"},
			want: "tz1QMREVfenun5z18veejSRjiM44y1QfrfXA", wantErr: false},
		{name: "address", args: args{publicKey: "sppk7aAV5AjmQPcph9SrrKBBeFwj15kMvnByjbvb9mqsTMgUm1ZoHxK"},
			want: "tz2Athio4Vrc1Ge2GdpFAY7dFceVN5eMy6oj", wantErr: false},
		{name: "address", args: args{publicKey: "p2pk68CeMSnZ8MhrW6zCJzGfS2VTsFUKK5GwB7Hem3UUuyQH2kHHeij"},
			want: "tz3YjPMAXtsCQYzVRoYFPg4n32WAvXEh6QQb", wantErr: false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := GetAddress(tt.args.publicKey)
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

func TestGenerateKeyPair(t *testing.T) {
	got, got1, err := GenerateKeyPair()
	require.NoError(t, err)
	t.Log(got)
	t.Log(got1)
	addr, err := GetAddressByPublicKey(got1)
	require.NoError(t, err)
	t.Log(addr)
}
