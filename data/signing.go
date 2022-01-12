package data

import (
	"github.com/ChainSQL/go-chainsql-api/common"
	"github.com/ChainSQL/go-chainsql-api/crypto"
)

func Sign(s Signer, key crypto.Key, sequence *uint32, keyType common.KeyType) error {
	s.InitialiseForSigning(key.Type())
	copy(s.GetPublicKey().Bytes(), key.Public(sequence))
	hash, msg, err := SigningHash(s, keyType)
	if err != nil {
		return err
	}
	sig, err := crypto.Sign(key, hash.Bytes(), sequence, append(s.SigningPrefix().Bytes(), msg...))
	if err != nil {
		return err
	}
	*s.GetSignature() = VariableLength(sig)
	hash, _, err = Raw(s,keyType)
	if err != nil {
		return err
	}
	copy(s.GetHash().Bytes(), hash.Bytes())
	return nil
}

func CheckSignature(s Signer, keyType common.KeyType) (bool, error) {
	hash, msg, err := SigningHash(s, keyType)
	if err != nil {
		return false, err
	}
	return crypto.Verify(s.GetPublicKey().Bytes(), hash.Bytes(), msg, s.GetSignature().Bytes())
}
