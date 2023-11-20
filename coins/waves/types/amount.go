/*
*
MIT License

Copyright (c) 2018 WavesPlatform
*/
package types

type Amount struct {
	AssetId []byte `json:"asset_id,omitempty"`
	Amount  int64  `json:"amount,omitempty"`
}
