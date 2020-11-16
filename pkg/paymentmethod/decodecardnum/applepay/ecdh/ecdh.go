package ecdh

import (
	"crypto"
	"crypto/ecdsa"
	"io"
)

// The main interface for ECDH key exchange.
type ECDH interface {
	GenerateKey(io.Reader) (crypto.PrivateKey, crypto.PublicKey, error)
	Marshal(crypto.PublicKey) []byte
	Unmarshal([]byte) (crypto.PublicKey, bool)
	GenerateSharedSecret(*ecdsa.PrivateKey, *ecdsa.PublicKey) ([]byte, error)
}
