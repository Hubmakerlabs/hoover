package utils

import (
	"fmt"
	"math/big"
	"strings"

	"github.com/Hubmakerlabs/hoover/pkg/arweave/go-everpay/token/schema"
	"github.com/shopspring/decimal"
)

// Tag: for special token tag cal, Tag is the unique identifier of token
// notice: Arweave address is Case Sensitive
func Tag(chainType, tokenSymbol, tokenID string) string {
	// process tokenId
	var id string
	switch chainType {
	case schema.ChainTypeArweave:
		id = tokenID
	case schema.ChainTypeCrossArEth: // now only AR token
		ids := strings.Split(tokenID, ",")
		if len(ids) != 2 {
			return "err_invalid_token"
		}

		ids[1] = strings.ToLower(ids[1])
		id = strings.Join(ids, ",")
	default: // "ethereum", "avalanche" and so on evm chain
		id = strings.ToLower(tokenID)
	}

	return strings.ToLower(chainType+"-"+tokenSymbol) + "-" + id
}

func TagDecode(tokenTag string) (chainType, tokenSymbol, tokenId string,
	err error) {
	ss := strings.SplitN(tokenTag, "-", 3)
	if len(ss) != 3 {
		err = fmt.Errorf("tokenTag incorrect; tokenTag: %s", tokenTag)
		return
	}
	chainType = ss[0]
	tokenSymbol = ss[1]
	tokenId = ss[2]
	return
}

// ZoomBurnAmount zoom burn tx amount
func ZoomBurnAmount(everTokenDecimal, targetDecimal int,
	everAmount *big.Int) (targetAmount *big.Int, err error) {
	diffDecimal := targetDecimal - everTokenDecimal
	amount := decimal.NewFromBigInt(everAmount, int32(diffDecimal))

	// amount must be integer,if it is float return err
	amountStr := amount.String()
	if strings.Contains(amountStr, ".") {
		return nil, fmt.Errorf("everAmount zoom can not be float type; amount: %s",
			amountStr)
	}

	return amount.BigInt(), nil
}

// ZoomMintAmount token mint tx amount split
func ZoomMintAmount(everTokenDecimal int, targetAmount *big.Int,
	targetDecimal int) (everAmount *big.Int, err error) {
	diffDecimal := everTokenDecimal - targetDecimal
	amount := decimal.NewFromBigInt(targetAmount,
		int32(diffDecimal)) // if diffDecimal < 0, the amount could become 0
	return amount.BigInt(), nil
}
