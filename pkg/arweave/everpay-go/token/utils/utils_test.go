package utils

import (
	"math/big"
	"testing"

	"github.com/Hubmakerlabs/hoover/pkg/arweave/everpay-go/token/schema"

	"github.com/stretchr/testify/assert"
)

func TestTag(t *testing.T) {
	assert.Equal(t, "ethereum-usdt-0xabc",
		Tag(schema.ChainTypeEth, "USDT", "0xABC"))
	assert.Equal(t, "err_invalid_token",
		Tag(schema.ChainTypeCrossArEth, "AR", "0xABC"))
	assert.Equal(t, "arweave,ethereum-ar-ABC,0xabc",
		Tag(schema.ChainTypeCrossArEth, "AR", "ABC,0xABC"))
}

func TestToken_SplitTargetAmount(t *testing.T) {
	everTokenDecimal := 18
	targetAmount := big.NewInt(1)
	amount, err := ZoomMintAmount(everTokenDecimal, targetAmount, 0)
	assert.NoError(t, err)
	assert.Equal(t, "1000000000000000000", amount.String())

	targetAmount = big.NewInt(1)
	amount, err = ZoomMintAmount(everTokenDecimal, targetAmount, 18)
	assert.NoError(t, err)
	assert.Equal(t, "1", amount.String())

	everTokenDecimal = 6
	targetAmount = big.NewInt(100)
	amount, err = ZoomMintAmount(everTokenDecimal, targetAmount, 10)
	assert.NoError(t, err)
	assert.Equal(t, "0",
		amount.String()) // return 0, because diffDecimal = -4, amount == 0.001 but amount must integer so is 0

	targetAmount = big.NewInt(999999)
	amount, err = ZoomMintAmount(everTokenDecimal, targetAmount, 10)
	assert.NoError(t, err)
	assert.Equal(t, "99", amount.String())

	targetAmount = big.NewInt(-999999)
	amount, err = ZoomMintAmount(everTokenDecimal, targetAmount, 10)
	assert.NoError(t, err)
	assert.Equal(t, "-99", amount.String())

	everTokenDecimal = 6
	targetAmount = big.NewInt(-199999)
	amount, err = ZoomMintAmount(everTokenDecimal, targetAmount, 6)
	assert.NoError(t, err)
	assert.Equal(t, "-199999", amount.String())

	targetAmount = big.NewInt(90000)
	amount, err = ZoomMintAmount(everTokenDecimal, targetAmount, 6)
	assert.NoError(t, err)
	assert.Equal(t, "90000", amount.String())

	// arweave pst token
	everTokenDecimal = 18
	targetAmount = big.NewInt(1)
	amount, err = ZoomMintAmount(everTokenDecimal, targetAmount, 0)
	assert.NoError(t, err)
	assert.Equal(t, "1000000000000000000", amount.String())
}

func TestToken_CombineBurnAmount(t *testing.T) {
	everTokenDecimal := 18
	targetAmount := big.NewInt(9000000000000000000)
	targetDecimal := 0
	amount, err := ZoomBurnAmount(everTokenDecimal, targetDecimal, targetAmount)
	assert.NoError(t, err)
	assert.Equal(t, "9", amount.String())

	targetAmount = big.NewInt(8900000000000000000)
	amount, err = ZoomBurnAmount(everTokenDecimal, targetDecimal, targetAmount)
	assert.Error(t, err, `targetAmount combine is float; amount: 8.9`)

	everTokenDecimal = 6
	targetAmount = big.NewInt(666)
	targetDecimal = 12
	amount, err = ZoomBurnAmount(everTokenDecimal, targetDecimal, targetAmount)
	assert.NoError(t, err)
	assert.Equal(t, "666000000", amount.String())

	everTokenDecimal = 6
	targetAmount = big.NewInt(666)
	targetDecimal = 6
	amount, err = ZoomBurnAmount(everTokenDecimal, targetDecimal, targetAmount)
	assert.NoError(t, err)
	assert.Equal(t, "666", amount.String())

	everTokenDecimal = 60
	targetAmount, ok := new(big.Int).SetString("10000000000000000000000000000000000000000000000000000000000001",
		10)
	assert.Equal(t, true, ok)

	amount, err = ZoomBurnAmount(everTokenDecimal, 0, targetAmount)
	assert.Equal(t,
		`everAmount zoom can not be float type; amount: 10.000000000000000000000000000000000000000000000000000000000001`,
		err.Error())

	// pst

	everTokenDecimal = 18
	targetAmount, ok = new(big.Int).SetString("3000000000000000000", 10)
	assert.Equal(t, true, ok)
	amount, err = ZoomBurnAmount(everTokenDecimal, 0, targetAmount)
	assert.NoError(t, err)
	assert.Equal(t, "3", amount.String())
}

func TestTagDecode(t *testing.T) {
	tag := "arweave-ardrive--8a6rexfkpfwwuyvo98wzsfzh0d6vjui-butjvlwojq"
	chainType, tokenSymbol, tokenId, err := TagDecode(tag)
	assert.NoError(t, err)
	assert.Equal(t, "arweave", chainType)
	assert.Equal(t, "ardrive", tokenSymbol)
	assert.Equal(t, "-8a6rexfkpfwwuyvo98wzsfzh0d6vjui-butjvlwojq", tokenId)
}
