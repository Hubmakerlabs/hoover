package pay

import (
	"testing"

	"github.com/Hubmakerlabs/hoover/pkg/arweave/go-everpay/pay/schema"
	tokSchema "github.com/Hubmakerlabs/hoover/pkg/arweave/go-everpay/token/schema"
	"github.com/stretchr/testify/assert"
)

func Test_getTxDataTargetChainType(t *testing.T) {
	// burn tx could have json data, could have targetChainType
	tx := schema.Transaction{
		Action:    tokSchema.TxActionBurn,
		ChainType: tokSchema.ChainTypeEth,
		Data:      "",
	}
	targetChainType, err := GetTargetChainTypeFromData(tx.Data, tx.Action,
		tx.ChainType)
	assert.NoError(t, err)
	assert.Equal(t, tokSchema.ChainTypeEth, targetChainType)

	tx = schema.Transaction{
		Action:    tokSchema.TxActionBurn,
		ChainType: tokSchema.ChainTypeEth,
		Data:      "{\"aaa\":\"bbb\"}",
	}
	targetChainType, err = GetTargetChainTypeFromData(tx.Data, tx.Action,
		tx.ChainType)
	assert.NoError(t, err)
	assert.Equal(t, tokSchema.ChainTypeEth, targetChainType)

	tx = schema.Transaction{
		Action:    tokSchema.TxActionBurn,
		ChainType: tokSchema.ChainTypeEth,
		Data:      "{\"targetChainType\":\"mmmm\"}",
	}
	targetChainType, err = GetTargetChainTypeFromData(tx.Data, tx.Action,
		tx.ChainType)
	assert.NoError(t, err)
	assert.Equal(t, "mmmm", targetChainType)

	// mint tx must have json data
	tx = schema.Transaction{
		Action:    tokSchema.TxActionMint,
		ChainType: tokSchema.ChainTypeEth,
		Data:      "aaaa",
	}
	targetChainType, err = GetTargetChainTypeFromData(tx.Data, tx.Action,
		tx.ChainType)
	assert.Equal(t, "invalid character 'a' looking for beginning of value",
		err.Error())
	assert.Equal(t, "", targetChainType)

	tx = schema.Transaction{
		Action:    tokSchema.TxActionMint,
		ChainType: tokSchema.ChainTypeEth,
		Data:      "{\"targetChainType\":\"mmmm\"}",
	}
	targetChainType, err = GetTargetChainTypeFromData(tx.Data, tx.Action,
		tx.ChainType)
	assert.NoError(t, err)
	assert.Equal(t, "mmmm", targetChainType)

	tx = schema.Transaction{
		Action:    tokSchema.TxActionMint,
		ChainType: tokSchema.ChainTypeArweave,
		Data:      "{\"aaa\":\"bbb\"}",
	}
	targetChainType, err = GetTargetChainTypeFromData(tx.Data, tx.Action,
		tx.ChainType)
	assert.NoError(t, err)
	assert.Equal(t, tokSchema.ChainTypeArweave, targetChainType)

	// AR token burn or mint must have data{targetChainType}
	tx = schema.Transaction{
		Action:    tokSchema.TxActionMint,
		ChainType: tokSchema.ChainTypeCrossArEth,
		Data:      "{\"targetChainType\":\"arweave\"}",
	}
	targetChainType, err = GetTargetChainTypeFromData(tx.Data, tx.Action,
		tx.ChainType)
	assert.NoError(t, err)
	assert.Equal(t, "arweave", targetChainType)

	tx = schema.Transaction{
		Action:    tokSchema.TxActionMint,
		ChainType: tokSchema.ChainTypeCrossArEth,
		Data:      "",
	}
	targetChainType, err = GetTargetChainTypeFromData(tx.Data, tx.Action,
		tx.ChainType)
	assert.Equal(t, "unexpected end of JSON input", err.Error())
	assert.Equal(t, "", targetChainType)

	tx = schema.Transaction{
		Action:    tokSchema.TxActionMint,
		ChainType: tokSchema.ChainTypeCrossArEth,
		Data:      "{\"tttt\":\"bbb\"}",
	}
	targetChainType, err = GetTargetChainTypeFromData(tx.Data, tx.Action,
		tx.ChainType)
	assert.NoError(t, err)
	assert.Equal(t, tokSchema.ChainTypeArweave, targetChainType)

	tx = schema.Transaction{
		Action:    tokSchema.TxActionBurn,
		ChainType: tokSchema.ChainTypeCrossArEth,
		Data:      "",
	}
	targetChainType, err = GetTargetChainTypeFromData(tx.Data, tx.Action,
		tx.ChainType)
	assert.NoError(t, err)
	assert.Equal(t, tokSchema.ChainTypeArweave, targetChainType)

	tx = schema.Transaction{
		Action:    tokSchema.TxActionBurn,
		ChainType: tokSchema.ChainTypeCrossArEth,
		Data:      "{\"tttt\":\"bbb\"}",
	}
	targetChainType, err = GetTargetChainTypeFromData(tx.Data, tx.Action,
		tx.ChainType)
	assert.NoError(t, err)
	assert.Equal(t, tokSchema.ChainTypeArweave, targetChainType)
}
