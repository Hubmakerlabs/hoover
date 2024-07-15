package schema

const (
	TxVersionV1 = "v1"

	TxActionTransfer = "transfer"
	TxActionMint     = "mint"
	TxActionBurn     = "burn"

	ArAddress   = "AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA"
	EvmAddress  = "0x0000000000000000000000000000000000000000"
	EthAddress  = EvmAddress
	MoonAddress = EvmAddress // moonbeam native token GLMR
	CfxAddress  = EvmAddress
)
