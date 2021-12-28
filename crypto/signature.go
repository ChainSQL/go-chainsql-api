package crypto

import (
	"crypto/ed25519"
	"encoding/hex"
	"fmt"

	"github.com/btcsuite/btcd/btcec"
)

func Sign(key Key, hash []byte, sequence *uint32, msg []byte) ([]byte, error) {
	switch key.Type() {
	case Ed25519:
		return signEd25519(key.Private(sequence), msg)
	case ECDSA:
		return signECDSA(key.Private(sequence), hash)
	case SoftGMAlg:
		return signSoftSM(key.Private(sequence), msg)
	default:
		return nil, fmt.Errorf("Unknown private key format")
	}
}

func Verify(publicKey, hash, msg, signature []byte) (bool, error) {
	switch publicKey[0] {
	case 0xED:
		return verifyEd25519(publicKey, signature, msg)
	case 0x02, 0x03:
		return verifyECDSA(publicKey, signature, hash)
	default:
		return false, fmt.Errorf("Unknown public key format")
	}
}

func signEd25519(privateKey, msg []byte) ([]byte, error) {
	return ed25519.Sign(privateKey, msg)[:], nil
}

func verifyEd25519(pubKey, signature, msg []byte) (bool, error) {
	switch {
	case len(pubKey) != ed25519.PublicKeySize+1:
		return false, fmt.Errorf("Wrong public key length: %d", len(pubKey))
	case pubKey[0] != 0xED:
		return false, fmt.Errorf("Wrong public format:")
	case len(signature) != ed25519.SignatureSize:
		return false, fmt.Errorf("Wrong Signature length: %d", len(signature))
	default:
		return ed25519.Verify(pubKey[1:], msg, signature), nil
	}
}

// Returns DER encoded signature from input hash
func signECDSA(privateKey, hash []byte) ([]byte, error) {
	priv, _ := btcec.PrivKeyFromBytes(btcec.S256(), privateKey)
	sig, err := priv.Sign(hash)
	if err != nil {
		return nil, err
	}
	return sig.Serialize(), nil
}

// Verifies a hash using DER encoded signature
func verifyECDSA(pubKey, signature, hash []byte) (bool, error) {
	sig, err := btcec.ParseDERSignature(signature, btcec.S256())
	if err != nil {
		return false, err
	}
	pk, err := btcec.ParsePubKey(pubKey, btcec.S256())
	if err != nil {
		return false, nil
	}
	return sig.Verify(hash, pk), nil
}

func signSoftSM(privateKey, msg []byte) ([]byte, error) {
	keyPair, err := PrivKeyFromBytes(privateKey)
	if err != nil {
		return nil, err
	}
	sign, err := SmSign(keyPair, msg)
	if err != nil {
		return nil, err
	}
	signByte, err := hex.DecodeString(sign) // 转码
	return signByte, nil
}
