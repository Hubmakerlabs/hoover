package schema

const (
	TxVersionV1 = "v1"

	TxActionTransfer        = "transfer"
	TxActionMint            = "mint"
	TxActionBurn            = "burn"
	TxActionTransferOwner   = "transferOwner"
	TxActionAddWhiteList    = "addWhiteList"
	TxActionRemoveWhiteList = "removeWhiteList"
	TxActionPauseWhiteList  = "pauseWhiteList"
	TxActionAddBlackList    = "addBlackList"
	TxActionRemoveBlackList = "removeBlackList"
	TxActionPauseBlackList  = "pauseBlackList"
	TxActionPause           = "pause"

	ArAddress     = "AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA"
	EvmAddress    = "0x0000000000000000000000000000000000000000"
	EthAddress    = EvmAddress
	MoonAddress   = EvmAddress // moonbeam native token GLMR
	CfxAddress    = EvmAddress
	BscAddress    = EvmAddress
	PlatonAddress = EvmAddress
)

// token type
const (
	TNS101Type = 101
	TNS102Type = 102
)

const (
	ZeroAddress = "0x0000000000000000000000000000000000000000"
)
