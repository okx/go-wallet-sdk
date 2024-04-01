package atomical

import (
	"github.com/stretchr/testify/assert"
	"testing"

	"github.com/btcsuite/btcd/chaincfg"
)

func TestAtomicalTransfer(t *testing.T) {
	network := &chaincfg.MainNetParams

	commitTxPrevOutputList := make([]*PrevOutput, 0)
	commitTxPrevOutputList = append(commitTxPrevOutputList, &PrevOutput{
		TxId:       "c54e134783156d85de4f8f281669b5c53c6748f886252092c25b0249fb7ff95a",
		VOut:       0,
		Amount:     1200,
		Address:    "1PY7wFLq74G1yDkuvM9g5isWbwrZ1DiM2Y",
		PrivateKey: "cPnvkvUYyHcSSS26iD1dkrJdV7k1RoUqJLhn3CYxpo398PdLVE22",
		Data:       AtomicalDatas{{AtomicalId: "9527efa43262636d8f5917fc763fbdd09333e4b387afd6d4ed7a905a127b27b4i0", Type: "FT"}},
	})

	commitTxPrevOutputList = append(commitTxPrevOutputList, &PrevOutput{
		TxId:       "c54e134783156d85de4f8f281669b5c53c6748f886252092c25b0249fb7ff95a",
		VOut:       1,
		Amount:     140700,
		Address:    "1PY7wFLq74G1yDkuvM9g5isWbwrZ1DiM2Y",
		PrivateKey: "cPnvkvUYyHcSSS26iD1dkrJdV7k1RoUqJLhn3CYxpo398PdLVE22",
	})

	outputs := make([]*Output, 0)
	outputs = append(outputs, &Output{
		Amount:  600,
		Address: "1PY7wFLq74G1yDkuvM9g5isWbwrZ1DiM2Y",
		Data:    []*AtomicalData{{AtomicalId: "9527efa43262636d8f5917fc763fbdd09333e4b387afd6d4ed7a905a127b27b4i0", Type: "FT"}},
	})

	request := &AtomicalRequest{
		Inputs:   commitTxPrevOutputList,
		Outputs:  outputs,
		FeePerB:  100,
		DustSize: 546,
		Address:  "1PY7wFLq74G1yDkuvM9g5isWbwrZ1DiM2Y",
	}
	txs, err := AtomicalTransfer(network, request)
	assert.Nil(t, err)
	assert.Equal(t, len(txs.RevealTxs), 0)
	assert.Equal(t, txs.CommitTxFee, int64(40600))
	assert.Equal(t, txs.CommitTx, "02000000025af97ffb49025bc292202586f848673cc5b56916288f4fde856d158347134ec5000000006a4730440220602df3a3b209ed42aa14fb23ce3aeef9534b9a5d0b05b5a9deca10eaf935625602200c4efadea3c9a638d9322d2c9edab3204465f3a1ea95ee20178c9ea44ae37e8b01210357bbb2d4a9cb8a2357633f201b9c518c2795ded682b7913c6beef3fe23bd6d2ffdffffff5af97ffb49025bc292202586f848673cc5b56916288f4fde856d158347134ec5010000006a473044022002bd581aadb20cdb8fac55e88258b2e69df38a07b3a4d7ba109b53e7cd293e79022054c746a64cc3ac065819fc37ad6cadd10db4e31f15d15309f2f075fe4e00360001210357bbb2d4a9cb8a2357633f201b9c518c2795ded682b7913c6beef3fe23bd6d2ffdffffff0358020000000000001976a914f735738be9baec76cda18724e82c08f50c41a58088ac58020000000000001976a914f735738be9baec76cda18724e82c08f50c41a58088ac04870100000000001976a914f735738be9baec76cda18724e82c08f50c41a58088ac00000000")
}

func TestAtomicalNFTTransfer(t *testing.T) {
	network := &chaincfg.MainNetParams

	commitTxPrevOutputList := make([]*PrevOutput, 0)
	commitTxPrevOutputList = append(commitTxPrevOutputList, &PrevOutput{
		TxId:       "c54e134783156d85de4f8f281669b5c53c6748f886252092c25b0249fb7ff95a",
		VOut:       0,
		Amount:     1200,
		Address:    "1PY7wFLq74G1yDkuvM9g5isWbwrZ1DiM2Y",
		PrivateKey: "cPnvkvUYyHcSSS26iD1dkrJdV7k1RoUqJLhn3CYxpo398PdLVE22",
		Data:       AtomicalDatas{{AtomicalId: "9527efa43262636d8f5917fc763fbdd09333e4b387afd6d4ed7a905a127b27b4i0", Type: "NFT"}},
	})

	commitTxPrevOutputList = append(commitTxPrevOutputList, &PrevOutput{
		TxId:       "c54e134783156d85de4f8f281669b5c53c6748f886252092c25b0249fb7ff95a",
		VOut:       1,
		Amount:     140700,
		Address:    "1PY7wFLq74G1yDkuvM9g5isWbwrZ1DiM2Y",
		PrivateKey: "cPnvkvUYyHcSSS26iD1dkrJdV7k1RoUqJLhn3CYxpo398PdLVE22",
	})

	outputs := make([]*Output, 0)
	outputs = append(outputs, &Output{
		Amount:  600,
		Address: "1PY7wFLq74G1yDkuvM9g5isWbwrZ1DiM2Y",
		Data:    []*AtomicalData{{AtomicalId: "9527efa43262636d8f5917fc763fbdd09333e4b387afd6d4ed7a905a127b27b4i0", Type: "NFT"}},
	})

	request := &AtomicalRequest{
		Inputs:   commitTxPrevOutputList,
		Outputs:  outputs,
		FeePerB:  100,
		DustSize: 546,
		Address:  "1PY7wFLq74G1yDkuvM9g5isWbwrZ1DiM2Y",
	}
	txs, err := AtomicalTransfer(network, request)
	assert.Nil(t, err)
	assert.Equal(t, len(txs.RevealTxs), 0)
	assert.Equal(t, txs.CommitTxFee, int64(37200))
	assert.Equal(t, txs.CommitTx, "02000000025af97ffb49025bc292202586f848673cc5b56916288f4fde856d158347134ec5000000006a47304402207ecb66373f06b0c6c6323c76938e79734c9fb8b60c1bb83138b929783c123f9d02207a4d78552b6b09c6ef8bce647a06e9ebd09962a459241df06f6cfb4ee88d273101210357bbb2d4a9cb8a2357633f201b9c518c2795ded682b7913c6beef3fe23bd6d2ffdffffff5af97ffb49025bc292202586f848673cc5b56916288f4fde856d158347134ec5010000006a473044022044c1cbaae6b8bace7eb5605eee76437027c98554621b189c18cca3b932031a7702207d5e9190291b9cf3b2b5144fcb5150bbb1372f1f482923a86e62f4d05c494bc201210357bbb2d4a9cb8a2357633f201b9c518c2795ded682b7913c6beef3fe23bd6d2ffdffffff0258020000000000001976a914f735738be9baec76cda18724e82c08f50c41a58088aca4960100000000001976a914f735738be9baec76cda18724e82c08f50c41a58088ac00000000")
}
