package ethrpc

import (
	"math/big"
)

type EthereumAPI interface {
	Web3ClientVersion() (string, error)
	Web3Sha3(data []byte) (string, error)
	NetVersion() (string, error)
	NetListening() (bool, error)
	NetPeerCount() (int, error)
	EthProtocolVersion() (string, error)
	EthSyncing() (*Syncing, error)
	EthCoinbase() (string, error)
	EthMining() (bool, error)
	EthHashrate() (int, error)
	EthGasPrice() (big.Int, error)
	EthAccounts() ([]string, error)
	EthBlockNumber() (int, error)
	EthGetBalance(address, block string) (big.Int, error)
	EthGetStorageAt(data string, position int, tag string) (string, error)
	EthGetTransactionCount(address, block string) (int, error)
	EthGetBlockTransactionCountByHash(hash string) (int, error)
	EthGetBlockTransactionCountByNumber(number int) (int, error)
	EthGetUncleCountByBlockHash(hash string) (int, error)
	EthGetUncleCountByBlockNumber(number int) (int, error)
	EthGetCode(address, block string) (string, error)
	EthSign(address, data string) (string, error)
	EthSendTransaction(transaction T) (string, error)
	EthSendRawTransaction(data string) (string, error)
	EthCall(transaction T, tag string) (string, error)
	EthEstimateGas(transaction T) (int, error)
	EthGetBlockByHash(hash string, withTransactions bool) (*Block, error)
	EthGetBlockByNumber(number int, withTransactions bool) (*Block, error)
	EthGetTransactionByHash(hash string) (*Transaction, error)
	EthGetTransactionByBlockHashAndIndex(blockHash string,
		transactionIndex int) (*Transaction, error)
	EthGetTransactionByBlockNumberAndIndex(blockNumber, transactionIndex int) (*Transaction,
		error)
	EthGetTransactionReceipt(hash string) (*TransactionReceipt, error)
	EthGetCompilers() ([]string, error)
	EthNewFilter(params FilterParams) (string, error)
	EthNewBlockFilter() (string, error)
	EthNewPendingTransactionFilter() (string, error)
	EthUninstallFilter(filterID string) (bool, error)
	EthGetFilterChanges(filterID string) ([]Log, error)
	EthGetFilterLogs(filterID string) ([]Log, error)
	EthGetLogs(params FilterParams) ([]Log, error)
}

var _ EthereumAPI = (*EthRPC)(nil)