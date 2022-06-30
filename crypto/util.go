package crypto

import (
	"crypto/sha256"
	"crypto/sha512"
	"encoding/hex"
	"fmt"
	"hash"

	"github.com/ChainSQL/go-chainsql-api/common"
	"golang.org/x/crypto/ripemd160"
	"golang.org/x/crypto/sha3"
)

// Write operations in a hash.Hash never return an error

// Returns first 32 bytes of a SHA512 of the input bytes
func Sha512Half(b []byte) []byte {
	hasher := sha512.New()
	hasher.Write(b)
	return hasher.Sum(nil)[:32]
}

// Returns first 16 bytes of a SHA512 of the input bytes
func Sha512Quarter(b []byte) []byte {
	hasher := sha512.New()
	hasher.Write(b)
	return hasher.Sum(nil)[:16]
}

func DoubleSha256(b []byte) []byte {
	hasher := sha256.New()
	hasher.Write(b)
	sha := hasher.Sum(nil)
	hasher.Reset()
	hasher.Write(sha)
	return hasher.Sum(nil)
}

func Sha256RipeMD160(b []byte) []byte {
	ripe := ripemd160.New()
	sha := sha256.New()
	sha.Write(b)
	ripe.Write(sha.Sum(nil))
	return ripe.Sum(nil)
}

func H2B(s string) ([]byte, error) {
	return hex.DecodeString(s)
}

func B2H(b []byte) string {
	return fmt.Sprintf("%X", b)
}

func B2H32(b [32]byte) string {
	return fmt.Sprintf("%X", b)
}

type KeccakState interface {
	hash.Hash
	Read([]byte) (int, error)
}

// NewKeccakState creates a new KeccakState
func NewKeccakState() KeccakState {
	return sha3.NewLegacyKeccak256().(KeccakState)
}

// Keccak256 calculates and returns the Keccak256 hash of the input data.
func Keccak256(data ...[]byte) []byte {
	b := make([]byte, 32)
	d := NewKeccakState()
	for _, b := range data {
		d.Write(b)
	}
	d.Read(b)
	return b
}

// Keccak256Hash calculates and returns the Keccak256 hash of the input data,
// converting it to an internal Hash data structure.
func Keccak256Hash(data ...[]byte) (h common.Hash) {
	d := NewKeccakState()
	for _, b := range data {
		d.Write(b)
	}
	d.Read(h[:])
	return h
}
