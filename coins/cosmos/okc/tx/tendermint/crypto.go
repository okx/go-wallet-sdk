package tendermint

// An address is a []byte, but hex-encoded even in JSON.
// []byte leaves us the option to change the address length.
// Use an alias so Unmarshal methods (with ptr receivers) are available too.
type Address = HexBytes

type PubKey interface {
	Address() Address
	Bytes() []byte
	Equals(PubKey) bool
}

type PrivKey interface {
	Bytes() []byte
	PubKey() PubKey
	Equals(PrivKey) bool
}
