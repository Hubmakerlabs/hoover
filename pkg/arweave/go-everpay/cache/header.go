package cache

import (
	"sync"

	"github.com/Hubmakerlabs/hoover/pkg/arweave/goar/utils"
)

type Header struct {
	RootHash        []byte `json:"rootHash"`
	EverRootHash    []byte `json:"everRootHash"`
	ArId            string `json:"arId"`
	EverHash        []byte `json:"everHash"`
	Nonce           string `json:"nonce"` // ever tx nonce
	BalanceRootHash string `json:"balanceRootHash"`
	locker          sync.RWMutex
}

func (h *Header) SetHeader(arId string, everHash []byte, txNonce string,
	balanceRootHash string) {
	id, _ := utils.Base64Decode(arId)
	rootHash := RootHash(h.RootHash, id, everHash)
	everRootHash := RootHash(h.EverRootHash, everHash)
	h.RootHash = rootHash
	h.EverHash = everHash
	h.EverRootHash = everRootHash
	h.ArId = arId
	h.Nonce = txNonce
	h.BalanceRootHash = balanceRootHash
}

func (h *Header) GetHeader() Header {
	h.locker.RLock()
	defer h.locker.RUnlock()
	return Header{
		RootHash:        h.RootHash,
		EverRootHash:    h.EverRootHash,
		ArId:            h.ArId,
		EverHash:        h.EverHash,
		Nonce:           h.Nonce,
		BalanceRootHash: h.BalanceRootHash,
	}
}
