package data

import (
	"github.com/ChainSQL/go-chainsql-api/crypto"
)

func Sign(s Signer, key crypto.Key, sequence *uint32) error {
	switch key.Type() {
	case crypto.Ed25519:
		s.InitialiseForSigning(Ed25519)
		break
	case crypto.ECDSA:
		s.InitialiseForSigning(ECDSA)
		break
	case crypto.SoftGMAlg:
		s.InitialiseForSigning(SoftGMAlg)
		break
	default:
		s.InitialiseForSigning(ECDSA)
	}
	copy(s.GetPublicKey().Bytes(), key.Public(sequence))
	hash, msg, err := SigningHash(s)
	if err != nil {
		return err
	}
	sig, err := crypto.Sign(key, hash.Bytes(), sequence, append(s.SigningPrefix().Bytes(), msg...))
	if err != nil {
		return err
	}
	*s.GetSignature() = VariableLength(sig)
	hash, _, err = Raw(s)
	if err != nil {
		return err
	}
	copy(s.GetHash().Bytes(), hash.Bytes())
	return nil
}

func CheckSignature(s Signer) (bool, error) {
	hash, msg, err := SigningHash(s)
	if err != nil {
		return false, err
	}
	return crypto.Verify(s.GetPublicKey().Bytes(), hash.Bytes(), msg, s.GetSignature().Bytes())
}
