/**
Author： https://github.com/xssnick/tonutils-go
*/

package wallet

import (
	"bytes"
	"crypto/ed25519"
	"fmt"

	"github.com/okx/go-wallet-sdk/coins/ton/address"
	"github.com/okx/go-wallet-sdk/coins/ton/tlb"
	"github.com/okx/go-wallet-sdk/coins/ton/tvm/cell"
)

const DefaultSubwallet = 698983191
const VenomDefaultSubwallet = 1269378442

func AddressFromPubKey(key ed25519.PublicKey, version VersionConfig, subwallet uint32) (*address.Address, error) {
	state, err := GetStateInit(key, version, subwallet)
	if err != nil {
		return nil, fmt.Errorf("failed to get state: %w", err)
	}

	stateCell, err := tlb.ToCell(state)
	if err != nil {
		return nil, fmt.Errorf("failed to get state cell: %w", err)
	}

	addr := address.NewAddress(0, 0, stateCell.Hash())
	addr.SetBounce(false)
	return addr, nil
}

func GetWalletVersion(account *tlb.Account) Version {
	if !account.IsActive || account.State.Status != tlb.AccountStatusActive {
		return Unknown
	}

	for v := range walletCodeHex {
		code, ok := walletCode[v]
		if !ok {
			continue
		}
		if bytes.Equal(account.Code.Hash(), code.Hash()) {
			return v
		}
	}

	return Unknown
}

func GetStateInit(pubKey ed25519.PublicKey, version VersionConfig, subWallet uint32) (*tlb.StateInit, error) {
	switch version.(type) {
	case ConfigV5R1Final:
		subWallet = 0
	}
	var ver Version
	switch v := version.(type) {
	case Version:
		ver = v
		switch ver {
		case HighloadV3:
			return nil, fmt.Errorf("use ConfigHighloadV3 for highload v3 spec")
		case V5R1Final:
			return nil, fmt.Errorf("use ConfigV5R1Final for V5 spec")
		}
	case ConfigHighloadV3:
		ver = HighloadV3
	case ConfigV5R1Final:
		ver = V5R1Final
	}

	code, ok := walletCode[ver]
	if !ok {
		return nil, fmt.Errorf("cannot get code: %w", ErrUnsupportedWalletVersion)
	}

	var data *cell.Cell
	switch ver {
	case V3R1, V3R2, VenomV3:
		data = cell.BeginCell().
			MustStoreUInt(0, 32).                 // seqno
			MustStoreUInt(uint64(subWallet), 32). // sub wallet
			MustStoreSlice(pubKey, 256).
			EndCell()
	case V4R1, V4R2:
		data = cell.BeginCell().
			MustStoreUInt(0, 32). // seqno
			MustStoreUInt(uint64(subWallet), 32).
			MustStoreSlice(pubKey, 256).
			MustStoreDict(nil). // empty dict of plugins
			EndCell()
	case V5R1Final:
		config := version.(ConfigV5R1Final)

		// Create WalletId instance
		walletId := V5R1ID{
			NetworkGlobalID: config.NetworkGlobalID, // -3 Testnet, -239 Mainnet
			WorkChain:       config.Workchain,
			SubwalletNumber: uint16(subWallet),
			WalletVersion:   0, // Wallet Version
		}

		data = cell.BeginCell().
			MustStoreBoolBit(true).                           // storeUint(1, 1) - boolean flag for context type
			MustStoreUInt(0, 32).                             // Sequence number, hardcoded as 0
			MustStoreUInt(uint64(walletId.Serialized()), 32). // Serializing WalletId into 32-bit integer
			MustStoreSlice(pubKey, 256).                      // Storing the public key
			MustStoreDict(nil).                               // Storing an empty plugins dictionary
			EndCell()
	case HighloadV2R2, HighloadV2Verified:
		data = cell.BeginCell().
			MustStoreUInt(uint64(subWallet), 32).
			MustStoreUInt(0, 64). // last cleaned
			MustStoreSlice(pubKey, 256).
			MustStoreDict(nil). // old queries
			EndCell()
	case HighloadV3:
		timeout := version.(ConfigHighloadV3).MessageTTL
		if timeout >= 1<<22 {
			return nil, fmt.Errorf("too big timeout")
		}

		data = cell.BeginCell().
			MustStoreSlice(pubKey, 256).
			MustStoreUInt(uint64(subWallet), 32).
			MustStoreUInt(0, 66).
			MustStoreUInt(uint64(timeout), 22).
			EndCell()
	default:
		return nil, ErrUnsupportedWalletVersion
	}

	return &tlb.StateInit{
		Data: data,
		Code: code,
	}, nil
}
