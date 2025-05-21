package sui

import (
	"bytes"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"github.com/okx/go-wallet-sdk/crypto/ed25519"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
)

var (
	suiObjects = []*SuiObjectRef{{
		Digest:   "ESUg3nLfPmcMK2vf8kAyyX967w1whtgv8dk6pZhNHh6N",
		ObjectId: "0x0cd3e81f2130b922a25f113f017d8f76cae8c2f9d7ebed690e56754a0b3a5784",
		Version:  60784,
	}, {
		Digest:   "A3Nk1uPDmLLaYgBDPM235ZhTHBJVhXUD74mRVn9zZx4Z",
		ObjectId: "0x1d892361074249073f82613ad08b387bdec26185ea013d379f5ae0bb2a611ebc",
		Version:  60785,
	}, {
		Digest:   "2bacji3hre1MZiatXVjqQ1yVZXXfq5rw7yPyHwSH9CPj",
		ObjectId: "0x8565ab3ac7072abd8f7a0d4e81974a0c8669defaac9340dda7083e62542cd2f9",
		Version:  60783,
	}}

	seed            = "uemYAwkvsf/a7q2DdoMKNHWP7DlDhLLmgUh6coTtp94="
	recipient       = "0x19a8a5cadebbe8f73364cb1ce4d9aa2d0ffb62a3098c7ff93bd0426960f56406"
	sender          = "0x19a8a5cadebbe8f73364cb1ce4d9aa2d0ffb62a3098c7ff93bd0426960f56406"
	gasBudget       = uint64(2976000)
	gasPrice        = uint64(1000)
	expectTxBytes   = "AAACAAgBAAAAAAAAAAAgGailyt676PczZMsc5NmqLQ/7YqMJjH/5O9BCaWD1ZAYCAgABAQAAAQECAAABAQAZqKXK3rvo9zNkyxzk2aotD/tiowmMf/k70EJpYPVkBgMM0+gfITC5IqJfET8BfY92yujC+dfr7WkOVnVKCzpXhHDtAAAAAAAAIMeuvZGR3WQ3pQiWLbk7muGO5aih1NL8t9/oK2hFUZUnHYkjYQdCSQc/gmE60Is4e97CYYXqAT03n1rguyphHrxx7QAAAAAAACCGVO1kgYBA1g3K7CLQPUgcez0JMvsNV07m3F3E0ZG5WoVlqzrHByq9j3oNToGXSgyGad76rJNA3acIPmJULNL5b+0AAAAAAAAgF7dxv5oQMrjUfoOZZGMMbDIzaIewGNL0Ze1UOMWKshIZqKXK3rvo9zNkyxzk2aotD/tiowmMf/k70EJpYPVkBugDAAAAAAAAAGktAAAAAAAA"
	expectSignature = "AKY+OqtyuSZtCCQ2xaKSbkoyJm4d2uieAcNbMHkZXdcbqRbZIZt1hLom+/tXwXhkVbk4shmecHFoj9ABl7b2XweXPl+82/o6Rt93H1ojjfMbkpJDm+Rnx1AAjN7Nvi7fnQ=="
	txJson          = `{"version":1,"sender":"0x19a8a5cadebbe8f73364cb1ce4d9aa2d0ffb62a3098c7ff93bd0426960f56406","expiration":{"Epoch":100},"gasConfig":{"price":"1","budget":"110","payment":[{"digest":"3Gei4443UwTNcGak4QXgeUvJjR7ZW7Qx8dPAURCx4hkG","objectId":"0xcde3e7ce68b62f5c0446b81946eeb443f62228ab121061afefbfac3fc231bfee","version":4}]},"inputs":[{"kind":"Input","value":100000,"index":0,"type":"pure"},{"kind":"Input","value":"0x19a8a5cadebbe8f73364cb1ce4d9aa2d0ffb62a3098c7ff93bd0426960f56406","index":1,"type":"pure"}],"transactions":[{"kind":"SplitCoins","coin":{"kind":"GasCoin"},"amounts":[{"kind":"Input","value":100000,"index":0,"type":"pure"}]},{"kind":"TransferObjects","objects":[{"kind":"Result","index":0}],"address":{"kind":"Input","value":"0x19a8a5cadebbe8f73364cb1ce4d9aa2d0ffb62a3098c7ff93bd0426960f56406","index":1,"type":"pure"}}]}`

	txJson2 = `{
	"version": 1,
	"sender": "0x19a8a5cadebbe8f73364cb1ce4d9aa2d0ffb62a3098c7ff93bd0426960f56406",
	"gasConfig": {
		"price": "1000",
		"budget": "100000000",
		"payment": [{
			"objectId": "0x26c26d20be986dbd99d33ff1b3a0bb16437039d3de85b4cb9d56a3f57066ef54",
			"digest": "CmebPaLJ6ggh6sTiDpfu5aSVZHHMGedyQCLQ7dwUFaPH",
			"version": 1978794
		}]
	},
	"inputs": [{
		"kind": "Input",
		"value": {
			"Object": {
				"ImmOrOwned": {
					"digest": "AyPDRVnTtwnw9tsTDcSpojN7hn8TiVr7S1SzuxD8b3CN",
					"version": 1978794,
					"objectId": "0x5b800c9e0a73512244d0f89bedad2c8c1a3fdf6799d10ad49a9bc89167cdbcab"
				}
			}
		},
		"index": 0,
		"type": "object"
	}, {
		"kind": "Input",
		"value": {
			"Pure": [1, 0, 0, 0, 0, 0, 0, 0]
		},
		"index": 1,
		"type": "pure"
	}, {
		"kind": "Input",
		"value": {
			"Pure": [25,168,165,202,222,187,232,247,51,100,203,28,228,217,170,45,15,251,98,163,9,140,127,249,59,208,66,105,96,245,100,6]
		},
		"index": 2,
		"type": "pure"
	}],
	"transactions": [{
		"kind": "SplitCoins",
		"coin": {
			"kind": "Input",
			"value": "0x5b800c9e0a73512244d0f89bedad2c8c1a3fdf6799d10ad49a9bc89167cdbcab",
			"index": 0,
			"type": "object"
		},
		"amounts": [{
			"kind": "Input",
			"value": 1,
			"index": 1,
			"type": "pure"
		}]
	}, {
		"kind": "TransferObjects",
		"objects": [{
			"kind": "Result",
			"index": 0
		}],
		"address": {
			"kind": "Input",
			"value": "0x19a8a5cadebbe8f73364cb1ce4d9aa2d0ffb62a3098c7ff93bd0426960f56406",
			"index": 2,
			"type": "pure"
		}
	}]
}`

	// Multiple token objects, first merged, and then split the specified amount to each other. In this example, there are two tokens in the inputs.
	txJson3 = `{
	"version": 1,
	"sender": "0x19a8a5cadebbe8f73364cb1ce4d9aa2d0ffb62a3098c7ff93bd0426960f56406",
	"gasConfig": {
		"price": "1000",
		"budget": "3000000",
		"payment": [{
			"objectId": "0x26c26d20be986dbd99d33ff1b3a0bb16437039d3de85b4cb9d56a3f57066ef54",
			"digest": "HSWkfqRq34bR4gS89xDbMBjSxANhpckJyxwdQyQLptUX",
			"version": 1978797
		}]
	},
	"inputs": [{
		"kind": "Input",
		"value": {
			"Object": {
				"ImmOrOwned": {
					"digest": "C1eBtxoHXE7fgdUbaqxrk2RCTR1XQ9BdreWbjPxwbPQ3",
					"version": 1978797,
					"objectId": "0x8db1eac813b301f0d6585dd9d470d6afc018e77b11acfd90892b3be3083b74e4"
				}
			}
		},
		"index": 0,
		"type": "object"
	}, {
		"kind": "Input",
		"value": {
			"Object": {
				"ImmOrOwned": {
					"digest": "7JaJkUiU5TKEKo4CtDbnVVUJMSxNiLQbJwFGmQBfwZze",
					"version": 1978797,
					"objectId": "0x4302d06bfae37024b523fae867d02ca697f2c374efa4d393de70044a1f04c81a"
				}
			}
		},
		"index": 1,
		"type": "object"
	}, {
		"kind": "Input",
		"value": {
			"Pure": [1, 0, 0, 0, 0, 0, 0, 0]
		},
		"index": 2,
		"type": "pure"
	}, {
		"kind": "Input",
		"value": {
			"Pure": [25,168,165,202,222,187,232,247,51,100,203,28,228,217,170,45,15,251,98,163,9,140,127,249,59,208,66,105,96,245,100,6]
		},
		"index": 3,
		"type": "pure"
	}],
	"transactions": [{
		"kind": "MergeCoins",
		"destination": {
			"kind": "Input",
			"value": "0x4302d06bfae37024b523fae867d02ca697f2c374efa4d393de70044a1f04c81a",
			"index": 1,
			"type": "object"
		},
		"sources": [{
			"kind": "Input",
			"value": "0x8db1eac813b301f0d6585dd9d470d6afc018e77b11acfd90892b3be3083b74e4",
			"index": 0,
			"type": "object"
		}]
	}, {
		"kind": "SplitCoins",
		"coin": {
			"kind": "Input",
			"value": "0x4302d06bfae37024b523fae867d02ca697f2c374efa4d393de70044a1f04c81a",
			"index": 1,
			"type": "object"
		},
		"amounts": [{
			"kind": "Input",
			"value": 1,
			"index": 2,
			"type": "pure"
		}]
	}, {
		"kind": "TransferObjects",
		"objects": [{
			"kind": "Result",
			"index": 1
		}],
		"address": {
			"kind": "Input",
			"value": "0x19a8a5cadebbe8f73364cb1ce4d9aa2d0ffb62a3098c7ff93bd0426960f56406",
			"index": 3,
			"type": "pure"
		}
	}]
}`

	// Multiple token objects, first merged, and then split the specified amount to each other. In this example, there are two tokens in the inputs..
	txJson4 = `{
	"version": 1,
	"sender": "0x19a8a5cadebbe8f73364cb1ce4d9aa2d0ffb62a3098c7ff93bd0426960f56406",
	"gasConfig": {
		"price": "1000",
		"budget": "3000000",
		"payment": [{
			"objectId": "0x26c26d20be986dbd99d33ff1b3a0bb16437039d3de85b4cb9d56a3f57066ef54",
			"digest": "AFQQ8sCHFQcoE9yBCFnhCojeieyiDexvonV2ojpcSZYq",
			"version": 1978799
		}, {
			"objectId": "0xe76d8b01fc25a7035c4748c7000f75450d726af4fe8ef00f6c1f800665ba0463",
			"digest": "7e7Y64CAiTZdjLBnz4Ncytt7t9V2XdWGbp4XGSYu1tRM",
			"version": 1978801
		}]
	},
	"inputs": [{
		"kind": "Input",
		"value": {
			"Object": {
				"ImmOrOwned": {
					"digest": "JAWCBZc7V6z4kBRBJCSc7qBoVvnXvkJrpPWPwFowQBBH",
					"version": 1978802,
					"objectId": "0x558f1d047d46214b1fdb756cd1c134d5960833c54dc719972b65f309931df4cd"
				}
			}
		},
		"index": 0,
		"type": "object"
	}, {
		"kind": "Input",
		"value": {
			"Object": {
				"ImmOrOwned": {
					"digest": "2k6TeD147XXphLzRrx1LRBTckeqtxgLyM8tpu5vmSKLR",
					"version": 1978800,
					"objectId": "0x98732f7a8388174dcd58873f39ab327ab1e6854ae56c77b9cee093128495ed70"
				}
			}
		},
		"index": 1,
		"type": "object"
	}, {
		"kind": "Input",
		"value": {
			"Object": {
				"ImmOrOwned": {
					"digest": "EPf6c5jNmQSZYtJ8rP2qMkvjz5fnK8JggapwjRCGDTYg",
					"version": 1978798,
					"objectId": "0xe67290a7b36c3a753a50b833aeedeeacba4058db6a970dec382c31b044814fdd"
				}
			}
		},
		"index": 2,
		"type": "object"
	}, {
		"kind": "Input",
		"value": {
			"Object": {
				"ImmOrOwned": {
					"digest": "EfkDithrS2xWznsiz5GqAFFKAeSoHLei32ocGFte3NfA",
					"version": 1978798,
					"objectId": "0x4302d06bfae37024b523fae867d02ca697f2c374efa4d393de70044a1f04c81a"
				}
			}
		},
		"index": 3,
		"type": "object"
	}, {
		"kind": "Input",
		"value": {
			"Pure": [1, 0, 0, 0, 0, 0, 0, 0]
		},
		"index": 4,
		"type": "pure"
	}, {
		"kind": "Input",
		"value": {
			"Pure": [25,168,165,202,222,187,232,247,51,100,203,28,228,217,170,45,15,251,98,163,9,140,127,249,59,208,66,105,96,245,100,6]
		},
		"index": 5,
		"type": "pure"
	}],
	"transactions": [{
		"kind": "MergeCoins",
		"destination": {
			"kind": "Input",
			"value": "0x4302d06bfae37024b523fae867d02ca697f2c374efa4d393de70044a1f04c81a",
			"index": 3,
			"type": "object"
		},
		"sources": [{
			"kind": "Input",
			"value": "0x558f1d047d46214b1fdb756cd1c134d5960833c54dc719972b65f309931df4cd",
			"index": 0,
			"type": "object"
		}, {
			"kind": "Input",
			"value": "0x98732f7a8388174dcd58873f39ab327ab1e6854ae56c77b9cee093128495ed70",
			"index": 1,
			"type": "object"
		}, {
			"kind": "Input",
			"value": "0xe67290a7b36c3a753a50b833aeedeeacba4058db6a970dec382c31b044814fdd",
			"index": 2,
			"type": "object"
		}]
	}, {
		"kind": "SplitCoins",
		"coin": {
			"kind": "Input",
			"value": "0x4302d06bfae37024b523fae867d02ca697f2c374efa4d393de70044a1f04c81a",
			"index": 3,
			"type": "object"
		},
		"amounts": [{
			"kind": "Input",
			"value": 1,
			"index": 4,
			"type": "pure"
		}]
	}, {
		"kind": "TransferObjects",
		"objects": [{
			"kind": "Result",
			"index": 1
		}],
		"address": {
			"kind": "Input",
			"value": "0x19a8a5cadebbe8f73364cb1ce4d9aa2d0ffb62a3098c7ff93bd0426960f56406",
			"index": 5,
			"type": "pure"
		}
	}]
}`
	txJson5 = `{"version":1,"sender":"0x19a8a5cadebbe8f73364cb1ce4d9aa2d0ffb62a3098c7ff93bd0426960f56406","gasConfig":{"price":"815","budget":"5582000","payment":[{"objectId":"0x26c26d20be986dbd99d33ff1b3a0bb16437039d3de85b4cb9d56a3f57066ef54","version":1978809,"digest":"8m7eNZYu9uZsbghPJy8gKxQq27iMPYF2uE1LHbBVJQyS"}]},"inputs":[{"kind":"Input","value":{"Pure":[1,0,0,0,0,0,0,0]},"index":0,"type":"pure"},{"kind":"Input","value":{"Pure":[2,0,0,0,0,0,0,0]},"index":1,"type":"pure"},{"kind":"Input","value":{"Pure":[3,0,0,0,0,0,0,0]},"index":2,"type":"pure"},{"kind":"Input","value":{"Pure":[4,0,0,0,0,0,0,0]},"index":3,"type":"pure"},{"kind":"Input","value":{"Pure":[5,0,0,0,0,0,0,0]},"index":4,"type":"pure"},{"kind":"Input","value":{"Pure":[6,0,0,0,0,0,0,0]},"index":5,"type":"pure"},{"kind":"Input","value":{"Pure":[7,0,0,0,0,0,0,0]},"index":6,"type":"pure"},{"kind":"Input","value":{"Pure":[25,168,165,202,222,187,232,247,51,100,203,28,228,217,170,45,15,251,98,163,9,140,127,249,59,208,66,105,96,245,100,6]},"index":7,"type":"pure"}],"transactions":[{"kind":"SplitCoins","coin":{"kind":"GasCoin"},"amounts":[{"kind":"Input","value":1,"index":0,"type":"pure"},{"kind":"Input","value":2,"index":1,"type":"pure"},{"kind":"Input","value":3,"index":2,"type":"pure"},{"kind":"Input","value":4,"index":3,"type":"pure"},{"kind":"Input","value":5,"index":4,"type":"pure"},{"kind":"Input","value":6,"index":5,"type":"pure"},{"kind":"Input","value":7,"index":6,"type":"pure"}]},{"kind":"TransferObjects","objects":[{"kind":"NestedResult","index":0,"resultIndex":0},{"kind":"NestedResult","index":0,"resultIndex":1},{"kind":"NestedResult","index":0,"resultIndex":2},{"kind":"NestedResult","index":0,"resultIndex":3},{"kind":"NestedResult","index":0,"resultIndex":4},{"kind":"NestedResult","index":0,"resultIndex":5},{"kind":"NestedResult","index":0,"resultIndex":6}],"address":{"kind":"Input","value":"0x19a8a5cadebbe8f73364cb1ce4d9aa2d0ffb62a3098c7ff93bd0426960f56406","index":7,"type":"pure"}}]}`

	mergeJson6 = `{"version":1,"sender":"0x19a8a5cadebbe8f73364cb1ce4d9aa2d0ffb62a3098c7ff93bd0426960f56406","gasConfig":{"price":"815","budget":"1630000","payment":[{"objectId":"0x26c26d20be986dbd99d33ff1b3a0bb16437039d3de85b4cb9d56a3f57066ef54","version":1978808,"digest":"3do9tdbdBCUMCn5rNweGt2Ag41fjubwr4WPkyMCiR6zv"}]},"inputs":[{"kind":"Input","value":{"Object":{"ImmOrOwned":{"digest":"DVy7dp7u9UWgh1JA1EbRfALxYPL9qRJhSXC5PZguWRMy","version":1978808,"objectId":"0xa7917d2b15ec1660b8a6658c29f99e5ba6feeb0cde66c7c5d072bd69b8574e30"}}},"index":0,"type":"object"},{"kind":"Input","value":{"Object":{"ImmOrOwned":{"digest":"GyKYSumghy7yGij9bKcf3wWL7f2M5nSLBiiE82TitmkT","version":1978808,"objectId":"0xd228b5a3f214a9a345399e8948a9844e5a0a02bbda221c2e5ba8b63ede9390d4"}}},"index":1,"type":"object"},{"kind":"Input","value":{"Object":{"ImmOrOwned":{"digest":"H4LGpfZMZuLWdiDdkJRBgT1pSRrcEkLyfWnkCwYtayLG","version":1978808,"objectId":"0xa21833a375925ff6352db9b1ef547a8da9cfcb717f60efa3edd243a1e3cee1b0"}}},"index":2,"type":"object"}],"transactions":[{"kind":"MergeCoins","destination":{"kind":"Input","value":"0xa21833a375925ff6352db9b1ef547a8da9cfcb717f60efa3edd243a1e3cee1b0","index":2,"type":"object"},"sources":[{"kind":"Input","value":"0xa7917d2b15ec1660b8a6658c29f99e5ba6feeb0cde66c7c5d072bd69b8574e30","index":0,"type":"object"},{"kind":"Input","value":"0xd228b5a3f214a9a345399e8948a9844e5a0a02bbda221c2e5ba8b63ede9390d4","index":1,"type":"object"}]}]}`

	mulJSon = `{"version":1,"sender":"0x19a8a5cadebbe8f73364cb1ce4d9aa2d0ffb62a3098c7ff93bd0426960f56406","gasConfig":{"price":"815","budget":"9534000","payment":[{"objectId":"0x26c26d20be986dbd99d33ff1b3a0bb16437039d3de85b4cb9d56a3f57066ef54","version":1978814,"digest":"8neVLSnZEDGjWq5ynYMP7sbCT7Tr1bLXnDTjw9dYuX58"}]},"inputs":[{"kind":"Input","value":{"Pure":[1,0,0,0,0,0,0,0]},"index":0,"type":"pure"},{"kind":"Input","value":{"Pure":[2,0,0,0,0,0,0,0]},"index":1,"type":"pure"},{"kind":"Input","value":{"Pure":[25,168,165,202,222,187,232,247,51,100,203,28,228,217,170,45,15,251,98,163,9,140,127,249,59,208,66,105,96,245,100,6]},"index":2,"type":"pure"},{"kind":"Input","value":{"Pure":[33,93,58,103,217,81,235,213,180,83,180,64,73,121,23,181,250,194,137,15,199,241,131,88,50,45,55,46,47,19,4,93]},"index":3,"type":"pure"}],"transactions":[{"kind":"SplitCoins","coin":{"kind":"GasCoin"},"amounts":[{"kind":"Input","value":1,"index":0,"type":"pure"},{"kind":"Input","value":2,"index":1,"type":"pure"}]},{"kind":"TransferObjects","objects":[{"kind":"NestedResult","index":0,"resultIndex":0}],"address":{"kind":"Input","value":"0x19a8a5cadebbe8f73364cb1ce4d9aa2d0ffb62a3098c7ff93bd0426960f56406","index":2,"type":"pure"}},{"kind":"TransferObjects","objects":[{"kind":"NestedResult","index":0,"resultIndex":1}],"address":{"kind":"Input","value":"0x19a8a5cadebbe8f73364cb1ce4d9aa2d0ffb62a3098c7ff93bd0426960f56406","index":3,"type":"pure"}}]}`
)

func toJson(v interface{}) string {
	b, err := json.MarshalIndent(v, "", " ")
	if err != nil {
		panic(err)
	}
	return string(b)
}

func TestNewAddress(t *testing.T) {
	b, err := base64.StdEncoding.DecodeString(seed)
	require.NoError(t, err)
	addr := NewAddress(hex.EncodeToString(b[0:32]))
	require.Equal(t, sender, addr)
	require.True(t, ValidateAddress(addr))
}

func TestDecodeHexString(t *testing.T) {
	b, err := DecodeHexString("0x")
	assert.Equal(t, err, nil)
	assert.NotEqual(t, b, nil)
	b, err = DecodeHexString("0x2")
	require.NoError(t, err)
	if len(b) > 32 {
		t.Fatal("invalid address")
	}
	if len(b) < 32 {
		bb := make([]byte, 32)
		copy(bb[32-len(b):], b)
		b = bb
	}
	require.Equal(t, len(b), 32)
}

func TestHash(t *testing.T) {
	hash, err := Hash("AAACAAgKAAAAAAAAAAAgIV06Z9lR69W0U7RASXkXtfrCiQ/H8YNYMi03Li8TBF0CAgABAQAAAQECAAABAQDVe00TKfTZWRmf5edRrrAAsoV3HkGalEtkPXnPpqdKVQEsj/ETTRUUJzPWmMYTJDqFMN7S+XJMNQeOa3wm21wUxfZMFgAAAAAAIGqg947MNtrnD2nvhiKwnhEfZPtbc3gy0FwOc2Aj8brx1XtNEyn02VkZn+XnUa6wALKFdx5BmpRLZD15z6anSlXoAwAAAAAAABCQLQAAAAAAAA==")
	require.NoError(t, err)
	require.Equal(t, "9Gm79kmn9XPVKVCqAJBd5MsmCKvLmGb3JamKyZcPBQRh", hash)
}

func TestSignMessage(t *testing.T) {
	b, err := base64.StdEncoding.DecodeString(seed)
	assert.NoError(t, err)
	seedHex := hex.EncodeToString(b[0:32])
	pubKey, err := ed25519.PublicKeyFromSeed(seedHex)
	assert.NoError(t, err)

	message := "im from okx"
	signature, err := SignMessage(message, seedHex)
	assert.NoError(t, err)
	assert.Equal(t, "AHhQeP0IE0F0fXfrbUSIEZPD974IR9Yx6DKBjPC5dP9kKYK9LwSg2EQyIvn/+iN3KEzipByemwM0Dq/F786tEwmXPl+82/o6Rt93H1ojjfMbkpJDm+Rnx1AAjN7Nvi7fnQ==", signature)

	err = VerifyMessage(message, base64.StdEncoding.EncodeToString(pubKey), signature)
	assert.NoError(t, err)

	hash, err := hex.DecodeString("ddb521e9f8756257e16cbb657feb022ba4c270939990e3bf0194e1330be44082")
	assert.NoError(t, err)
	sign, err := base64.StdEncoding.DecodeString(signature)
	assert.NoError(t, err)
	err = VerifySign(pubKey, sign, hash)
	assert.NoError(t, err)
}

func TestExecute(t *testing.T) {
	b, err := base64.StdEncoding.DecodeString(seed)
	require.NoError(t, err)
	addr := NewAddress(hex.EncodeToString(b[0:32]))
	pay := &PaySuiRequest{suiObjects, 1, 0}
	raw, err := json.Marshal(pay)
	require.NoError(t, err)
	res, err := Execute(&Request{Data: string(raw)}, addr, recipient, gasBudget, gasPrice, hex.EncodeToString(b[0:32]))
	require.NoError(t, err)
	var sx SignedTransaction
	err = json.Unmarshal([]byte(res), &sx)
	require.NoError(t, err)
	if sx.TxBytes != expectTxBytes || sx.Signature != expectSignature {
		buf1, _ := base64.StdEncoding.DecodeString(res)
		buf2, _ := base64.StdEncoding.DecodeString(expectTxBytes)
		t.Log("buf1 : ", buf1)
		t.Log("buf2 : ", buf2)
		t.Fatal(" TransferSui fail")
	}
}

func TestGenerate(t *testing.T) {
	k, err := GenerateKey()
	require.NoError(t, err)
	addr := NewAddress(hex.EncodeToString(k.Seed()))
	b, err := base64.StdEncoding.DecodeString(base64.StdEncoding.EncodeToString(k.Seed()))
	require.NoError(t, err)
	addr2 := NewAddress(hex.EncodeToString(b[0:32]))
	assert.Equal(t, addr2, addr)
}

func TestExecuteToken(t *testing.T) {
	b, err := base64.StdEncoding.DecodeString("uemYAwkvsf/a7q2DdoMKNHWP7DlDhLLmgUh6coTtp94=")
	require.NoError(t, err)
	pay := Pay{}
	err = json.Unmarshal([]byte(txJson2), &pay)
	require.NoError(t, err)
	data, err := pay.Build()
	require.NoError(t, err)
	tx, err := SignTransaction(base64.StdEncoding.EncodeToString(data), hex.EncodeToString(b[0:32]))
	require.NoError(t, err)
	assert.Equal(t, "AAADAQBbgAyeCnNRIkTQ+JvtrSyMGj/fZ5nRCtSam8iRZ828q6oxHgAAAAAAIJQq00aqPos3CP2pfPJoCO7VA4whbkbYcgECudT4zZnbAAgBAAAAAAAAAAAgGailyt676PczZMsc5NmqLQ/7YqMJjH/5O9BCaWD1ZAYCAgEAAAEBAQABAQIAAAECABmopcreu+j3M2TLHOTZqi0P+2KjCYx/+TvQQmlg9WQGASbCbSC+mG29mdM/8bOguxZDcDnT3oW0y51Wo/VwZu9UqjEeAAAAAAAgruCeRMC8oChpzsEOEdHLoBl1f12SLj/qPn7rYFx5PlAZqKXK3rvo9zNkyxzk2aotD/tiowmMf/k70EJpYPVkBugDAAAAAAAAAOH1BQAAAAAA", tx.TxBytes)
	assert.Equal(t, "AIz5cKG2KskS4trV1c8O3WYtKnC2aOWyuQjlwTlSB7ODpwajlJvg8kJaGjqcTzWi1YsjBNAx0owheR30UukPoAWXPl+82/o6Rt93H1ojjfMbkpJDm+Rnx1AAjN7Nvi7fnQ==", tx.Signature)
}

func TestExecuteToken2(t *testing.T) {
	b, err := base64.StdEncoding.DecodeString("uemYAwkvsf/a7q2DdoMKNHWP7DlDhLLmgUh6coTtp94=")
	require.NoError(t, err)
	pay := Pay{}
	err = json.Unmarshal([]byte(txJson3), &pay)
	require.NoError(t, err)
	data, err := pay.Build()
	require.NoError(t, err)
	tx, err := SignTransaction(base64.StdEncoding.EncodeToString(data), hex.EncodeToString(b[0:32]))
	require.NoError(t, err)
	require.Equal(t, tx.TxBytes, "AAAEAQCNserIE7MB8NZYXdnUcNavwBjnexGs/ZCJKzvjCDt05K0xHgAAAAAAIKOamsmHCCi+Pa1UmQcFnHhpDUaojMJ73ygqGAIVQzjgAQBDAtBr+uNwJLUj+uhn0Cyml/LDdO+k05PecARKHwTIGq0xHgAAAAAAIF2mp5EgLMSb5OAysoBRvnYiquv3oe86JD7nwVricRdfAAgBAAAAAAAAAAAgGailyt676PczZMsc5NmqLQ/7YqMJjH/5O9BCaWD1ZAYDAwEBAAEBAAACAQEAAQECAAEBAgEAAQMAGailyt676PczZMsc5NmqLQ/7YqMJjH/5O9BCaWD1ZAYBJsJtIL6Ybb2Z0z/xs6C7FkNwOdPehbTLnVaj9XBm71StMR4AAAAAACD0RCcdA6WA0cHdn1N+8S6RUwoeNzqMlBgf335aJ9520Bmopcreu+j3M2TLHOTZqi0P+2KjCYx/+TvQQmlg9WQG6AMAAAAAAADAxi0AAAAAAAA=")
	require.Equal(t, tx.Signature, "AHlnpxG13z0pfjKaq2pljgn1VCVIi/UhKfL8KBfciQ+iSLbEEZGOp7JnKVxsf1lQejGc4KBqwTonv2W/pV5T8gSXPl+82/o6Rt93H1ojjfMbkpJDm+Rnx1AAjN7Nvi7fnQ==")
}

func TestEqual2(t *testing.T) {
	pay := Pay{}
	err := json.Unmarshal([]byte(txJson3), &pay)
	require.NoError(t, err)
	data, err := pay.Build()
	require.NoError(t, err)

	coins := []*SuiObjectRef{
		{ObjectId: "0x26c26d20be986dbd99d33ff1b3a0bb16437039d3de85b4cb9d56a3f57066ef54", Digest: "HSWkfqRq34bR4gS89xDbMBjSxANhpckJyxwdQyQLptUX", Version: 1978797},
	}
	tokens := []*SuiObjectRef{
		{ObjectId: "0x8db1eac813b301f0d6585dd9d470d6afc018e77b11acfd90892b3be3083b74e4", Digest: "C1eBtxoHXE7fgdUbaqxrk2RCTR1XQ9BdreWbjPxwbPQ3", Version: 1978797},
		{ObjectId: "0x4302d06bfae37024b523fae867d02ca697f2c374efa4d393de70044a1f04c81a", Digest: "7JaJkUiU5TKEKo4CtDbnVVUJMSxNiLQbJwFGmQBfwZze", Version: 1978797},
	}
	data2, err := BuildTokenTx("0x19a8a5cadebbe8f73364cb1ce4d9aa2d0ffb62a3098c7ff93bd0426960f56406", "0x19a8a5cadebbe8f73364cb1ce4d9aa2d0ffb62a3098c7ff93bd0426960f56406", coins, tokens, 1, 0, 3000000, 1000)
	require.NoError(t, err)
	require.Equal(t, data2, data)
}

func TestEqualToken2(t *testing.T) {
	pay := Pay{}
	err := json.Unmarshal([]byte(txJson2), &pay)
	require.NoError(t, err)
	data, err := pay.Build()
	require.NoError(t, err)
	coins := []*SuiObjectRef{{ObjectId: "0x26c26d20be986dbd99d33ff1b3a0bb16437039d3de85b4cb9d56a3f57066ef54", Digest: "CmebPaLJ6ggh6sTiDpfu5aSVZHHMGedyQCLQ7dwUFaPH", Version: 1978794}}
	tokens := []*SuiObjectRef{{ObjectId: "0x5b800c9e0a73512244d0f89bedad2c8c1a3fdf6799d10ad49a9bc89167cdbcab", Digest: "AyPDRVnTtwnw9tsTDcSpojN7hn8TiVr7S1SzuxD8b3CN", Version: 1978794}}
	data2, err := BuildTokenTx("0x19a8a5cadebbe8f73364cb1ce4d9aa2d0ffb62a3098c7ff93bd0426960f56406", "0x19a8a5cadebbe8f73364cb1ce4d9aa2d0ffb62a3098c7ff93bd0426960f56406", coins, tokens, 1, 0, 100000000, 1000)
	require.NoError(t, err)
	require.Equal(t, data2, data)
}

func TestExecuteToken4(t *testing.T) {
	key := "b9e99803092fb1ffdaeead8376830a34758fec394384b2e681487a7284eda7de"
	pay := Pay{}
	err := json.Unmarshal([]byte(txJson4), &pay)
	require.NoError(t, err)
	data, err := pay.Build()
	require.NoError(t, err)
	tx, err := SignTransaction(base64.StdEncoding.EncodeToString(data), key)
	require.NoError(t, err)
	require.Equal(t, tx.TxBytes, "AAAGAQBVjx0EfUYhSx/bdWzRwTTVlggzxU3HGZcrZfMJkx30zbIxHgAAAAAAIP8F6dZtD3aV6nZnA2ncgk22qW2zVkI89P1VlWmNzjn0AQCYcy96g4gXTc1Yhz85qzJ6seaFSuVsd7nO4JMShJXtcLAxHgAAAAAAIBnl1ajZS5HU5b8W1Lo7yYdECFgSkTIfZWVtuEuvEYu2AQDmcpCns2w6dTpQuDOu7e6sukBY22qXDew4LDGwRIFP3a4xHgAAAAAAIMb1yKoUyEV3meUa0NhstjKR44mtbxN8Gq4Y0TtuAn+FAQBDAtBr+uNwJLUj+uhn0Cyml/LDdO+k05PecARKHwTIGq4xHgAAAAAAIMsU3nrSe89SLm/oZAXC2dptBpHpgbBvn2XoYu2lWCsZAAgBAAAAAAAAAAAgGailyt676PczZMsc5NmqLQ/7YqMJjH/5O9BCaWD1ZAYDAwEDAAMBAAABAQABAgACAQMAAQEEAAEBAgEAAQUAGailyt676PczZMsc5NmqLQ/7YqMJjH/5O9BCaWD1ZAYCJsJtIL6Ybb2Z0z/xs6C7FkNwOdPehbTLnVaj9XBm71SvMR4AAAAAACCJacVX9Ay3CEGPFf/g5xtvzALSn6uOxevghkr9+/ZxrudtiwH8JacDXEdIxwAPdUUNcmr0/o7wD2wfgAZlugRjsTEeAAAAAAAgYqgAcxgVmPx0gkoZALMSCL1hCGQczatrehuCcublyxAZqKXK3rvo9zNkyxzk2aotD/tiowmMf/k70EJpYPVkBugDAAAAAAAAwMYtAAAAAAAA")
	require.Equal(t, tx.Signature, "AL8xBYt1p/UwuBhNoHyxBFZ380hEPVCMrudIIw39cR6jQiieP8+zYLQ9qtLHuZi4JzrJa3HVExV/iu2qRIqpWwqXPl+82/o6Rt93H1ojjfMbkpJDm+Rnx1AAjN7Nvi7fnQ==")

}

func TestExecuteSplit2(t *testing.T) {
	key := "b9e99803092fb1ffdaeead8376830a34758fec394384b2e681487a7284eda7de"
	addr := NewAddress(key)
	split, err := BuildSplitTx(addr, addr, []*SuiObjectRef{{Digest: "8m7eNZYu9uZsbghPJy8gKxQq27iMPYF2uE1LHbBVJQyS", ObjectId: "0x26c26d20be986dbd99d33ff1b3a0bb16437039d3de85b4cb9d56a3f57066ef54", Version: 1978809}}, []uint64{1, 2, 3, 4, 5, 6, 7}, 0, 5582000, 815)
	require.NoError(t, err)
	data := toJson(&SplitSuiRequest{Coins: []*SuiObjectRef{{Digest: "8m7eNZYu9uZsbghPJy8gKxQq27iMPYF2uE1LHbBVJQyS", ObjectId: "0x26c26d20be986dbd99d33ff1b3a0bb16437039d3de85b4cb9d56a3f57066ef54", Version: 1978809}}, Amounts: []uint64{1, 2, 3, 4, 5, 6, 7}})
	res, err := PrepareTx(&Request{Data: data, Type: Split}, addr, 5582000, 815, addr)
	require.NoError(t, err)
	require.Equal(t, base64.StdEncoding.EncodeToString(split), res)
}

func TestExecuteSplit(t *testing.T) {
	key := "b9e99803092fb1ffdaeead8376830a34758fec394384b2e681487a7284eda7de"
	addr := NewAddress(key)
	split, err := BuildSplitTx(addr, addr, []*SuiObjectRef{{Digest: "8m7eNZYu9uZsbghPJy8gKxQq27iMPYF2uE1LHbBVJQyS", ObjectId: "0x26c26d20be986dbd99d33ff1b3a0bb16437039d3de85b4cb9d56a3f57066ef54", Version: 1978809}}, []uint64{1, 2, 3, 4, 5, 6, 7}, 0, 5582000, 815)
	require.NoError(t, err)
	tx, err := SignTransaction(base64.StdEncoding.EncodeToString(split), key)
	require.NoError(t, err)
	pay := Pay{}
	err = json.Unmarshal([]byte(txJson5), &pay)
	require.NoError(t, err)
	data, err := pay.Build()
	require.NoError(t, err)
	assert.Equal(t, split, data)
	tx, err = SignTransaction(base64.StdEncoding.EncodeToString(data), key)
	require.NoError(t, err)
	require.Equal(t, tx.TxBytes, "AAAIAAgBAAAAAAAAAAAIAgAAAAAAAAAACAMAAAAAAAAAAAgEAAAAAAAAAAAIBQAAAAAAAAAACAYAAAAAAAAAAAgHAAAAAAAAAAAgGailyt676PczZMsc5NmqLQ/7YqMJjH/5O9BCaWD1ZAYCAgAHAQAAAQEAAQIAAQMAAQQAAQUAAQYAAQcDAAAAAAMAAAEAAwAAAgADAAADAAMAAAQAAwAABQADAAAGAAEHABmopcreu+j3M2TLHOTZqi0P+2KjCYx/+TvQQmlg9WQGASbCbSC+mG29mdM/8bOguxZDcDnT3oW0y51Wo/VwZu9UuTEeAAAAAAAgc07geh9mPnCD+GJ0yLZkf47CF6K7ApzKJC8B22LwkS0ZqKXK3rvo9zNkyxzk2aotD/tiowmMf/k70EJpYPVkBi8DAAAAAAAAsCxVAAAAAAAA")
	require.Equal(t, tx.Signature, "AI8b8D4XO51mpZ6LXmIGhTkZPRgV/GBH9pRrM+hWNVi8HXXiwOWKbd+C4cRlUzzr6/j6Qyv8PVeN6Aoi3APi3weXPl+82/o6Rt93H1ojjfMbkpJDm+Rnx1AAjN7Nvi7fnQ==")
}

func TestExecuteMerge2(t *testing.T) {
	key := "b9e99803092fb1ffdaeead8376830a34758fec394384b2e681487a7284eda7de"
	addr := NewAddress(key)
	merge, err := BuildMergeTx(addr, []*SuiObjectRef{{Digest: "3do9tdbdBCUMCn5rNweGt2Ag41fjubwr4WPkyMCiR6zv", ObjectId: "0x26c26d20be986dbd99d33ff1b3a0bb16437039d3de85b4cb9d56a3f57066ef54", Version: 1978808}},
		[]*SuiObjectRef{{Digest: "DVy7dp7u9UWgh1JA1EbRfALxYPL9qRJhSXC5PZguWRMy", ObjectId: "0xa7917d2b15ec1660b8a6658c29f99e5ba6feeb0cde66c7c5d072bd69b8574e30", Version: 1978808},
			{Digest: "GyKYSumghy7yGij9bKcf3wWL7f2M5nSLBiiE82TitmkT", ObjectId: "0xd228b5a3f214a9a345399e8948a9844e5a0a02bbda221c2e5ba8b63ede9390d4", Version: 1978808},
			{Digest: "H4LGpfZMZuLWdiDdkJRBgT1pSRrcEkLyfWnkCwYtayLG", ObjectId: "0xa21833a375925ff6352db9b1ef547a8da9cfcb717f60efa3edd243a1e3cee1b0", Version: 1978808}},
		0, 1630000, 815)
	require.NoError(t, err)
	data := toJson(&MergeSuiRequest{Coins: []*SuiObjectRef{{Digest: "3do9tdbdBCUMCn5rNweGt2Ag41fjubwr4WPkyMCiR6zv", ObjectId: "0x26c26d20be986dbd99d33ff1b3a0bb16437039d3de85b4cb9d56a3f57066ef54", Version: 1978808}},
		Objects: []*SuiObjectRef{{Digest: "DVy7dp7u9UWgh1JA1EbRfALxYPL9qRJhSXC5PZguWRMy", ObjectId: "0xa7917d2b15ec1660b8a6658c29f99e5ba6feeb0cde66c7c5d072bd69b8574e30", Version: 1978808},
			{Digest: "GyKYSumghy7yGij9bKcf3wWL7f2M5nSLBiiE82TitmkT", ObjectId: "0xd228b5a3f214a9a345399e8948a9844e5a0a02bbda221c2e5ba8b63ede9390d4", Version: 1978808},
			{Digest: "H4LGpfZMZuLWdiDdkJRBgT1pSRrcEkLyfWnkCwYtayLG", ObjectId: "0xa21833a375925ff6352db9b1ef547a8da9cfcb717f60efa3edd243a1e3cee1b0", Version: 1978808}}})
	res, err := PrepareTx(&Request{Data: data, Type: Merge}, addr, 1630000, 815, addr)
	require.NoError(t, err)
	assert.Equal(t, base64.StdEncoding.EncodeToString(merge), res)
}

func TestExecuteMerge(t *testing.T) {
	key := "b9e99803092fb1ffdaeead8376830a34758fec394384b2e681487a7284eda7de"
	addr := NewAddress(key)
	merge, err := BuildMergeTx(addr, []*SuiObjectRef{{Digest: "3do9tdbdBCUMCn5rNweGt2Ag41fjubwr4WPkyMCiR6zv", ObjectId: "0x26c26d20be986dbd99d33ff1b3a0bb16437039d3de85b4cb9d56a3f57066ef54", Version: 1978808}},
		[]*SuiObjectRef{{Digest: "DVy7dp7u9UWgh1JA1EbRfALxYPL9qRJhSXC5PZguWRMy", ObjectId: "0xa7917d2b15ec1660b8a6658c29f99e5ba6feeb0cde66c7c5d072bd69b8574e30", Version: 1978808},
			{Digest: "GyKYSumghy7yGij9bKcf3wWL7f2M5nSLBiiE82TitmkT", ObjectId: "0xd228b5a3f214a9a345399e8948a9844e5a0a02bbda221c2e5ba8b63ede9390d4", Version: 1978808},
			{Digest: "H4LGpfZMZuLWdiDdkJRBgT1pSRrcEkLyfWnkCwYtayLG", ObjectId: "0xa21833a375925ff6352db9b1ef547a8da9cfcb717f60efa3edd243a1e3cee1b0", Version: 1978808}},
		0, 1630000, 815)
	require.NoError(t, err)
	tx, err := SignTransaction(base64.StdEncoding.EncodeToString(merge), key)
	require.NoError(t, err)
	pay := Pay{}
	err = json.Unmarshal([]byte(mergeJson6), &pay)
	require.NoError(t, err)
	data, err := pay.Build()
	require.NoError(t, err)
	assert.Equal(t, merge, data)
	tx, err = SignTransaction(base64.StdEncoding.EncodeToString(data), key)
	require.NoError(t, err)
	require.Equal(t, tx.TxBytes, "AAADAQCnkX0rFewWYLimZYwp+Z5bpv7rDN5mx8XQcr1puFdOMLgxHgAAAAAAILm39FsWLIpL4vns+V6pEcpL0PN7VJ8ORuFfODlAlVpIAQDSKLWj8hSpo0U5nolIqYROWgoCu9oiHC5bqLY+3pOQ1LgxHgAAAAAAIO1My3/LXFaB95we+up5H3Luqoeyrs1yd2DnI6RIrraQAQCiGDOjdZJf9jUtubHvVHqNqc/LcX9g76Pt0kOh487hsLgxHgAAAAAAIO6VhmpyEnX3NyXrpnzYsmhaT7No/umICWshdYQTiyH1AQMBAgACAQAAAQEAGailyt676PczZMsc5NmqLQ/7YqMJjH/5O9BCaWD1ZAYBJsJtIL6Ybb2Z0z/xs6C7FkNwOdPehbTLnVaj9XBm71S4MR4AAAAAACAnJHf6AKPRorbjT2n++V1UWg5HVnWH6V/NoIeD7ff0Axmopcreu+j3M2TLHOTZqi0P+2KjCYx/+TvQQmlg9WQGLwMAAAAAAAAw3xgAAAAAAAA=")
	require.Equal(t, tx.Signature, "AETPBLNuaCSDYuzJ86hVdIbBrfYZgMfVc4cHo8ievO9RGj6UdcXrlkW041MqRbW7RdvCTJ07DJUMyuMjEp9fMQaXPl+82/o6Rt93H1ojjfMbkpJDm+Rnx1AAjN7Nvi7fnQ==")
}

func TestExecuteMul(t *testing.T) {
	key := "b9e99803092fb1ffdaeead8376830a34758fec394384b2e681487a7284eda7de"
	addr := NewAddress(key)
	mul, err := BuildMulTx(addr, []*SuiObjectRef{{Digest: "8neVLSnZEDGjWq5ynYMP7sbCT7Tr1bLXnDTjw9dYuX58", ObjectId: "0x26c26d20be986dbd99d33ff1b3a0bb16437039d3de85b4cb9d56a3f57066ef54", Version: 1978814}}, map[string]uint64{"0x19a8a5cadebbe8f73364cb1ce4d9aa2d0ffb62a3098c7ff93bd0426960f56406": 1, "0x215d3a67d951ebd5b453b440497917b5fac2890fc7f18358322d372e2f13045d": 2}, 0, 9534000, 815)
	require.NoError(t, err)
	tx, err := SignTransaction(base64.StdEncoding.EncodeToString(mul), key)
	require.NoError(t, err)
	pay := Pay{}
	err = json.Unmarshal([]byte(mulJSon), &pay)
	require.NoError(t, err)
	data, err := pay.Build()
	require.NoError(t, err)
	tx, err = SignTransaction(base64.StdEncoding.EncodeToString(data), key)
	require.NoError(t, err)

	require.Equal(t, tx.TxBytes, "AAAEAAgBAAAAAAAAAAAIAgAAAAAAAAAAIBmopcreu+j3M2TLHOTZqi0P+2KjCYx/+TvQQmlg9WQGACAhXTpn2VHr1bRTtEBJeRe1+sKJD8fxg1gyLTcuLxMEXQMCAAIBAAABAQABAQMAAAAAAQIAAQEDAAABAAEDABmopcreu+j3M2TLHOTZqi0P+2KjCYx/+TvQQmlg9WQGASbCbSC+mG29mdM/8bOguxZDcDnT3oW0y51Wo/VwZu9UvjEeAAAAAAAgc7NVVGHTJpRDCtFbtTFvU3yDOVEYvLbkomev5JKeybcZqKXK3rvo9zNkyxzk2aotD/tiowmMf/k70EJpYPVkBi8DAAAAAAAAMHqRAAAAAAAA")
	require.Equal(t, tx.Signature, "AJzm3G/4aUioFDmaP32dUboklrbG9TDp+dvbgMWtuEiKLEZti/psv4rIl6wDFdR91AoiNyrdl0VuNmkZQCWNzQWXPl+82/o6Rt93H1ojjfMbkpJDm+Rnx1AAjN7Nvi7fnQ==")
}

func TestEqualToken4(t *testing.T) {
	pay1 := Pay{}
	if err := json.Unmarshal([]byte(txJson4), &pay1); err != nil {
		t.Fatal(err)
	}
	data, err := pay1.Build()
	require.NoError(t, err)
	coins := []*SuiObjectRef{
		{ObjectId: "0x26c26d20be986dbd99d33ff1b3a0bb16437039d3de85b4cb9d56a3f57066ef54", Digest: "AFQQ8sCHFQcoE9yBCFnhCojeieyiDexvonV2ojpcSZYq", Version: 1978799},
		{ObjectId: "0xe76d8b01fc25a7035c4748c7000f75450d726af4fe8ef00f6c1f800665ba0463", Digest: "7e7Y64CAiTZdjLBnz4Ncytt7t9V2XdWGbp4XGSYu1tRM", Version: 1978801},
	}
	tokens := []*SuiObjectRef{
		{ObjectId: "0x558f1d047d46214b1fdb756cd1c134d5960833c54dc719972b65f309931df4cd", Digest: "JAWCBZc7V6z4kBRBJCSc7qBoVvnXvkJrpPWPwFowQBBH", Version: 1978802},
		{ObjectId: "0x98732f7a8388174dcd58873f39ab327ab1e6854ae56c77b9cee093128495ed70", Digest: "2k6TeD147XXphLzRrx1LRBTckeqtxgLyM8tpu5vmSKLR", Version: 1978800},
		{ObjectId: "0xe67290a7b36c3a753a50b833aeedeeacba4058db6a970dec382c31b044814fdd", Digest: "EPf6c5jNmQSZYtJ8rP2qMkvjz5fnK8JggapwjRCGDTYg", Version: 1978798},
		{ObjectId: "0x4302d06bfae37024b523fae867d02ca697f2c374efa4d393de70044a1f04c81a", Digest: "EfkDithrS2xWznsiz5GqAFFKAeSoHLei32ocGFte3NfA", Version: 1978798},
	}
	data2, err := BuildTokenTx("0x19a8a5cadebbe8f73364cb1ce4d9aa2d0ffb62a3098c7ff93bd0426960f56406", "0x19a8a5cadebbe8f73364cb1ce4d9aa2d0ffb62a3098c7ff93bd0426960f56406", coins, tokens, 1, 0, 3000000, 1000)
	require.NoError(t, err)
	if !bytes.Equal(data2, data) {
		for i := 0; i < len(data) && i < len(data2); i++ {
			if data[i] != data2[i] {
				t.Log("diff : ", i, data2[i], data[i])
			}
		}
		t.Fatal()
	}
}

func TestStake(t *testing.T) {
	data, err := BuildStakeTx("0x19a8a5cadebbe8f73364cb1ce4d9aa2d0ffb62a3098c7ff93bd0426960f56406", "0x72169c90b7ea87f8101285c849c09cacced9968f83aa30786dad546bb94c78ab",
		[]*SuiObjectRef{{ObjectId: "0x26c26d20be986dbd99d33ff1b3a0bb16437039d3de85b4cb9d56a3f57066ef54", Digest: "AMGM65x2qTfM4kfPjbv7Aqpap6MBiVVa4W8hrakgvPjB", Version: 1978816}},
		1000000000, 0, 9644512, 820)
	require.NoError(t, err)
	key := "b9e99803092fb1ffdaeead8376830a34758fec394384b2e681487a7284eda7de"
	tx, err := SignTransaction(base64.StdEncoding.EncodeToString(data), key)
	require.NoError(t, err)
	require.Equal(t, tx.TxBytes, "AAADAAgAypo7AAAAAAEBAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAUBAAAAAAAAAAEAIHIWnJC36of4EBKFyEnAnKzO2ZaPg6oweG2tVGu5THirAgIAAQEAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAwpzdWlfc3lzdGVtEXJlcXVlc3RfYWRkX3N0YWtlAAMBAQACAAABAgAZqKXK3rvo9zNkyxzk2aotD/tiowmMf/k70EJpYPVkBgEmwm0gvphtvZnTP/GzoLsWQ3A5096FtMudVqP1cGbvVMAxHgAAAAAAIIrqJpLzvTYHdoUCDU1KmmXxL+/TEZdmF0A8w9aqCH7+Gailyt676PczZMsc5NmqLQ/7YqMJjH/5O9BCaWD1ZAY0AwAAAAAAAOApkwAAAAAAAA==")
	require.Equal(t, tx.Signature, "AMyhKkgz6zoV22pl+JGrPUqe0ttQWPgArpG4oB/j+sKUP9rU6IkVpzCvH+RYiHp8CvBGBgmkC4YIVWA0xKBSBQOXPl+82/o6Rt93H1ojjfMbkpJDm+Rnx1AAjN7Nvi7fnQ==")
}

func TestWithdraw(t *testing.T) {
	data, err := BuildWithdrawStakeTx("0x19a8a5cadebbe8f73364cb1ce4d9aa2d0ffb62a3098c7ff93bd0426960f56406",
		[]*SuiObjectRef{{ObjectId: "0x26c26d20be986dbd99d33ff1b3a0bb16437039d3de85b4cb9d56a3f57066ef54", Digest: "AihBh2VjG96NDTCw1HvZj8TtWEpmYZUh9rb1D92oQ7Ak", Version: 5656730}},
		&SuiObjectRef{Digest: "CkmUVCkHFWyjH27Zg5xTd5xbUQZt1BReQnvtq2zeT6zW", Version: 5656730, ObjectId: "0x194acb4ec803ef63f15331efa9e701b4a334cf417fa15432d736d90978ce43e4"}, 0, 9534000, 820)
	require.NoError(t, err)
	key := "b9e99803092fb1ffdaeead8376830a34758fec394384b2e681487a7284eda7de"
	tx, err := SignTransaction(base64.StdEncoding.EncodeToString(data), key)
	require.NoError(t, err)
	require.Equal(t, tx.TxBytes, "AAACAQEAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAABQEAAAAAAAAAAQEAGUrLTsgD72PxUzHvqecBtKM0z0F/oVQy1zbZCXjOQ+SaUFYAAAAAACCuptFaCeiByBnx2kjvoGQSG8HVSQDZY4/GJYRmToFZ+wEAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAMKc3VpX3N5c3RlbRZyZXF1ZXN0X3dpdGhkcmF3X3N0YWtlAAIBAAABAQAZqKXK3rvo9zNkyxzk2aotD/tiowmMf/k70EJpYPVkBgEmwm0gvphtvZnTP/GzoLsWQ3A5096FtMudVqP1cGbvVJpQVgAAAAAAIJBnbps0rg7WV3qFRlezNq7ebV0Hry0IL/LATrprMgElGailyt676PczZMsc5NmqLQ/7YqMJjH/5O9BCaWD1ZAY0AwAAAAAAADB6kQAAAAAAAA==")
	require.Equal(t, tx.Signature, "AEG1Da49kGDx/ngBD3goKO28d58g3sTkV4z6KCAm7uk8Bv1j5MUlr92wZM6KO1eKsNk5qpupu1gzh2UcADAIpgCXPl+82/o6Rt93H1ojjfMbkpJDm+Rnx1AAjN7Nvi7fnQ==")
}

func TestEqualToken5(t *testing.T) {
	pay := Pay{}
	err := json.Unmarshal([]byte(txJson4), &pay)
	require.NoError(t, err)
	coins := []*SuiObjectRef{
		{ObjectId: "0x26c26d20be986dbd99d33ff1b3a0bb16437039d3de85b4cb9d56a3f57066ef54", Digest: "AFQQ8sCHFQcoE9yBCFnhCojeieyiDexvonV2ojpcSZYq", Version: 1978799},
		{ObjectId: "0xe76d8b01fc25a7035c4748c7000f75450d726af4fe8ef00f6c1f800665ba0463", Digest: "7e7Y64CAiTZdjLBnz4Ncytt7t9V2XdWGbp4XGSYu1tRM", Version: 1978801},
	}
	tokens := []*SuiObjectRef{
		{ObjectId: "0x558f1d047d46214b1fdb756cd1c134d5960833c54dc719972b65f309931df4cd", Digest: "JAWCBZc7V6z4kBRBJCSc7qBoVvnXvkJrpPWPwFowQBBH", Version: 1978802},
		{ObjectId: "0x98732f7a8388174dcd58873f39ab327ab1e6854ae56c77b9cee093128495ed70", Digest: "2k6TeD147XXphLzRrx1LRBTckeqtxgLyM8tpu5vmSKLR", Version: 1978800},
		{ObjectId: "0xe67290a7b36c3a753a50b833aeedeeacba4058db6a970dec382c31b044814fdd", Digest: "EPf6c5jNmQSZYtJ8rP2qMkvjz5fnK8JggapwjRCGDTYg", Version: 1978798},
		{ObjectId: "0x4302d06bfae37024b523fae867d02ca697f2c374efa4d393de70044a1f04c81a", Digest: "EfkDithrS2xWznsiz5GqAFFKAeSoHLei32ocGFte3NfA", Version: 1978798},
	}
	data, err := BuildTokenTx("0x19a8a5cadebbe8f73364cb1ce4d9aa2d0ffb62a3098c7ff93bd0426960f56406", "0x19a8a5cadebbe8f73364cb1ce4d9aa2d0ffb62a3098c7ff93bd0426960f56406", coins, tokens, 1, 0, 3000000, 1000)
	require.NoError(t, err)
	expected := "0000060100558f1d047d46214b1fdb756cd1c134d5960833c54dc719972b65f309931df4cdb2311e000000000020ff05e9d66d0f7695ea76670369dc824db6a96db356423cf4fd5595698dce39f4010098732f7a8388174dcd58873f39ab327ab1e6854ae56c77b9cee093128495ed70b0311e00000000002019e5d5a8d94b91d4e5bf16d4ba3bc9874408581291321f65656db84baf118bb60100e67290a7b36c3a753a50b833aeedeeacba4058db6a970dec382c31b044814fddae311e000000000020c6f5c8aa14c8457799e51ad0d86cb63291e389ad6f137c1aae18d13b6e027f8501004302d06bfae37024b523fae867d02ca697f2c374efa4d393de70044a1f04c81aae311e000000000020cb14de7ad27bcf522e6fe86405c2d9da6d0691e981b06f9f65e862eda5582b1900080100000000000000002019a8a5cadebbe8f73364cb1ce4d9aa2d0ffb62a3098c7ff93bd0426960f564060303010300030100000101000102000201030001010400010102010001050019a8a5cadebbe8f73364cb1ce4d9aa2d0ffb62a3098c7ff93bd0426960f564060226c26d20be986dbd99d33ff1b3a0bb16437039d3de85b4cb9d56a3f57066ef54af311e0000000000208969c557f40cb708418f15ffe0e71b6fcc02d29fab8ec5ebe0864afdfbf671aee76d8b01fc25a7035c4748c7000f75450d726af4fe8ef00f6c1f800665ba0463b1311e00000000002062a80073181598fc74824a1900b31208bd6108641ccdab6b7a1b8272e6e5cb1019a8a5cadebbe8f73364cb1ce4d9aa2d0ffb62a3098c7ff93bd0426960f56406e803000000000000c0c62d000000000000"
	require.Equal(t, expected, hex.EncodeToString(data))
}

func TestGetAddressByPubKey(t *testing.T) {
	pri := "31342f041c5b54358074b4579231c8a300be65e687dff020bc7779598b42897a"

	addr1 := NewAddress(pri)
	p, err := ed25519.PublicKeyFromSeed(pri)
	require.NoError(t, err)
	pub := hex.EncodeToString(p)
	assert.Equal(t, "bcc5c3f2165a5ee30e836878d4205fd1023225a176977d1a84d6a422e9e27d2d", pub)

	addr2, err := GetAddressByPubKey("bcc5c3f2165a5ee30e836878d4205fd1023225a176977d1a84d6a422e9e27d2d")
	require.NoError(t, err)
	assert.Equal(t, addr1, addr2)
}
