package cache

import (
	"encoding/json"
	"fmt"
	"strconv"
	"sync"

	"github.com/Hubmakerlabs/hoover/pkg/arweave/ethrpc"
	"github.com/Hubmakerlabs/hoover/pkg/arweave/go-everpay/account"
	"github.com/Hubmakerlabs/hoover/pkg/arweave/go-everpay/cache/schema"
	comm "github.com/Hubmakerlabs/hoover/pkg/arweave/go-everpay/common"
	"github.com/Hubmakerlabs/hoover/pkg/arweave/go-everpay/pay"
	paySchema "github.com/Hubmakerlabs/hoover/pkg/arweave/go-everpay/pay/schema"
	tokSchema "github.com/Hubmakerlabs/hoover/pkg/arweave/go-everpay/token/schema"
	arTypes "github.com/Hubmakerlabs/hoover/pkg/arweave/goar/types"
	"github.com/Hubmakerlabs/hoover/pkg/arweave/goar/utils"
	"github.com/Hubmakerlabs/hoover/pkg/eth/common"
	"github.com/Hubmakerlabs/hoover/pkg/eth/crypto"
)

var log = comm.NewLog("cache")

type Cache struct {
	rootHash     []byte // hash generate with arid & everHash
	everRootHash []byte // hash generate only with everHash

	pendingTxQueue    []*schema.TxResponse            // pending everTx queue
	txs               []*schema.TxResponse            // list of tx
	txByHash          map[string]*schema.TxResponse   // everHash -> Tx
	txsByAcc          map[string][]*schema.TxResponse // accid -> []Tx
	mintedByChainHash map[string]*schema.TxResponse   // chainHash -> Tx

	burningByEverHash    map[string]*schema.TxResponse // everHash -> Tx
	expressingByEverHash map[string]*schema.TxResponse // everHash -> Tx

	txCount int64 // everTx count, the same as latest schema.TxResponse.RawId

	lock sync.RWMutex
}

func New() *Cache {
	return &Cache{
		txs:               []*schema.TxResponse{},
		txByHash:          map[string]*schema.TxResponse{},
		txsByAcc:          map[string][]*schema.TxResponse{},
		mintedByChainHash: map[string]*schema.TxResponse{},

		burningByEverHash:    map[string]*schema.TxResponse{},
		expressingByEverHash: map[string]*schema.TxResponse{},
	}
}

func (c *Cache) SnapSet(txs []*schema.TxResponse) {
	c.txs = txs
	for _, tx := range txs {
		c.txCount += 1
		tx.RawId = c.txCount
		c.txByHash[tx.EverHash] = tx
		c.addTxsByAcc(tx)
		c.addMintedByChainHash(tx)
		c.addBurningTx(tx)
		c.addExpressingTx(tx)

		id, _ := utils.Base64Decode(tx.ID)
		everHash := common.FromHex(tx.EverHash)
		c.rootHash = RootHash(c.rootHash, id, everHash)
		c.everRootHash = RootHash(c.everRootHash, everHash)
	}
}

func (c *Cache) GetTxByHash(everHash string) *schema.TxResponse {
	c.lock.RLock()
	defer c.lock.RUnlock()

	return c.txByHash[everHash]
}

func (c *Cache) AddTx(tx paySchema.Transaction,
	internalErr *paySchema.InternalErr) {
	c.lock.Lock()
	defer c.lock.Unlock()

	nonce, _ := strconv.ParseInt(tx.Nonce, 10, 64)
	everHash := tx.HexHash()

	c.txCount += 1
	rTx := schema.TxResponse{
		RawId:             c.txCount,
		ID:                tx.ArTxID,
		TokenSymbol:       tx.TokenSymbol,
		Action:            tx.Action,
		From:              tx.From,
		To:                tx.To,
		Amount:            tx.Amount,
		Fee:               tx.Fee,
		FeeRecipient:      tx.FeeRecipient,
		Nonce:             nonce,
		TokenID:           tx.TokenID,
		ChainType:         tx.ChainType,
		ChainID:           tx.ChainID,
		Data:              tx.Data,
		Version:           tx.Version,
		Sig:               tx.Sig,
		EverHash:          everHash,
		Status:            schema.TxStatusConfirmed,
		InternalStatus:    InternalErrToStatus(internalErr).Marshal(),
		Timestamp:         tx.ArTimestamp,
		TargetChainTxHash: "",
		Express:           schema.Express{},
	}

	c.txs = append(c.txs, &rTx)
	c.txByHash[everHash] = &rTx
	c.addTxsByAcc(&rTx)
	c.addMintedByChainHash(&rTx)
	c.pushPendingTx(&rTx)

	c.addBurningTx(&rTx)
	c.addExpressingTx(&rTx)
}

func (c *Cache) addTxsByAcc(tx *schema.TxResponse) {
	_, from, err := account.IDCheck(tx.From)
	if err != nil {
		log.Error("add to cache failed, IDCheck failed", "from", from, "err",
			err)
		return
	}
	_, to, err := account.IDCheck(tx.To)
	if err != nil {
		log.Error("add to cache failed, IDCheck failed", "to", from, "err", err)
		return
	}

	// need add to cache address
	needAddMap := map[string]struct{}{
		from: {},
		to:   {},
	}

	switch tx.Action {
	case tokSchema.TxActionBurn:
		// burn(withdraw) tx do not show in receiver's txs, For front-end display
		if from != to { // delete to
			delete(needAddMap, to)
		}
		break
	case paySchema.TxActionBundle:
		// decode bundle
		bundleData := paySchema.BundleData{}
		if err = json.Unmarshal([]byte(tx.Data), &bundleData); err != nil {
			log.Error("can not unmarshal bundle data", "err", err, "everHash",
				tx.EverHash)
			return
		}
		// bundle item from and to de-duplication
		for _, item := range bundleData.Bundle.Items {
			needAddMap[item.From] = struct{}{}
			needAddMap[item.To] = struct{}{}
		}
	}

	// add to cache
	for addr := range needAddMap {
		if _, ok := c.txsByAcc[addr]; !ok {
			c.txsByAcc[addr] = make([]*schema.TxResponse, 0)
		}
		c.txsByAcc[addr] = append(c.txsByAcc[addr], tx)
	}
}

func (c *Cache) addMintedByChainHash(tx *schema.TxResponse) {
	if tx.Action != tokSchema.TxActionMint {
		return
	}
	targetChainTxHash, err := GetMintTargetTxHash(tx.ChainType, tx.Data)
	if err != nil {
		return
	}
	tx.TargetChainTxHash = targetChainTxHash
	c.mintedByChainHash[targetChainTxHash] = tx
}

func (c *Cache) addBurningTx(tx *schema.TxResponse) {
	if tx.Action != tokSchema.TxActionBurn {
		return
	}

	if tx.TargetChainTxHash != "" {
		return
	}

	c.burningByEverHash[tx.EverHash] = tx
}

func (c *Cache) pushPendingTx(tx *schema.TxResponse) {
	// 存在则不插入
	for _, tt := range c.pendingTxQueue {
		if tt.EverHash == tx.EverHash {
			return
		}
	}

	c.pendingTxQueue = append(c.pendingTxQueue, tx)
}

func (c *Cache) addExpressingTx(tx *schema.TxResponse) {
	exp := struct {
		AppId          string
		WithdrawAction string
	}{}
	if err := json.Unmarshal([]byte(tx.Data), &exp); err != nil {
		return
	}

	if exp.AppId != "express" && exp.WithdrawAction != "pay" {
		return
	}

	c.expressingByEverHash[tx.EverHash] = tx
}

func GetMintTargetTxHash(everTxChainType, everTxData string) (targetChainTxHash string,
	err error) {
	targetChainType, err := pay.GetTargetChainTypeFromData(everTxData,
		tokSchema.TxActionMint, everTxChainType)
	if err != nil {
		return "", err
	}
	switch targetChainType { // todo more chain
	case tokSchema.OracleEthChainType, tokSchema.OracleMoonChainType, tokSchema.OracleCfxChainType:
		ethTx := ethrpc.Transaction{}
		if err := json.Unmarshal([]byte(everTxData), &ethTx); err != nil {
			log.Error("tx data unmarshal failed", "data", everTxData, "err",
				err)
			return "", err
		}
		return ethTx.Hash, nil
	case tokSchema.OracleArweaveChainType:
		arTx := arTypes.Transaction{}
		if err = json.Unmarshal([]byte(everTxData), &arTx); err != nil {
			log.Error("tx data unmarshal failed, data", everTxData, "err", err)
			return
		}
		return arTx.ID, nil
	default:
		err = fmt.Errorf("not support this targetChainType: %s",
			targetChainType)
		return
	}
}

// RootHash: packaged unique rootHash
func RootHash(prevRoot []byte, ids ...[]byte) []byte {
	rawData := prevRoot
	for _, id := range ids {
		rawData = append(rawData, id...)
	}
	return crypto.Keccak256(rawData)
}

func InternalErrToStatus(internalErr *paySchema.InternalErr) schema.InternalStatus {
	internalState := schema.InternalStatus{}
	if internalErr == nil {
		internalState = schema.InternalStatus{Status: schema.InternalStatusSuccess}
	} else {
		internalState = schema.InternalStatus{
			Status:      schema.InternalStatusFailed,
			InternalErr: internalErr,
		}
	}

	return internalState
}
