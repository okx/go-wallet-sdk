package cardano

type MultiAssets struct {
	PolicyId string   `json:"policyId"`
	Assets   []*Asset `json:"assets"`
}
type Asset struct {
	AssetName string `json:"assetName"`
	Amount    uint64 `json:"amount"`
}

type TxIn struct {
	TxId       string         `json:"txId"`
	Index      uint64         `json:"index"`
	Amount     uint64         `json:"amount"`
	MultiAsset []*MultiAssets `json:"multiAsset"`
}

type TxData struct {
	PrvKey        string         `json:"privateKey"`
	Inputs        []*TxIn        `json:"inputs"`
	ToAddress     string         `json:"toAddress"`
	Amount        uint64         `json:"amount"`
	MultiAsset    []*MultiAssets `json:"multiAsset"`
	ChangeAddress string         `json:"changeAddress"`
	Max           bool           `json:"max"`
	TTL           uint64         `json:"ttl"`
}

type MinFeeData struct {
	Valid  bool   `json:"valid"`
	Fee    uint64 `json:"fee"`
	Change uint64 `json:"change"`
}
