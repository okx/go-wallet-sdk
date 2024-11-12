package starknet

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestJsonToTypedData(t *testing.T) {
	jsonMsg := `{
    "accountAddress" : "0x06c3c93eeb1643740a80a338b9346c0c9a06177bfcc098a6d86e353532090ae4",
    "typedData" : {
          "types" : {
              "StarkNetDomain": [
                  { "name" : "name", "type" : "felt" },
                  { "name" : "version", "type" : "felt" },
                  { "name" : "chainId", "type" : "felt" }
              ],
              "Person" : [
                  { "name": "name", "type" : "felt" },
                  { "name": "wallet", "type" : "felt" }
              ],
              "Mail": [
                  { "name": "from", "type": "Person" },
                  { "name": "to", "type": "Person" },
                  { "name": "contents", "type": "felt" }
              ]
          },
          "primaryType" : "Mail",
          "domain" : {
              "name" : "StarkNet Mail",
              "version" : "1",
              "chainId" : "1"
          },
          "message" : {
              "from" : {
                  "name" : "Cow",
                  "wallet" : "0xCD2a3d9F938E13CD947Ec05AbC7FE734Df8DD826"
              },
              "to": {
                  "name" : "Bob",
                  "wallet" : "0xbBbBBBBbbBBBbbbBbbBbbbbBBbBbbbbBbBbbBBbB"
              },
              "contents" : "Hello, Bob!"
          }
  }
}`
	hash, err := GetMessageHashWithJson(jsonMsg)
	if err != nil {
		t.Fatal(hash)
	}

	fmt.Println(hash)
	assert.Equal(t, hash, "0x45514f85d4e7e2d3db3aac059a5d937f6c5d0f61f87ba25fa138c038248ce7a")
}

type Mail struct {
	From     Person
	To       Person
	Contents string
}

type Person struct {
	Name   string
	Wallet string
}

func MockTypedDataMail() (ttd TypedData) {
	exampleTypes := make(map[string]TypeDef)
	domDefs := []Definition{{"name", "felt"}, {"version", "felt"}, {"chainId", "felt"}}
	exampleTypes["StarkNetDomain"] = TypeDef{Definitions: domDefs}
	mailDefs := []Definition{{"from", "Person"}, {"to", "Person"}, {"contents", "felt"}}
	exampleTypes["Mail"] = TypeDef{Definitions: mailDefs}
	persDefs := []Definition{{"name", "felt"}, {"wallet", "felt"}}
	exampleTypes["Person"] = TypeDef{Definitions: persDefs}

	dm := Domain{
		Name:    "StarkNet Mail",
		Version: "1",
		ChainId: "1",
	}
	fmt.Println(exampleTypes)
	ttd, _ = NewTypedData(exampleTypes, "Mail", dm)
	return ttd
}
