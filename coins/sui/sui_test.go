package sui

import (
	"bytes"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"github.com/stretchr/testify/assert"
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

	seed            = "//todo please replace your base64 sui key"
	recipient       = "0x215d3a67d951ebd5b453b440497917b5fac2890fc7f18358322d372e2f13045d"
	sender          = "0x11e0c681404673ddb16c929a852b352b9386a74455f9638d4c69f4c57d7f5f1a"
	gasBudget       = uint64(2976000)
	gasPrice        = uint64(1000)
	expectTxBytes   = "AAACAAgBAAAAAAAAAAAgIV06Z9lR69W0U7RASXkXtfrCiQ/H8YNYMi03Li8TBF0CAgABAQAAAQECAAABAQAR4MaBQEZz3bFskpqFKzUrk4anRFX5Y41MafTFfX9fGgMM0+gfITC5IqJfET8BfY92yujC+dfr7WkOVnVKCzpXhHDtAAAAAAAAIMeuvZGR3WQ3pQiWLbk7muGO5aih1NL8t9/oK2hFUZUnHYkjYQdCSQc/gmE60Is4e97CYYXqAT03n1rguyphHrxx7QAAAAAAACCGVO1kgYBA1g3K7CLQPUgcez0JMvsNV07m3F3E0ZG5WoVlqzrHByq9j3oNToGXSgyGad76rJNA3acIPmJULNL5b+0AAAAAAAAgF7dxv5oQMrjUfoOZZGMMbDIzaIewGNL0Ze1UOMWKshIR4MaBQEZz3bFskpqFKzUrk4anRFX5Y41MafTFfX9fGugDAAAAAAAAAGktAAAAAAAA"
	expectSignature = "AOZ3pwZYm3+RTYytY5pLOP4/0/ebx4y30qlJ8kW4fuYz88uUrcNHuplqEg5odQed31qpCy5StRen63cNgCjHeQK7Mod9GfbAH3tu3DBQ5omVbwy3xwqsFGXiRww2tmu/2w=="
	txJson          = `{"version":1,"sender":"0x11e0c681404673ddb16c929a852b352b9386a74455f9638d4c69f4c57d7f5f1a","expiration":{"Epoch":100},"gasConfig":{"price":"1","budget":"110","payment":[{"digest":"3Gei4443UwTNcGak4QXgeUvJjR7ZW7Qx8dPAURCx4hkG","objectId":"0xcde3e7ce68b62f5c0446b81946eeb443f62228ab121061afefbfac3fc231bfee","version":4}]},"inputs":[{"kind":"Input","value":100000,"index":0,"type":"pure"},{"kind":"Input","value":"0x215d3a67d951ebd5b453b440497917b5fac2890fc7f18358322d372e2f13045d","index":1,"type":"pure"}],"transactions":[{"kind":"SplitCoins","coin":{"kind":"GasCoin"},"amounts":[{"kind":"Input","value":100000,"index":0,"type":"pure"}]},{"kind":"TransferObjects","objects":[{"kind":"Result","index":0}],"address":{"kind":"Input","value":"0x215d3a67d951ebd5b453b440497917b5fac2890fc7f18358322d372e2f13045d","index":1,"type":"pure"}}]}`

	txJson2 = `{
	"version": 1,
	"sender": "0x1e7fa5fd46bdf8ec1291ca52084bdbeeabde6b3bab3a5d9e6cf83d78061de619",
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
			"Pure": [33, 93, 58, 103, 217, 81, 235, 213, 180, 83, 180, 64, 73, 121, 23, 181, 250, 194, 137, 15, 199, 241, 131, 88, 50, 45, 55, 46, 47, 19, 4, 93]
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
			"value": "0x215d3a67d951ebd5b453b440497917b5fac2890fc7f18358322d372e2f13045d",
			"index": 2,
			"type": "pure"
		}
	}]
}`

	// Multiple token objects, first merged, and then split the specified amount to each other. In this example, there are two tokens in the inputs.
	txJson3 = `{
	"version": 1,
	"sender": "0x1e7fa5fd46bdf8ec1291ca52084bdbeeabde6b3bab3a5d9e6cf83d78061de619",
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
			"Pure": [30, 127, 165, 253, 70, 189, 248, 236, 18, 145, 202, 82, 8, 75, 219, 238, 171, 222, 107, 59, 171, 58, 93, 158, 108, 248, 61, 120, 6, 29, 230, 25]
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
			"value": "0x1e7fa5fd46bdf8ec1291ca52084bdbeeabde6b3bab3a5d9e6cf83d78061de619",
			"index": 3,
			"type": "pure"
		}
	}]
}`

	// Multiple token objects, first merged, and then split the specified amount to each other. In this example, there are two tokens in the inputs..
	txJson4 = `{
	"version": 1,
	"sender": "0x1e7fa5fd46bdf8ec1291ca52084bdbeeabde6b3bab3a5d9e6cf83d78061de619",
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
			"Pure": [30, 127, 165, 253, 70, 189, 248, 236, 18, 145, 202, 82, 8, 75, 219, 238, 171, 222, 107, 59, 171, 58, 93, 158, 108, 248, 61, 120, 6, 29, 230, 25]
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
			"value": "0x1e7fa5fd46bdf8ec1291ca52084bdbeeabde6b3bab3a5d9e6cf83d78061de619",
			"index": 5,
			"type": "pure"
		}
	}]
}`
	txJson5 = `{"version":1,"sender":"0x1e7fa5fd46bdf8ec1291ca52084bdbeeabde6b3bab3a5d9e6cf83d78061de619","gasConfig":{"price":"815","budget":"5582000","payment":[{"objectId":"0x26c26d20be986dbd99d33ff1b3a0bb16437039d3de85b4cb9d56a3f57066ef54","version":1978809,"digest":"8m7eNZYu9uZsbghPJy8gKxQq27iMPYF2uE1LHbBVJQyS"}]},"inputs":[{"kind":"Input","value":{"Pure":[1,0,0,0,0,0,0,0]},"index":0,"type":"pure"},{"kind":"Input","value":{"Pure":[2,0,0,0,0,0,0,0]},"index":1,"type":"pure"},{"kind":"Input","value":{"Pure":[3,0,0,0,0,0,0,0]},"index":2,"type":"pure"},{"kind":"Input","value":{"Pure":[4,0,0,0,0,0,0,0]},"index":3,"type":"pure"},{"kind":"Input","value":{"Pure":[5,0,0,0,0,0,0,0]},"index":4,"type":"pure"},{"kind":"Input","value":{"Pure":[6,0,0,0,0,0,0,0]},"index":5,"type":"pure"},{"kind":"Input","value":{"Pure":[7,0,0,0,0,0,0,0]},"index":6,"type":"pure"},{"kind":"Input","value":{"Pure":[30,127,165,253,70,189,248,236,18,145,202,82,8,75,219,238,171,222,107,59,171,58,93,158,108,248,61,120,6,29,230,25]},"index":7,"type":"pure"}],"transactions":[{"kind":"SplitCoins","coin":{"kind":"GasCoin"},"amounts":[{"kind":"Input","value":1,"index":0,"type":"pure"},{"kind":"Input","value":2,"index":1,"type":"pure"},{"kind":"Input","value":3,"index":2,"type":"pure"},{"kind":"Input","value":4,"index":3,"type":"pure"},{"kind":"Input","value":5,"index":4,"type":"pure"},{"kind":"Input","value":6,"index":5,"type":"pure"},{"kind":"Input","value":7,"index":6,"type":"pure"}]},{"kind":"TransferObjects","objects":[{"kind":"NestedResult","index":0,"resultIndex":0},{"kind":"NestedResult","index":0,"resultIndex":1},{"kind":"NestedResult","index":0,"resultIndex":2},{"kind":"NestedResult","index":0,"resultIndex":3},{"kind":"NestedResult","index":0,"resultIndex":4},{"kind":"NestedResult","index":0,"resultIndex":5},{"kind":"NestedResult","index":0,"resultIndex":6}],"address":{"kind":"Input","value":"0x1e7fa5fd46bdf8ec1291ca52084bdbeeabde6b3bab3a5d9e6cf83d78061de619","index":7,"type":"pure"}}]}`

	mergeJson6 = `{"version":1,"sender":"0x1e7fa5fd46bdf8ec1291ca52084bdbeeabde6b3bab3a5d9e6cf83d78061de619","gasConfig":{"price":"815","budget":"1630000","payment":[{"objectId":"0x26c26d20be986dbd99d33ff1b3a0bb16437039d3de85b4cb9d56a3f57066ef54","version":1978808,"digest":"3do9tdbdBCUMCn5rNweGt2Ag41fjubwr4WPkyMCiR6zv"}]},"inputs":[{"kind":"Input","value":{"Object":{"ImmOrOwned":{"digest":"DVy7dp7u9UWgh1JA1EbRfALxYPL9qRJhSXC5PZguWRMy","version":1978808,"objectId":"0xa7917d2b15ec1660b8a6658c29f99e5ba6feeb0cde66c7c5d072bd69b8574e30"}}},"index":0,"type":"object"},{"kind":"Input","value":{"Object":{"ImmOrOwned":{"digest":"GyKYSumghy7yGij9bKcf3wWL7f2M5nSLBiiE82TitmkT","version":1978808,"objectId":"0xd228b5a3f214a9a345399e8948a9844e5a0a02bbda221c2e5ba8b63ede9390d4"}}},"index":1,"type":"object"},{"kind":"Input","value":{"Object":{"ImmOrOwned":{"digest":"H4LGpfZMZuLWdiDdkJRBgT1pSRrcEkLyfWnkCwYtayLG","version":1978808,"objectId":"0xa21833a375925ff6352db9b1ef547a8da9cfcb717f60efa3edd243a1e3cee1b0"}}},"index":2,"type":"object"}],"transactions":[{"kind":"MergeCoins","destination":{"kind":"Input","value":"0xa21833a375925ff6352db9b1ef547a8da9cfcb717f60efa3edd243a1e3cee1b0","index":2,"type":"object"},"sources":[{"kind":"Input","value":"0xa7917d2b15ec1660b8a6658c29f99e5ba6feeb0cde66c7c5d072bd69b8574e30","index":0,"type":"object"},{"kind":"Input","value":"0xd228b5a3f214a9a345399e8948a9844e5a0a02bbda221c2e5ba8b63ede9390d4","index":1,"type":"object"}]}]}`

	mulJSon = `{"version":1,"sender":"0x1e7fa5fd46bdf8ec1291ca52084bdbeeabde6b3bab3a5d9e6cf83d78061de619","gasConfig":{"price":"815","budget":"9534000","payment":[{"objectId":"0x26c26d20be986dbd99d33ff1b3a0bb16437039d3de85b4cb9d56a3f57066ef54","version":1978814,"digest":"8neVLSnZEDGjWq5ynYMP7sbCT7Tr1bLXnDTjw9dYuX58"}]},"inputs":[{"kind":"Input","value":{"Pure":[1,0,0,0,0,0,0,0]},"index":0,"type":"pure"},{"kind":"Input","value":{"Pure":[2,0,0,0,0,0,0,0]},"index":1,"type":"pure"},{"kind":"Input","value":{"Pure":[30,127,165,253,70,189,248,236,18,145,202,82,8,75,219,238,171,222,107,59,171,58,93,158,108,248,61,120,6,29,230,25]},"index":2,"type":"pure"},{"kind":"Input","value":{"Pure":[33,93,58,103,217,81,235,213,180,83,180,64,73,121,23,181,250,194,137,15,199,241,131,88,50,45,55,46,47,19,4,93]},"index":3,"type":"pure"}],"transactions":[{"kind":"SplitCoins","coin":{"kind":"GasCoin"},"amounts":[{"kind":"Input","value":1,"index":0,"type":"pure"},{"kind":"Input","value":2,"index":1,"type":"pure"}]},{"kind":"TransferObjects","objects":[{"kind":"NestedResult","index":0,"resultIndex":0}],"address":{"kind":"Input","value":"0x1e7fa5fd46bdf8ec1291ca52084bdbeeabde6b3bab3a5d9e6cf83d78061de619","index":2,"type":"pure"}},{"kind":"TransferObjects","objects":[{"kind":"NestedResult","index":0,"resultIndex":1}],"address":{"kind":"Input","value":"0x215d3a67d951ebd5b453b440497917b5fac2890fc7f18358322d372e2f13045d","index":3,"type":"pure"}}]}`
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
	if err != nil {
		t.Fatal(err)
	}
	addr := NewAddress(hex.EncodeToString(b[0:32]))
	fmt.Println("addr", addr)
	if addr != sender {
		t.Fatal("invalid address")
	}
	if !ValidateAddress(addr) {
		t.Fail()
	}
}

func TestDecodeHexString(t *testing.T) {
	b, err := DecodeHexString("0x")
	assert.Equal(t, err, nil)
	assert.NotEqual(t, b, nil)

	b, err = DecodeHexString("0x2")
	if err != nil {
		t.Fatal(err)
	}
	if len(b) > 32 {
		t.Fatal("invalid address")
	}
	if len(b) < 32 {
		bb := make([]byte, 32)
		copy(bb[32-len(b):], b)
		b = bb
	}
	assert.Equal(t, len(b), 32)
}

func TestHash(t *testing.T) {
	hash, err := Hash("AAACAAgKAAAAAAAAAAAgIV06Z9lR69W0U7RASXkXtfrCiQ/H8YNYMi03Li8TBF0CAgABAQAAAQECAAABAQDVe00TKfTZWRmf5edRrrAAsoV3HkGalEtkPXnPpqdKVQEsj/ETTRUUJzPWmMYTJDqFMN7S+XJMNQeOa3wm21wUxfZMFgAAAAAAIGqg947MNtrnD2nvhiKwnhEfZPtbc3gy0FwOc2Aj8brx1XtNEyn02VkZn+XnUa6wALKFdx5BmpRLZD15z6anSlXoAwAAAAAAABCQLQAAAAAAAA==")
	fmt.Println(hash, err)
	if err != nil {
		t.FailNow()
	}
	if hash != "9Gm79kmn9XPVKVCqAJBd5MsmCKvLmGb3JamKyZcPBQRh" {
		t.FailNow()
	}
}

func TestSignMessage(t *testing.T) {
	b, err := base64.StdEncoding.DecodeString(seed)
	if err != nil {
		t.Fatal(err)
	}
	r, err := SignMessage("im from okx", hex.EncodeToString(b[0:32]))
	if err != nil {
		t.Error(err)
		t.FailNow()
	}
	fmt.Println(r, err)
	assert.Equal(t, r, "AKSQ4JWYP2zFlwaoFc80qkwBSaVUNM+ZkzvFWg/YOk93Qwn64VoXzTENTNJZg3mC/II8SQzzTMTwUdi4IiCxVAO7Mod9GfbAH3tu3DBQ5omVbwy3xwqsFGXiRww2tmu/2w==")
}

// https://explorer.sui.io/txblock/6ZJAwnt3HU1NNviTEtWcXWwW2sMd46gEGGXvis8TK9ao?network=testnet
func TestExecute(t *testing.T) {
	b, err := base64.StdEncoding.DecodeString(seed)
	if err != nil {
		t.Fatal(err)
	}
	addr := NewAddress(hex.EncodeToString(b[0:32]))
	pay := &PaySuiRequest{suiObjects, 1, 0}
	raw, err := json.Marshal(pay)
	if err != nil {
		t.Fatal(err)
	}
	res, err := Execute(&Request{Data: string(raw)}, addr, recipient, gasBudget, gasPrice, hex.EncodeToString(b[0:32]))
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(res)
	var sx SignedTransaction
	if err := json.Unmarshal([]byte(res), &sx); err != nil {
		t.Fatal(err)
		return
	}
	if sx.TxBytes != expectTxBytes || sx.Signature != expectSignature {
		buf1, _ := base64.StdEncoding.DecodeString(res)
		buf2, _ := base64.StdEncoding.DecodeString(expectTxBytes)
		fmt.Println("result", buf1)
		fmt.Println("expectTxBytes", buf2)
		t.Fatal(" TransferSui fail")
	}
}

func TestExecuteToken(t *testing.T) {
	b, err := base64.StdEncoding.DecodeString("//todo please replace your base64 sui key")
	if err != nil {
		t.Fatal(err)
	}
	addr := NewAddress(hex.EncodeToString(b[0:32]))
	fmt.Println("address", addr)
	pay := Pay{}
	if err := json.Unmarshal([]byte(txJson2), &pay); err != nil {
		t.Fatal(err)
	}
	data, err := pay.Build()
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(base64.StdEncoding.EncodeToString(data))
	tx, err := SignTransaction(base64.StdEncoding.EncodeToString(data), hex.EncodeToString(b[0:32]))
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(fmt.Sprintf("\"%s\",[\"%s\"]", tx.TxBytes, tx.Signature))
}

// https://suiexplorer.com/txblock/NVi6Exztk3iBDd7bHLuxh1euwmTf5XMkDCVu4hf4fWA
func TestExecuteToken2(t *testing.T) {
	b, err := base64.StdEncoding.DecodeString("//todo please replace your base64 sui key")
	if err != nil {
		t.Fatal(err)
	}
	addr := NewAddress(hex.EncodeToString(b[0:32]))
	fmt.Println("address", addr)
	pay := Pay{}
	if err := json.Unmarshal([]byte(txJson3), &pay); err != nil {
		t.Fatal(err)
	}
	data, err := pay.Build()
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(base64.StdEncoding.EncodeToString(data))
	tx, err := SignTransaction(base64.StdEncoding.EncodeToString(data), hex.EncodeToString(b[0:32]))
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(fmt.Sprintf("\"%s\",[\"%s\"]", tx.TxBytes, tx.Signature))
	assert.Equal(t, tx.TxBytes, "AAAEAQCNserIE7MB8NZYXdnUcNavwBjnexGs/ZCJKzvjCDt05K0xHgAAAAAAIKOamsmHCCi+Pa1UmQcFnHhpDUaojMJ73ygqGAIVQzjgAQBDAtBr+uNwJLUj+uhn0Cyml/LDdO+k05PecARKHwTIGq0xHgAAAAAAIF2mp5EgLMSb5OAysoBRvnYiquv3oe86JD7nwVricRdfAAgBAAAAAAAAAAAgHn+l/Ua9+OwSkcpSCEvb7qveazurOl2ebPg9eAYd5hkDAwEBAAEBAAACAQEAAQECAAEBAgEAAQMAHn+l/Ua9+OwSkcpSCEvb7qveazurOl2ebPg9eAYd5hkBJsJtIL6Ybb2Z0z/xs6C7FkNwOdPehbTLnVaj9XBm71StMR4AAAAAACD0RCcdA6WA0cHdn1N+8S6RUwoeNzqMlBgf335aJ9520B5/pf1GvfjsEpHKUghL2+6r3ms7qzpdnmz4PXgGHeYZ6AMAAAAAAADAxi0AAAAAAAA=")
	assert.Equal(t, tx.Signature, "ADfId3fKflgPGOEb+Tcw1e8Igrj/o4wowxKEWXNIsHGmZSQLcxeeMpVexyf6mN/9AZnIc/NGvWQ+NZgn2fVLvQK8xcPyFlpe4w6DaHjUIF/RAjIloXaXfRqE1qQi6eJ9LQ==")
}

func TestEqual2(t *testing.T) {
	b, err := base64.StdEncoding.DecodeString("//todo please replace your base64 sui key")
	if err != nil {
		t.Fatal(err)
	}
	addr := NewAddress(hex.EncodeToString(b[0:32]))
	fmt.Println("address", addr)
	pay := Pay{}
	if err := json.Unmarshal([]byte(txJson3), &pay); err != nil {
		t.Fatal(err)
	}
	data, err := pay.Build()
	if err != nil {
		t.Fatal(err)
	}

	coins := []*SuiObjectRef{
		{ObjectId: "0x26c26d20be986dbd99d33ff1b3a0bb16437039d3de85b4cb9d56a3f57066ef54", Digest: "HSWkfqRq34bR4gS89xDbMBjSxANhpckJyxwdQyQLptUX", Version: 1978797},
	}
	tokens := []*SuiObjectRef{
		{ObjectId: "0x8db1eac813b301f0d6585dd9d470d6afc018e77b11acfd90892b3be3083b74e4", Digest: "C1eBtxoHXE7fgdUbaqxrk2RCTR1XQ9BdreWbjPxwbPQ3", Version: 1978797},
		{ObjectId: "0x4302d06bfae37024b523fae867d02ca697f2c374efa4d393de70044a1f04c81a", Digest: "7JaJkUiU5TKEKo4CtDbnVVUJMSxNiLQbJwFGmQBfwZze", Version: 1978797},
	}
	data2, err := BuildTokenTx("0x1e7fa5fd46bdf8ec1291ca52084bdbeeabde6b3bab3a5d9e6cf83d78061de619", "0x1e7fa5fd46bdf8ec1291ca52084bdbeeabde6b3bab3a5d9e6cf83d78061de619", coins, tokens, 1, 0, 3000000, 1000)
	if err != nil {
		t.Fatal(err)
	}
	assert.Equal(t, data2, data)
}

func TestEncode(t *testing.T) {
	i := []byte{30, 127, 165, 253, 70, 189, 248, 236, 18, 145, 202, 82, 8, 75, 219, 238, 171, 222, 107, 59, 171, 58, 93, 158, 108, 248, 61, 120, 6, 29, 230, 25}
	fmt.Println(hex.EncodeToString(i))
}

func TestEqualToken2(t *testing.T) {
	pay1 := Pay{}
	if err := json.Unmarshal([]byte(txJson2), &pay1); err != nil {
		t.Fatal(err)
	}
	data, err := pay1.Build()
	if err != nil {
		t.Fatal(err)
	}
	coins := []*SuiObjectRef{{ObjectId: "0x26c26d20be986dbd99d33ff1b3a0bb16437039d3de85b4cb9d56a3f57066ef54", Digest: "CmebPaLJ6ggh6sTiDpfu5aSVZHHMGedyQCLQ7dwUFaPH", Version: 1978794}}
	tokens := []*SuiObjectRef{{ObjectId: "0x5b800c9e0a73512244d0f89bedad2c8c1a3fdf6799d10ad49a9bc89167cdbcab", Digest: "AyPDRVnTtwnw9tsTDcSpojN7hn8TiVr7S1SzuxD8b3CN", Version: 1978794}}
	data2, err := BuildTokenTx("0x1e7fa5fd46bdf8ec1291ca52084bdbeeabde6b3bab3a5d9e6cf83d78061de619", "0x215d3a67d951ebd5b453b440497917b5fac2890fc7f18358322d372e2f13045d", coins, tokens, 1, 0, 100000000, 1000)
	if err != nil {
		t.Fatal(err)
	}
	assert.Equal(t, data2, data)
}

// https://suiexplorer.com/txblock/BUvGPZXkWoWLHJpd8eFhGiSuoLH2Sq444TzTg5LcaNCB
func TestExecuteToken4(t *testing.T) {
	key := "//todo please replace your hex sui key"
	addr := NewAddress(key)
	fmt.Println("address", addr)
	pay := Pay{}
	if err := json.Unmarshal([]byte(txJson4), &pay); err != nil {
		t.Fatal(err)
	}
	data, err := pay.Build()
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(base64.StdEncoding.EncodeToString(data))
	tx, err := SignTransaction(base64.StdEncoding.EncodeToString(data), key)
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(fmt.Sprintf("\"%s\",[\"%s\"]", tx.TxBytes, tx.Signature))

	assert.Equal(t, tx.TxBytes, "AAAGAQBVjx0EfUYhSx/bdWzRwTTVlggzxU3HGZcrZfMJkx30zbIxHgAAAAAAIP8F6dZtD3aV6nZnA2ncgk22qW2zVkI89P1VlWmNzjn0AQCYcy96g4gXTc1Yhz85qzJ6seaFSuVsd7nO4JMShJXtcLAxHgAAAAAAIBnl1ajZS5HU5b8W1Lo7yYdECFgSkTIfZWVtuEuvEYu2AQDmcpCns2w6dTpQuDOu7e6sukBY22qXDew4LDGwRIFP3a4xHgAAAAAAIMb1yKoUyEV3meUa0NhstjKR44mtbxN8Gq4Y0TtuAn+FAQBDAtBr+uNwJLUj+uhn0Cyml/LDdO+k05PecARKHwTIGq4xHgAAAAAAIMsU3nrSe89SLm/oZAXC2dptBpHpgbBvn2XoYu2lWCsZAAgBAAAAAAAAAAAgHn+l/Ua9+OwSkcpSCEvb7qveazurOl2ebPg9eAYd5hkDAwEDAAMBAAABAQABAgACAQMAAQEEAAEBAgEAAQUAHn+l/Ua9+OwSkcpSCEvb7qveazurOl2ebPg9eAYd5hkCJsJtIL6Ybb2Z0z/xs6C7FkNwOdPehbTLnVaj9XBm71SvMR4AAAAAACCJacVX9Ay3CEGPFf/g5xtvzALSn6uOxevghkr9+/ZxrudtiwH8JacDXEdIxwAPdUUNcmr0/o7wD2wfgAZlugRjsTEeAAAAAAAgYqgAcxgVmPx0gkoZALMSCL1hCGQczatrehuCcublyxAef6X9Rr347BKRylIIS9vuq95rO6s6XZ5s+D14Bh3mGegDAAAAAAAAwMYtAAAAAAAA")
	assert.Equal(t, tx.Signature, "AHRfOdsXX7E1+Kw6sQ/z2aTiqOeAlBowq6xEQM1d2MogzNZCA9SjOvOzrnd1ViBDpKIJKUQHtz5t0kJUDIZ1ngK8xcPyFlpe4w6DaHjUIF/RAjIloXaXfRqE1qQi6eJ9LQ==")

}

// https://suiexplorer.com/txblock/E5iw9P4VjJugZcAAtRhqzBrtjTK6aHyrqyvGUe4QUmBA
func TestExecuteSplit2(t *testing.T) {
	key := "//todo please replace your hex sui key"
	addr := NewAddress(key)
	fmt.Println("address", addr)
	split, err := BuildSplitTx(addr, addr, []*SuiObjectRef{{Digest: "8m7eNZYu9uZsbghPJy8gKxQq27iMPYF2uE1LHbBVJQyS", ObjectId: "0x26c26d20be986dbd99d33ff1b3a0bb16437039d3de85b4cb9d56a3f57066ef54", Version: 1978809}}, []uint64{1, 2, 3, 4, 5, 6, 7}, 0, 5582000, 815)
	if err != nil {
		t.Fatal(err)
	}
	data := toJson(&SplitSuiRequest{Coins: []*SuiObjectRef{{Digest: "8m7eNZYu9uZsbghPJy8gKxQq27iMPYF2uE1LHbBVJQyS", ObjectId: "0x26c26d20be986dbd99d33ff1b3a0bb16437039d3de85b4cb9d56a3f57066ef54", Version: 1978809}}, Amounts: []uint64{1, 2, 3, 4, 5, 6, 7}})
	fmt.Println(data)
	res, err := prepareTx(&Request{Data: data, Type: Split}, addr, 5582000, 815, addr)
	if err != nil {
		t.Fatal(err)
	}
	assert.Equal(t, base64.StdEncoding.EncodeToString(split), res)
}

func TestExecuteSplit(t *testing.T) {
	key := "//todo please replace your hex sui key"
	addr := NewAddress(key)
	fmt.Println("address", addr)
	split, err := BuildSplitTx(addr, addr, []*SuiObjectRef{{Digest: "8m7eNZYu9uZsbghPJy8gKxQq27iMPYF2uE1LHbBVJQyS", ObjectId: "0x26c26d20be986dbd99d33ff1b3a0bb16437039d3de85b4cb9d56a3f57066ef54", Version: 1978809}}, []uint64{1, 2, 3, 4, 5, 6, 7}, 0, 5582000, 815)
	if err != nil {
		t.Fatal(err)
	}
	tx, err := SignTransaction(base64.StdEncoding.EncodeToString(split), key)
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(fmt.Sprintf("\"%s\",[\"%s\"]", tx.TxBytes, tx.Signature))
	pay := Pay{}
	if err := json.Unmarshal([]byte(txJson5), &pay); err != nil {
		t.Fatal(err)
	}
	data, err := pay.Build()
	if err != nil {
		t.Fatal(err)
	}
	assert.Equal(t, split, data)
	fmt.Println(base64.StdEncoding.EncodeToString(data))
	tx, err = SignTransaction(base64.StdEncoding.EncodeToString(data), key)
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(fmt.Sprintf("\"%s\",[\"%s\"]", tx.TxBytes, tx.Signature))

	assert.Equal(t, tx.TxBytes, "AAAIAAgBAAAAAAAAAAAIAgAAAAAAAAAACAMAAAAAAAAAAAgEAAAAAAAAAAAIBQAAAAAAAAAACAYAAAAAAAAAAAgHAAAAAAAAAAAgHn+l/Ua9+OwSkcpSCEvb7qveazurOl2ebPg9eAYd5hkCAgAHAQAAAQEAAQIAAQMAAQQAAQUAAQYAAQcDAAAAAAMAAAEAAwAAAgADAAADAAMAAAQAAwAABQADAAAGAAEHAB5/pf1GvfjsEpHKUghL2+6r3ms7qzpdnmz4PXgGHeYZASbCbSC+mG29mdM/8bOguxZDcDnT3oW0y51Wo/VwZu9UuTEeAAAAAAAgc07geh9mPnCD+GJ0yLZkf47CF6K7ApzKJC8B22LwkS0ef6X9Rr347BKRylIIS9vuq95rO6s6XZ5s+D14Bh3mGS8DAAAAAAAAsCxVAAAAAAAA")
	assert.Equal(t, tx.Signature, "AJNML/3c0Aczg2xhjpjmFj59bGm6lqQKI+iAk0DMUyQpjcxFvi/wEztnqB8tUaUMevyRepmZk7gpsElU4htEdwC8xcPyFlpe4w6DaHjUIF/RAjIloXaXfRqE1qQi6eJ9LQ==")
}

// https://suiexplorer.com/txblock/96nEuw4B6xZVxLUvQuq4AFaWVv7cqJaV27U8MmT9gBBW
func TestExecuteMerge2(t *testing.T) {
	key := "//todo please replace your hex sui key"
	addr := NewAddress(key)
	merge, err := BuildMergeTx(addr, []*SuiObjectRef{{Digest: "3do9tdbdBCUMCn5rNweGt2Ag41fjubwr4WPkyMCiR6zv", ObjectId: "0x26c26d20be986dbd99d33ff1b3a0bb16437039d3de85b4cb9d56a3f57066ef54", Version: 1978808}},
		[]*SuiObjectRef{{Digest: "DVy7dp7u9UWgh1JA1EbRfALxYPL9qRJhSXC5PZguWRMy", ObjectId: "0xa7917d2b15ec1660b8a6658c29f99e5ba6feeb0cde66c7c5d072bd69b8574e30", Version: 1978808},
			{Digest: "GyKYSumghy7yGij9bKcf3wWL7f2M5nSLBiiE82TitmkT", ObjectId: "0xd228b5a3f214a9a345399e8948a9844e5a0a02bbda221c2e5ba8b63ede9390d4", Version: 1978808},
			{Digest: "H4LGpfZMZuLWdiDdkJRBgT1pSRrcEkLyfWnkCwYtayLG", ObjectId: "0xa21833a375925ff6352db9b1ef547a8da9cfcb717f60efa3edd243a1e3cee1b0", Version: 1978808}},
		0, 1630000, 815)
	if err != nil {
		t.Fatal(err)
	}
	data := toJson(&MergeSuiRequest{Coins: []*SuiObjectRef{{Digest: "3do9tdbdBCUMCn5rNweGt2Ag41fjubwr4WPkyMCiR6zv", ObjectId: "0x26c26d20be986dbd99d33ff1b3a0bb16437039d3de85b4cb9d56a3f57066ef54", Version: 1978808}},
		Objects: []*SuiObjectRef{{Digest: "DVy7dp7u9UWgh1JA1EbRfALxYPL9qRJhSXC5PZguWRMy", ObjectId: "0xa7917d2b15ec1660b8a6658c29f99e5ba6feeb0cde66c7c5d072bd69b8574e30", Version: 1978808},
			{Digest: "GyKYSumghy7yGij9bKcf3wWL7f2M5nSLBiiE82TitmkT", ObjectId: "0xd228b5a3f214a9a345399e8948a9844e5a0a02bbda221c2e5ba8b63ede9390d4", Version: 1978808},
			{Digest: "H4LGpfZMZuLWdiDdkJRBgT1pSRrcEkLyfWnkCwYtayLG", ObjectId: "0xa21833a375925ff6352db9b1ef547a8da9cfcb717f60efa3edd243a1e3cee1b0", Version: 1978808}}})
	fmt.Println(data)
	res, err := prepareTx(&Request{Data: data, Type: Merge}, addr, 1630000, 815, addr)
	if err != nil {
		t.Fatal(err)
	}
	assert.Equal(t, base64.StdEncoding.EncodeToString(merge), res)
}

func TestExecuteMerge(t *testing.T) {
	key := "//todo please replace your hex sui key"
	addr := NewAddress(key)
	fmt.Println("address", addr)
	merge, err := BuildMergeTx(addr, []*SuiObjectRef{{Digest: "3do9tdbdBCUMCn5rNweGt2Ag41fjubwr4WPkyMCiR6zv", ObjectId: "0x26c26d20be986dbd99d33ff1b3a0bb16437039d3de85b4cb9d56a3f57066ef54", Version: 1978808}},
		[]*SuiObjectRef{{Digest: "DVy7dp7u9UWgh1JA1EbRfALxYPL9qRJhSXC5PZguWRMy", ObjectId: "0xa7917d2b15ec1660b8a6658c29f99e5ba6feeb0cde66c7c5d072bd69b8574e30", Version: 1978808},
			{Digest: "GyKYSumghy7yGij9bKcf3wWL7f2M5nSLBiiE82TitmkT", ObjectId: "0xd228b5a3f214a9a345399e8948a9844e5a0a02bbda221c2e5ba8b63ede9390d4", Version: 1978808},
			{Digest: "H4LGpfZMZuLWdiDdkJRBgT1pSRrcEkLyfWnkCwYtayLG", ObjectId: "0xa21833a375925ff6352db9b1ef547a8da9cfcb717f60efa3edd243a1e3cee1b0", Version: 1978808}},
		0, 1630000, 815)
	if err != nil {
		t.Fatal(err)
	}
	tx, err := SignTransaction(base64.StdEncoding.EncodeToString(merge), key)
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(fmt.Sprintf("\"%s\",[\"%s\"]", tx.TxBytes, tx.Signature))
	pay := Pay{}
	if err := json.Unmarshal([]byte(mergeJson6), &pay); err != nil {
		t.Fatal(err)
	}
	data, err := pay.Build()
	if err != nil {
		t.Fatal(err)
	}
	assert.Equal(t, merge, data)
	fmt.Println(base64.StdEncoding.EncodeToString(data))
	tx, err = SignTransaction(base64.StdEncoding.EncodeToString(data), key)
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(fmt.Sprintf("\"%s\",[\"%s\"]", tx.TxBytes, tx.Signature))

	assert.Equal(t, tx.TxBytes, "AAADAQCnkX0rFewWYLimZYwp+Z5bpv7rDN5mx8XQcr1puFdOMLgxHgAAAAAAILm39FsWLIpL4vns+V6pEcpL0PN7VJ8ORuFfODlAlVpIAQDSKLWj8hSpo0U5nolIqYROWgoCu9oiHC5bqLY+3pOQ1LgxHgAAAAAAIO1My3/LXFaB95we+up5H3Luqoeyrs1yd2DnI6RIrraQAQCiGDOjdZJf9jUtubHvVHqNqc/LcX9g76Pt0kOh487hsLgxHgAAAAAAIO6VhmpyEnX3NyXrpnzYsmhaT7No/umICWshdYQTiyH1AQMBAgACAQAAAQEAHn+l/Ua9+OwSkcpSCEvb7qveazurOl2ebPg9eAYd5hkBJsJtIL6Ybb2Z0z/xs6C7FkNwOdPehbTLnVaj9XBm71S4MR4AAAAAACAnJHf6AKPRorbjT2n++V1UWg5HVnWH6V/NoIeD7ff0Ax5/pf1GvfjsEpHKUghL2+6r3ms7qzpdnmz4PXgGHeYZLwMAAAAAAAAw3xgAAAAAAAA=")
	assert.Equal(t, tx.Signature, "AER61pCwYdZPx3akQQVMFhJmhqRTMx86AoMVqfYuzj0z5On8cCd4kCHeK7CkKaYbKWx/T7dDGfB5eSsPRLJNtQy8xcPyFlpe4w6DaHjUIF/RAjIloXaXfRqE1qQi6eJ9LQ==")
}

// https://suiexplorer.com/txblock/36PnrdSemNtEWFXabDF1WJTT8ANuzdFwX7gmYHSE5PPx
func TestExecuteMul(t *testing.T) {
	key := "//todo please replace your hex sui key"
	addr := NewAddress(key)
	fmt.Println("address", addr)
	mul, err := BuildMulTx(addr, []*SuiObjectRef{{Digest: "8neVLSnZEDGjWq5ynYMP7sbCT7Tr1bLXnDTjw9dYuX58", ObjectId: "0x26c26d20be986dbd99d33ff1b3a0bb16437039d3de85b4cb9d56a3f57066ef54", Version: 1978814}}, map[string]uint64{"0x1e7fa5fd46bdf8ec1291ca52084bdbeeabde6b3bab3a5d9e6cf83d78061de619": 1, "0x215d3a67d951ebd5b453b440497917b5fac2890fc7f18358322d372e2f13045d": 2}, 0, 9534000, 815)
	if err != nil {
		t.Fatal(err)
	}
	tx, err := SignTransaction(base64.StdEncoding.EncodeToString(mul), key)
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(fmt.Sprintf("\"%s\",[\"%s\"]", tx.TxBytes, tx.Signature))
	pay := Pay{}
	if err := json.Unmarshal([]byte(mulJSon), &pay); err != nil {
		t.Fatal(err)
	}
	data, err := pay.Build()
	if err != nil {
		t.Fatal(err)
	}
	//assert.Equal(t, merge, data)
	fmt.Println(base64.StdEncoding.EncodeToString(data))
	tx, err = SignTransaction(base64.StdEncoding.EncodeToString(data), key)
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(fmt.Sprintf("\"%s\",[\"%s\"]", tx.TxBytes, tx.Signature))

	assert.Equal(t, tx.TxBytes, "AAAEAAgBAAAAAAAAAAAIAgAAAAAAAAAAIB5/pf1GvfjsEpHKUghL2+6r3ms7qzpdnmz4PXgGHeYZACAhXTpn2VHr1bRTtEBJeRe1+sKJD8fxg1gyLTcuLxMEXQMCAAIBAAABAQABAQMAAAAAAQIAAQEDAAABAAEDAB5/pf1GvfjsEpHKUghL2+6r3ms7qzpdnmz4PXgGHeYZASbCbSC+mG29mdM/8bOguxZDcDnT3oW0y51Wo/VwZu9UvjEeAAAAAAAgc7NVVGHTJpRDCtFbtTFvU3yDOVEYvLbkomev5JKeybcef6X9Rr347BKRylIIS9vuq95rO6s6XZ5s+D14Bh3mGS8DAAAAAAAAMHqRAAAAAAAA")
	assert.Equal(t, tx.Signature, "APpwRQQF1TPnTa1siSo6Yqcgv3fuiL2s6rpiMSp9ijqpHUnN8wev4jdVD699+OphlV5dWhJ9mfVDLgEQ1smUdwa8xcPyFlpe4w6DaHjUIF/RAjIloXaXfRqE1qQi6eJ9LQ==")
}

func TestEqualToken4(t *testing.T) {
	pay1 := Pay{}
	if err := json.Unmarshal([]byte(txJson4), &pay1); err != nil {
		t.Fatal(err)
	}
	data, err := pay1.Build()
	if err != nil {
		t.Fatal(err)
	}
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
	data2, err := BuildTokenTx("0x1e7fa5fd46bdf8ec1291ca52084bdbeeabde6b3bab3a5d9e6cf83d78061de619", "0x1e7fa5fd46bdf8ec1291ca52084bdbeeabde6b3bab3a5d9e6cf83d78061de619", coins, tokens, 1, 0, 3000000, 1000)
	if err != nil {
		t.Fatal(err)
	}
	//fmt.Println(data2)
	if !bytes.Equal(data2, data) {
		fmt.Println(data2)
		fmt.Println(data)
		fmt.Println(len(data2), len(data))
		for i := 0; i < len(data) && i < len(data2); i++ {
			if data[i] != data2[i] {
				fmt.Println("不同", i, data2[i], data[i])
			}
		}
		t.Fatal()
	}
}

// 7dz7UduUssmJmNWCBeHEyZrxLsXnwMyzPrZKHpe85HaY
func TestStake(t *testing.T) {
	data, err := BuildStakeTx("0x1e7fa5fd46bdf8ec1291ca52084bdbeeabde6b3bab3a5d9e6cf83d78061de619", "0x72169c90b7ea87f8101285c849c09cacced9968f83aa30786dad546bb94c78ab",
		[]*SuiObjectRef{{ObjectId: "0x26c26d20be986dbd99d33ff1b3a0bb16437039d3de85b4cb9d56a3f57066ef54", Digest: "AMGM65x2qTfM4kfPjbv7Aqpap6MBiVVa4W8hrakgvPjB", Version: 1978816}},
		1000000000, 0, 9644512, 820)
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(base64.StdEncoding.EncodeToString(data))
	key := "//todo please replace your hex sui key"
	tx, err := SignTransaction(base64.StdEncoding.EncodeToString(data), key)
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(fmt.Sprintf("\"%s\",[\"%s\"]", tx.TxBytes, tx.Signature))
	assert.Equal(t, tx.TxBytes, "AAADAAgAypo7AAAAAAEBAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAUBAAAAAAAAAAEAIHIWnJC36of4EBKFyEnAnKzO2ZaPg6oweG2tVGu5THirAgIAAQEAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAwpzdWlfc3lzdGVtEXJlcXVlc3RfYWRkX3N0YWtlAAMBAQACAAABAgAef6X9Rr347BKRylIIS9vuq95rO6s6XZ5s+D14Bh3mGQEmwm0gvphtvZnTP/GzoLsWQ3A5096FtMudVqP1cGbvVMAxHgAAAAAAIIrqJpLzvTYHdoUCDU1KmmXxL+/TEZdmF0A8w9aqCH7+Hn+l/Ua9+OwSkcpSCEvb7qveazurOl2ebPg9eAYd5hk0AwAAAAAAAOApkwAAAAAAAA==")
	assert.Equal(t, tx.Signature, "AJWF7M68Rk24Sol0RMy3ZfA5SEifvZglxjov/TY2nmI21SovAMfScklsz++QV8eQSJrwbQdvDIK/IOh7YuDB/Aa8xcPyFlpe4w6DaHjUIF/RAjIloXaXfRqE1qQi6eJ9LQ==")
}

func TestWithdraw(t *testing.T) {
	data, err := BuildWithdrawStakeTx("0x1e7fa5fd46bdf8ec1291ca52084bdbeeabde6b3bab3a5d9e6cf83d78061de619",
		[]*SuiObjectRef{{ObjectId: "0x26c26d20be986dbd99d33ff1b3a0bb16437039d3de85b4cb9d56a3f57066ef54", Digest: "AihBh2VjG96NDTCw1HvZj8TtWEpmYZUh9rb1D92oQ7Ak", Version: 5656730}},
		&SuiObjectRef{Digest: "CkmUVCkHFWyjH27Zg5xTd5xbUQZt1BReQnvtq2zeT6zW", Version: 5656730, ObjectId: "0x194acb4ec803ef63f15331efa9e701b4a334cf417fa15432d736d90978ce43e4"}, 0, 9534000, 820)
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(base64.StdEncoding.EncodeToString(data))
	key := "//todo please replace your hex sui key"
	tx, err := SignTransaction(base64.StdEncoding.EncodeToString(data), key)
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(fmt.Sprintf("\"%s\",[\"%s\"]", tx.TxBytes, tx.Signature))
	assert.Equal(t, tx.TxBytes, "AAACAQEAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAABQEAAAAAAAAAAQEAGUrLTsgD72PxUzHvqecBtKM0z0F/oVQy1zbZCXjOQ+SaUFYAAAAAACCuptFaCeiByBnx2kjvoGQSG8HVSQDZY4/GJYRmToFZ+wEAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAMKc3VpX3N5c3RlbRZyZXF1ZXN0X3dpdGhkcmF3X3N0YWtlAAIBAAABAQAef6X9Rr347BKRylIIS9vuq95rO6s6XZ5s+D14Bh3mGQEmwm0gvphtvZnTP/GzoLsWQ3A5096FtMudVqP1cGbvVJpQVgAAAAAAIJBnbps0rg7WV3qFRlezNq7ebV0Hry0IL/LATrprMgElHn+l/Ua9+OwSkcpSCEvb7qveazurOl2ebPg9eAYd5hk0AwAAAAAAADB6kQAAAAAAAA==")
	assert.Equal(t, tx.Signature, "AP2tuYbwdmTTZRtJmYOQozNfEK4Cs6UGhgRXgSDULDpulWXS/bPqftk1wgNacLNTuFbvXrPoLINOXFHPRxH2uwe8xcPyFlpe4w6DaHjUIF/RAjIloXaXfRqE1qQi6eJ9LQ==")
}

func TestEqualToken5(t *testing.T) {
	pay1 := Pay{}
	if err := json.Unmarshal([]byte(txJson4), &pay1); err != nil {
		t.Fatal(err)
	}
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
	data2, err := BuildTokenTx("0x1e7fa5fd46bdf8ec1291ca52084bdbeeabde6b3bab3a5d9e6cf83d78061de619", "0x1e7fa5fd46bdf8ec1291ca52084bdbeeabde6b3bab3a5d9e6cf83d78061de619", coins, tokens, 1, 0, 3000000, 1000)
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(data2)
}
