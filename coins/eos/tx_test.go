package eos

import (
	"encoding/hex"
	"github.com/okx/go-wallet-sdk/coins/eos/types"
	"reflect"
	"testing"
)

func TestNewTransaction(t *testing.T) {
	blockID, err := hex.DecodeString("0012cf6247be7e2050090bd83b473369b705ba1d280cd55d3aef79998c784b9b")
	if err != nil {
		t.Error(err)
		return
	}
	opt := &types.TxOptions{
		ChainID:     []byte("eosio"),
		HeadBlockID: blockID,
	}
	type args struct {
		actions []*types.Action
		opts    *types.TxOptions
	}
	tests := []struct {
		name string
		args args
		want *types.Transaction
	}{
		{
			name: "test",
			args: args{
				opts: opt,
			},
			want: &types.Transaction{
				TransactionHeader: types.TransactionHeader{
					RefBlockNum:    uint16(0xcf62),
					RefBlockPrefix: uint32(3624601936),
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tx := NewTransaction(tt.args.actions, tt.args.opts)
			// check tx.RefBlockNum
			if tx.RefBlockNum != tt.want.RefBlockNum {
				t.Errorf("RefBlockNum = %d, want %d", tx.RefBlockNum, tt.want.RefBlockNum)
			}
			// check tx.RefBlockPrefix
			if tx.RefBlockPrefix != tt.want.RefBlockPrefix {
				t.Errorf("RefBlockPrefix = %d, want %d", tx.RefBlockPrefix, tt.want.RefBlockPrefix)
			}
		})
	}
}

func TestNewTransactionWithParams(t *testing.T) {
	quantity, _ := types.NewEOSAssetFromString("1.0000 EOS")
	type args struct {
		from     string
		to       string
		memo     string
		quantity types.Asset
		opts     *types.TxOptions
	}
	tests := []struct {
		name string
		args args
		want *types.Transaction
	}{
		{
			name: "test",
			args: args{
				from:     "dubuqingfeng",
				to:       "",
				memo:     "",
				quantity: quantity,
				opts: &types.TxOptions{
					ChainID: []byte("eosio"),
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tx := NewTransactionWithParams(tt.args.from, tt.args.to, tt.args.memo, tt.args.quantity, tt.args.opts)
			// check tx.Actions
			if len(tx.Actions) != 1 {
				t.Errorf("Actions = %d, want %d", len(tx.Actions), 1)
			}
			// check tx.Actions[0].Account
			if tx.Actions[0].Account != "eosio.token" {
				t.Errorf("Actions[0].Account = %s, want %s", tx.Actions[0].Account, "eosio.token")
			}
			// check tx.Actions[0].Name
			if tx.Actions[0].Name != "transfer" {
				t.Errorf("Actions[0].Name = %s, want %s", tx.Actions[0].Name, "transfer")
			}
			// check tx.Actions[0].Authorizations
			if len(tx.Actions[0].Authorization) != 1 {
				t.Errorf("Actions[0].Authorizations = %d, want %d", len(tx.Actions[0].Authorization), 1)
			}
			// check tx.Actions[0].Authorizations[0].Actor
			if tx.Actions[0].Authorization[0].Actor != types.AN(tt.args.from) {
				t.Errorf("Actions[0].Authorizations[0].Actor = %s, want %s", tx.Actions[0].Authorization[0].Actor, tt.args.from)
			}
			// check tx.Actions[0].Authorizations[0].Permission
			if tx.Actions[0].Authorization[0].Permission != "active" {
				t.Errorf("Actions[0].Authorizations[0].Permission = %s, want %s", tx.Actions[0].Authorization[0].Permission, "active")
			}
			// check tx.Actions[0].Data
			var transfer = tx.Actions[0].ActionData.Data.(types.Transfer)
			if transfer.From != types.AN(tt.args.from) {
				t.Errorf("transfer.From = %s, want %s", transfer.From, tt.args.from)
			}
			if transfer.To != types.AN(tt.args.to) {
				t.Errorf("transfer.To = %s, want %s", transfer.To, tt.args.to)
			}
			if transfer.Quantity != tt.args.quantity {
				t.Errorf("transfer.Quantity = %s, want %s", transfer.Quantity, tt.args.quantity)
			}
			if transfer.Memo != tt.args.memo {
				t.Errorf("transfer.Memo = %s, want %s", transfer.Memo, tt.args.memo)
			}
		})
	}
}

func TestSignTransaction(t *testing.T) {
	type args struct {
		wifKey      string
		tx          *types.Transaction
		chainID     types.Checksum256
		compression types.CompressionType
	}
	p, _ := GenerateKeyPair()
	tests := []struct {
		name    string
		args    args
		want    *types.SignedTransaction
		want1   *types.PackedTransaction
		wantErr bool
	}{
		{
			name: "test1",
			args: args{
				wifKey:      p,
				tx:          &types.Transaction{},
				chainID:     types.Checksum256{},
				compression: types.CompressionNone,
			},
			want:    nil,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, got1, err := SignTransaction(tt.args.wifKey, tt.args.tx, tt.args.chainID, tt.args.compression)
			if (err != nil) != tt.wantErr {
				t.Errorf("SignTransaction() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("SignTransaction() got = %v, want %v", got, tt.want)
			}
			if !reflect.DeepEqual(got1, tt.want1) {
				t.Errorf("SignTransaction() got1 = %v, want %v", got1, tt.want1)
			}
		})
	}
}

func TestSetRefBlock(t *testing.T) {
	tx := &types.Transaction{}
	blockID, err := hex.DecodeString("0012cf6247be7e2050090bd83b473369b705ba1d280cd55d3aef79998c784b9b")
	if err != nil {
		t.Error(err)
		return
	}
	tx.Fill(blockID, 0, 0, 0, 0)
	if tx.RefBlockNum != uint16(0xcf62) {
		t.Errorf("RefBlockNum = %d, want %d", tx.RefBlockNum, 0xcf62)
	}
	if tx.RefBlockPrefix == uint32(0xbe7e2050) {
		t.Errorf("RefBlockPrefix = %d, want %d", tx.RefBlockPrefix, 0xbe7e2050)
	}
	if tx.RefBlockPrefix != uint32(0xd80b0950) {
		t.Errorf("RefBlockPrefix = %d, want %d", tx.RefBlockPrefix, 0x90bd83b4)
	}
}
