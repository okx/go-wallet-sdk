/**
Author： https://github.com/xssnick/tonutils-go
*/

package wallet

import (
	"encoding/hex"
	"fmt"
	"testing"
)

func TestAddressFromPubKey(t *testing.T) {
	pkey, _ := hex.DecodeString("dcc39550bb494f4b493e7efe1aa18ea31470f33a2553c568cb74a17ed56790c1")

	a, err := AddressFromPubKey(pkey, V4R2, DefaultSubwallet)
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(a)
	if a.String() != "UQAwwdowWbBKkrnRlbY8CUEzy_pgK9pIvOKP2eqcD01EWl2U" {
		t.Fatal("v3 not match")
	}
}
