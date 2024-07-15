package token

import (
	"math/big"
	"sync"

	"github.com/Hubmakerlabs/hoover/pkg/arweave/go-everpay/token/schema"
	"github.com/Hubmakerlabs/hoover/pkg/arweave/go-everpay/token/utils"
)

type Token struct {
	ID           string // On Native-Chain tokenId; Special AR token: "AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA,0xcc9141efa8c20c7df0778748255b1487957811be"
	Symbol       string
	Decimals     int                           // On everPay decimals
	ChainType    string                        // On everPay chainType; Special AR token: "arweave,ethereum"
	ChainID      string                        // On everPay chainId; Special AR token: "0,1"(mainnet) or "0,42"(testnet)
	targetChains map[string]schema.TargetChain // key: targetChainType

	TotalSupply *big.Int
	Balances    map[string]*big.Int // account id -> balance
	Minted      map[string]bool     // chainHash -> bool

	// backup
	balancesBackup map[string]*big.Int

	lock sync.RWMutex
}

func New(tokenId, symbol, chainType, chainID string, everDecimals int,
	targetChainArr []schema.TargetChain) *Token {

	targetChainMap, err := genTargetInfoMap(targetChainArr)
	if err != nil {
		panic(err)
	}

	// everDecimals must >= chainDecimals
	for _, info := range targetChainMap {
		if info.Decimals > everDecimals {
			panic("everDecimals can not less than chainDecimals")
		}
	}

	return &Token{
		ID:        tokenId,
		Symbol:    symbol,
		Decimals:  everDecimals,
		ChainType: chainType,
		ChainID:   chainID,

		targetChains: targetChainMap,
		TotalSupply:  big.NewInt(0),
		Balances:     make(map[string]*big.Int),
		Minted:       make(map[string]bool),
	}
}

// Tag is the unique identifier of token
func (t *Token) Tag() string {
	return utils.Tag(t.ChainType, t.Symbol, t.ID)
}

func (t *Token) GetTargetChain(targetChainType string) (res schema.TargetChain,
	exist bool) {
	res, exist = t.targetChains[targetChainType]
	return
}

func (t *Token) GetTargetChains() map[string]schema.TargetChain {
	mmap := make(map[string]schema.TargetChain)
	for k, v := range t.targetChains {
		mmap[k] = v
	}
	return mmap
}

func genTargetInfoMap(targetChainArr []schema.TargetChain) (map[string]schema.TargetChain,
	error) {
	targetChainMap := make(map[string]schema.TargetChain)
	for _, info := range targetChainArr {
		targetChainMap[info.ChainType] = info
	}
	return targetChainMap, nil
}
