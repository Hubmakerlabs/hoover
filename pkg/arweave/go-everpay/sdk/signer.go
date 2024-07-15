package sdk

import (
	"errors"

	"github.com/minio/sha256-simd"

	"github.com/Hubmakerlabs/hoover/pkg/arweave/goar"
	"github.com/Hubmakerlabs/hoover/pkg/arweave/goar/utils"
	"github.com/Hubmakerlabs/hoover/pkg/arweave/goether"
	"github.com/Hubmakerlabs/hoover/pkg/eth/common/hexutil"
)

const (
	RSASignerType = "RSASigner"
	EccSignerType = "EccSigner"
)

func (s *SDK) Sign(msg string) (string, error) {
	switch s.signerType {
	case RSASignerType:
		signer := s.signer.(*goar.Signer)
		hash := sha256.Sum256([]byte(msg))
		sig, err := signer.SignMsg(hash[:])
		if err != nil {
			return "", err
		}
		return utils.Base64Encode(sig) + "," + signer.Owner(), nil
	case EccSignerType:
		signer := s.signer.(*goether.Signer)
		sig, err := signer.SignMsg([]byte(msg))
		if err != nil {
			return "", err
		}
		return hexutil.Encode(sig), nil
	default:
		return "", errors.New("not found signer")
	}
}

func reflectSigner(signer interface{}) (signerType string, signerAddr string,
	err error) {
	if s, ok := signer.(*goar.Signer); ok {
		signerType = RSASignerType
		signerAddr = s.Address
		return
	}
	if s, ok := signer.(*goether.Signer); ok {
		signerType = EccSignerType
		signerAddr = s.Address.String()
		return
	}
	err = errors.New("not support this signer")
	return
}
