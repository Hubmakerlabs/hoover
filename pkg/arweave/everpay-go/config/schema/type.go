package schema

type TokenFee struct {
	TokenTag    string            `json:"tokenTag"` // token tag
	TransferFee string            `json:"transferFee"`
	BundleFee   string            `json:"bundleFee"`
	BurnFeeMap  map[string]string `json:"burnFeeMap"` // key: targetChainType, val: burnFee
}

type EvmGasInfo struct {
	ChainType string  `json:"chainType"`
	Fastest   float64 `json:"fastest"`
	Source    string  `json:"source"`
}
