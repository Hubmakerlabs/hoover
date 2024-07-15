package utils

import (
	"crypto"
	"crypto/rand"
	"crypto/rsa"

	"github.com/minio/sha256-simd"
)

func GenerateRsaKey(bits int) (*rsa.PrivateKey, error) {
	return rsa.GenerateKey(rand.Reader, bits)
}

func Sign(msg []byte, prvKey *rsa.PrivateKey) ([]byte, error) {
	hashed := sha256.Sum256(msg)

	return rsa.SignPSS(rand.Reader, prvKey, crypto.SHA256, hashed[:],
		&rsa.PSSOptions{
			SaltLength: rsa.PSSSaltLengthAuto,
			Hash:       crypto.SHA256,
		})
}

func Verify(msg []byte, pubKey *rsa.PublicKey, sign []byte) error {
	hashed := sha256.Sum256(msg)

	return rsa.VerifyPSS(pubKey, crypto.SHA256, hashed[:], sign,
		&rsa.PSSOptions{
			SaltLength: rsa.PSSSaltLengthAuto,
			Hash:       crypto.SHA256,
		})
}
