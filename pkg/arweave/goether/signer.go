package goether

import (
	"crypto/ecdsa"
	"io/ioutil"
	"math/big"
	"strings"

	"github.com/Hubmakerlabs/hoover/pkg/eth/accounts"
	"github.com/Hubmakerlabs/hoover/pkg/eth/common"
	"github.com/Hubmakerlabs/hoover/pkg/eth/common/hexutil"
	"github.com/Hubmakerlabs/hoover/pkg/eth/core/types"
	"github.com/Hubmakerlabs/hoover/pkg/eth/crypto"
	"github.com/Hubmakerlabs/hoover/pkg/eth/crypto/ecies"
	"github.com/Hubmakerlabs/hoover/pkg/eth/signer/core/apitypes"
)

type Signer struct {
	Address common.Address
	key     *ecdsa.PrivateKey
}

func NewSigner(prvHex string) (*Signer, error) {
	k, err := crypto.HexToECDSA(prvHex)
	if err != nil {
		return nil, err
	}

	return &Signer{
		key:     k,
		Address: crypto.PubkeyToAddress(k.PublicKey),
	}, nil
}

func NewSignerFromPath(prvPath string) (*Signer, error) {
	b, err := ioutil.ReadFile(prvPath)
	if err != nil {
		return nil, err
	}

	return NewSigner(strings.TrimSpace(string(b)))
}

func (s Signer) GetPrivateKey() *ecdsa.PrivateKey {
	return s.key
}

func (s Signer) GetPublicKey() []byte {
	return crypto.FromECDSAPub(&s.key.PublicKey)
}

func (s Signer) GetPublicKeyHex() string {
	return hexutil.Encode(s.GetPublicKey())
}

// SignTx DynamicFeeTx
func (s *Signer) SignTx(
	nonce int, to common.Address, amount *big.Int,
	gasLimit int, gasTipCap *big.Int, gasFeeCap *big.Int,
	data []byte, chainID *big.Int,
) (tx *types.Transaction, err error) {
	baseTx := &types.DynamicFeeTx{
		Nonce:     uint64(nonce),
		GasTipCap: gasTipCap,
		GasFeeCap: gasFeeCap,
		Gas:       uint64(gasLimit),
		To:        &to,
		Value:     amount,
		Data:      data,
	}
	return types.SignNewTx(s.key, types.LatestSignerForChainID(chainID), baseTx)
}

func (s *Signer) SignLegacyTx(
	nonce int, to common.Address, amount *big.Int,
	gasLimit int, gasPrice *big.Int,
	data []byte, chainID *big.Int,
) (tx *types.Transaction, err error) {
	return types.SignTx(
		types.NewTransaction(
			uint64(nonce), to, amount,
			uint64(gasLimit), gasPrice, data),
		types.NewEIP155Signer(chainID),
		s.key,
	)
}

func (s Signer) SignMsg(msg []byte) (sig []byte, err error) {
	hash := accounts.TextHash(msg)
	sig, err = crypto.Sign(hash, s.key)
	if err != nil {
		return
	}

	sig[64] += 27
	return
}

func (s Signer) SignTypedData(typedData apitypes.TypedData) (sig []byte,
	err error) {
	hash, err := EIP712Hash(typedData)
	if err != nil {
		return
	}

	sig, err = crypto.Sign(hash, s.key)
	if err != nil {
		return
	}

	sig[64] += 27
	return
}

// Decrypt decrypt
func (s Signer) Decrypt(ct []byte) ([]byte, error) {
	eciesPriv := ecies.ImportECDSA(s.key)
	return eciesPriv.Decrypt(ct, nil, nil)
}
