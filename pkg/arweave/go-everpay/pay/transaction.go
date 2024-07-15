package pay

import (
	"encoding/json"
	"math/big"

	"github.com/Hubmakerlabs/hoover/pkg/arweave/go-everpay/common"

	"github.com/Hubmakerlabs/hoover/pkg/arweave/go-everpay/pay/schema"

	"github.com/Hubmakerlabs/hoover/pkg/arweave/go-everpay/account"
	tokSchema "github.com/Hubmakerlabs/hoover/pkg/arweave/go-everpay/token/schema"
)

var log = common.NewLog("tx")

func AsTokenTx(t schema.Transaction) (tokenTx tokSchema.Transaction,
	err error) {
	amount, ok := new(big.Int).SetString(t.Amount, 10)
	if !ok {
		log.Error("invalid amount", "amount", t.Amount)
		err = ERR_INVALID_AMOUNT
		return
	}
	fee, ok := new(big.Int).SetString(t.Fee, 10)
	if !ok {
		log.Error("invalid fee", "fee", t.Fee)
		err = ERR_INVALID_FEE
		return
	}

	_, from, err := account.IDCheck(t.From)
	if err != nil {
		return
	}
	_, to, err := account.IDCheck(t.To)
	if err != nil {
		return
	}
	feeRecipient := "" // 费率为 0 的情况下，feeRecipient 可以为空
	if t.Fee != "0" || t.FeeRecipient != "" {
		_, feeRecipient, err = account.IDCheck(t.FeeRecipient)
		if err != nil {
			return
		}
	}

	var targetChainType string
	if t.Action == tokSchema.TxActionMint || t.Action == tokSchema.TxActionBurn {
		targetChainType, err = GetTargetChainTypeFromData(t.Data, t.Action,
			t.ChainType)
		if err != nil {
			err = ERR_INVALID_TARGET_CHAIN_TYPE
			return
		}
	}

	// bundle auto convert to transfer
	action := t.Action
	if t.Action == schema.TxActionBundle {
		action = tokSchema.TxActionTransfer
	}

	tokenTx = tokSchema.Transaction{
		Action:          action,
		From:            from, // notice: Case Sensitive !!!
		To:              to,   // notice: Case Sensitive !!!
		Amount:          amount,
		Fee:             fee,
		FeeRecipient:    feeRecipient, // notice: Case Sensitive !!!
		Data:            t.Data,
		TargetChainType: targetChainType,
	}

	return tokenTx, nil
}

func GetTargetChainTypeFromData(txData, txAction, txChainType string) (string,
	error) {
	/*
		1. mint tx must have json txData
		2. if parsed targetChainType is "", then we think this is cross to native chain,so targetChainType == txChainType
	*/

	targetChain := struct{ TargetChainType string }{}
	err := json.Unmarshal([]byte(txData), &targetChain)
	// mint tx must have json txData
	if err != nil && txAction == tokSchema.TxActionMint {
		return "", err
	}

	targetChainType := targetChain.TargetChainType
	if targetChainType == "" {
		nativeChainType, err := tokSchema.GetEverToNativeChainType(txChainType)
		if err != nil {
			return "", err
		}
		targetChainType = nativeChainType
	}

	return targetChainType, nil
}
