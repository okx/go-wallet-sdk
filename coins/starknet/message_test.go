package starknet

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestVerifyMsgSign(t *testing.T) {
	ok := VerifyMsgSign("02f4a65ecea5351f49f181841bdddcdf62f600d0e4864755699386d42dd17e37", "0x1d6f9ddef6e87e75a12850ac21bfa1a9dccabf97d5efc311d16e1939df367c1", `{"publicKey":"0x02f4a65ecea5351f49f181841bdddcdf62f600d0e4864755699386d42dd17e37","publicKeyY":"0x0250990eae46b48f5dffbca10b2e71a25b62c943972bb16d8b36a09d927170af","signedDataR":"0x07294c0ab1106743e278c0e5b0e51d2f85a2d220523f9562885d554e7617ee30","signedDataS":"0x072ab5f6bb8e3451493677114564c72af598279f32229234978170303ae977f5"}`)
	assert.True(t, ok)
}
