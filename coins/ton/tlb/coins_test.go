/**
Author： https://github.com/xssnick/tonutils-go
*/

package tlb

import (
	"math/big"
	"reflect"
	"testing"
)

func TestCoins_MarshalJSON(t *testing.T) {
	tests := []struct {
		name    string
		coins   Coins
		want    string
		wantErr bool
	}{
		{
			name: "0.123456789 TON",
			coins: Coins{
				decimals: 9,
				val:      big.NewInt(123_456_789),
			},
			want:    "\"123456789\"",
			wantErr: false,
		},
		{
			name: "1 TON",
			coins: Coins{
				decimals: 9,
				val:      big.NewInt(1_000_000_000),
			},
			want:    "\"1000000000\"",
			wantErr: false,
		},
		{
			name: "123 TON",
			coins: Coins{
				decimals: 9,
				val:      big.NewInt(123_000_000_000),
			},
			want:    "\"123000000000\"",
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.coins.MarshalJSON()
			if (err != nil) != tt.wantErr {
				t.Errorf("MarshalJSON() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			wantBytes := []byte(tt.want)
			if !reflect.DeepEqual(got, wantBytes) {
				t.Errorf("MarshalJSON() got = %v, want %v", string(got), tt.want)
			}
		})
	}
}

func TestCoins_UnmarshalJSON(t *testing.T) {
	tests := []struct {
		name    string
		data    string
		want    Coins
		wantErr bool
	}{
		{
			name:    "empty invalid",
			data:    "",
			wantErr: true,
		},
		{
			name:    "empty",
			data:    "\"\"",
			wantErr: true,
		},
		{
			name:    "invalid",
			data:    "\"123a\"",
			wantErr: true,
		},
		{
			name: "0.123456789 TON",
			data: "\"123456789\"",
			want: Coins{
				decimals: 9,
				val:      big.NewInt(123_456_789),
			},
		},
		{
			name: "1 TON",
			data: "\"1000000000\"",
			want: Coins{
				decimals: 9,
				val:      big.NewInt(1_000_000_000),
			},
		},
		{
			name: "123 TON",
			data: "\"123000000000\"",
			want: Coins{
				decimals: 9,
				val:      big.NewInt(123_000_000_000),
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var coins Coins

			err := coins.UnmarshalJSON([]byte(tt.data))
			if (err != nil) != tt.wantErr {
				t.Errorf("UnmarshalJSON() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !reflect.DeepEqual(coins, tt.want) {
				t.Errorf("UnmarshalJSON() got = %v, want %v", coins, tt.want)
			}
		})
	}
}
